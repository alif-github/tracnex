package util

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type ToZipStruct struct {
	Filepath    string
	NewFileName string
}

func ZipFiles(filename string, files []string) error {
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err = newZipFile.Close(); err != nil {
			log.Println("Failed to close zip File", err)
		}
	}()

	zipWriter := zip.NewWriter(newZipFile)
	defer func() {
		if err = zipWriter.Close(); err != nil {
			log.Println("Failed to close zip File", err)
		}
	}()

	for _, file := range files {
		if err = AddFileToZip(zipWriter, file, ""); err != nil {
			return err
		}
	}
	return nil
}

func ZipFilesWithOtherName(filename string, files []ToZipStruct, readNewData bool) (temp []byte, err error) {
	defer func() {
		if readNewData {
			temp, err = readFile(filename)
		}
	}()

	newZipFile, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = newZipFile.Close()
	}()

	zipWriter := zip.NewWriter(newZipFile)
	defer func() {
		_ = zipWriter.Close()
	}()

	for _, file := range files {
		if err = AddFileToZipWithOtherName(zipWriter, file); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func readFile(filepath string) (result []byte, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer func() {
		_ = file.Close()
	}()

	stat, _ := file.Stat()

	result = make([]byte, stat.Size())
	_, err = file.Read(result)
	if err != nil {
		return
	}

	return
}

func AddFileToZipWithOtherName(zipWriter *zip.Writer, filename ToZipStruct) error {
	fileToZip, err := os.Open(filename.Filepath)
	if err != nil {
		return err
	}
	defer func() {
		if err = fileToZip.Close(); err != nil {
			log.Println("Failed to close zip File", err)
		}
	}()

	// Get the file information
	info, _ := fileToZip.Stat()

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return errors.New("ERROR_CREATING_ZIP_HEADER")
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename.NewFileName

	// Change to deflate to gain better compression
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return errors.New("ERROR_CREATING_ZIP_HEADER")
	}
	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return errors.New("ERROR_INSERTING_FILE_TO_ZIP")
	}
	return nil
}

func UnzipFromByte(data []byte) (output []*zip.File, err error) {
	unziped, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, errors.New("NOT_VALID_ZIP")
	}
	output = unziped.File
	return
}

func ReadZipFile(zf *zip.File) ([]byte, error) {
	f, err := zf.Open()
	if err != nil {
		return nil, errors.New("ERROR_READING_ZIP")
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Println("Failed to close file", err)
		}
	}()
	return ioutil.ReadAll(f)
}

func Unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}

	defer func() {
		if err = r.Close(); err != nil {
			log.Println("Failed to close file", err)
		}
	}()

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return filenames, err
			}
			continue
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		err = outFile.Close()
		if err != nil {
			return filenames, err
		}
		err = rc.Close()
		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}

func AddFileToZip(zipWriter *zip.Writer, basePath string, baseInZip string) error {
	basePath = strings.Replace(basePath, "\\", "/", -1)
	fileInfo, err := os.Stat(basePath)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		dat, err := ioutil.ReadFile(basePath)
		if err != nil {
			return err
		}

		f, err := zipWriter.Create(baseInZip + fileInfo.Name())
		if err != nil {
			return err
		}
		_, err = f.Write(dat)
		if err != nil {
			return err
		}
	} else {
		files, err := ioutil.ReadDir(basePath)
		fileName := strings.Split(basePath, "/")
		if err != nil {
			return err
		}
		if len(files) == 0 {
			zipName := baseInZip
			if baseInZip == "" {
				zipName = baseInZip + fileName[len(fileName)-1] + "/"
			}
			_, err := zipWriter.Create(zipName)
			if err != nil {
				return err
			}
		} else {
			if baseInZip == "" {
				baseInZip = fileName[len(fileName)-1] + "/"
			}
			for _, file := range files {
				if !file.IsDir() {
					newBase := basePath + "/" + file.Name()
					err = AddFileToZip(zipWriter, newBase, baseInZip)
					if err != nil {
						return err
					}
				} else if file.IsDir() {
					newBase := basePath + "/" + file.Name()
					err = AddFileToZip(zipWriter, newBase, baseInZip+file.Name()+"/")
					if err != nil {
						return err
					}
				}

			}
		}
	}
	return nil
}