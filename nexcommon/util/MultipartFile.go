package util

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

var multipartByReader = &multipart.Form{
	Value: make(map[string][]string),
	File:  make(map[string][]*multipart.FileHeader),
}

func ParseMultipartForm(r *http.Request, maxMemory int64, tempFile string) error {
	if r.MultipartForm == multipartByReader {
		return errors.New("http: multipart handled by MultipartReader")
	}
	if r.Form == nil {
		err := r.ParseForm()
		if err != nil {
			return err
		}
	}
	if r.MultipartForm != nil {
		return nil
	}

	mr, err := multipartReader(r, false)
	if err != nil {
		return err
	}

	f, err := readForm(mr, maxMemory, tempFile)
	if err != nil {
		return err
	}

	if r.PostForm == nil {
		r.PostForm = make(url.Values)
	}
	for k, v := range f.Value {
		r.Form[k] = append(r.Form[k], v...)
		// r.PostForm should also be populated. See Issue 9305.
		r.PostForm[k] = append(r.PostForm[k], v...)
	}

	r.MultipartForm = f

	return nil
}

func multipartReader(r *http.Request, allowMixed bool) (*multipart.Reader, error) {
	v := r.Header.Get("Content-Type")
	if v == "" {
		return nil, http.ErrNotMultipart
	}
	d, params, err := mime.ParseMediaType(v)
	if err != nil || !(d == "multipart/form-data" || allowMixed && d == "multipart/mixed") {
		return nil, http.ErrNotMultipart
	}
	boundary, ok := params["boundary"]
	if !ok {
		return nil, http.ErrMissingBoundary
	}
	return multipart.NewReader(r.Body, boundary), nil
}

func readForm(r *multipart.Reader, maxMemory int64, tempFile string) (form *multipart.Form, err error) {
	form = &multipart.Form{Value: make(map[string][]string), File: make(map[string][]*multipart.FileHeader)}
	defer func() {
		if err != nil {
			_ = form.RemoveAll()
		}
	}()

	maxValueBytes := maxMemory + int64(10<<20)
	_ = os.Remove(tempFile)

	err = ioutil.WriteFile(tempFile, nil, 0660)
	if err != nil {
		return
	}

	for {
		p, err := r.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		name := p.FormName()
		if name == "" {
			continue
		}
		filename := p.FileName()

		var b bytes.Buffer

		if filename == "" {
			n, err := io.CopyN(&b, p, maxValueBytes+1)
			if err != nil && err != io.EOF {
				return nil, err
			}
			maxValueBytes -= n
			if maxValueBytes < 0 {
				return nil, multipart.ErrMessageTooLarge
			}
			form.Value[name] = append(form.Value[name], b.String())
			continue
		}

		fh := &multipart.FileHeader{
			Filename: filename,
			Header:   p.Header,
		}

		err = func() error {
			var fileTemp *os.File
			defer func() {
				if fileTemp != nil {
					_ = fileTemp.Close()
				}
			}()

			fileTemp, err = os.OpenFile(tempFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModeAppend)
			if err != nil {
				fileTemp, err = ioutil.TempFile("", "multipart-")
				if err != nil {
					return err
				}
				return err
			}

			size, err := io.Copy(fileTemp, io.MultiReader(&b, p))
			if cerr := fileTemp.Close(); err == nil {
				err = cerr
			}
			if err != nil {
				_ = os.Remove(fileTemp.Name())
				return err
			}
			fh.Size = size

			return nil
		}()
		if err != nil {
			return nil, err
		}
		form.File[name] = append(form.File[name], fh)
	}

	return form, nil
}

func ReadMultipartFile(request *http.Request, key string) (file multipart.File, filename string, size int64, errs error) {
	file, handler, errs := request.FormFile(key)
	if errs != nil {
		return
	}

	defer func() {
		errs = file.Close()
		if errs != nil {
			return
		}
	}()

	return file, handler.Filename, handler.Size, nil
}

func getFileTemp(tempFile string, b bytes.Buffer, p *multipart.Part, fh *multipart.FileHeader) (err error) {
	var fileTemp *os.File
	defer func() {
		if fileTemp != nil {
			_ = fileTemp.Close()
		}
	}()

	fileTemp, err = os.OpenFile(tempFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModeAppend)
	if err != nil {
		fileTemp, err = ioutil.TempFile("", "multipart-")
		if err != nil {
			return err
		}
		return err
	}

	size, err := io.Copy(fileTemp, io.MultiReader(&b, p))
	if cerr := fileTemp.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(fileTemp.Name())
		return err
	}
	fh.Filename = fileTemp.Name()
	fh.Size = size

	return nil
}
