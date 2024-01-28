package util

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"net/smtp"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"path/filepath"
	"strconv"
	"strings"
)

type RequestEmail struct {
	SenderName  string
	To          []mail.Address
	Cc          []mail.Address
	Bcc         []mail.Address
	Subject     string
	Body        string
	Attachments map[string][]byte
	Embeds      map[string][]byte
	ContentType string
}

type TemplateData struct {
	Title       string        `json:"Title"`
	Name        string        `json:"Name"`
	Description template.HTML `json:"Description"`
	Link        string        `json:"Link"`
	DueDate     string        `json:"DueDate"`
}

type TemplateDataActivationUserNexMile struct {
	Purpose     string        `json:"Purpose"`
	Name        string        `json:"Name"`
	ClientType  string        `json:"ClientType"`
	UniqueID1   string        `json:"UniqueID1"`
	CompanyName string        `json:"CompanyName"`
	UniqueID2   string        `json:"UniqueID2"`
	BranchName  string        `json:"BranchName"`
	SalesmanID  string        `json:"SalesmanID"`
	UserID      string        `json:"UserID"`
	Password    string        `json:"Password"`
	OTP         string        `json:"OTP"`
	Email       string        `json:"Email"`
	ClientID    string        `json:"ClientID"`
	AuthUserID  int64         `json:"AuthUserID"`
	Description template.HTML `json:"Description"`
	Link        string        `json:"Link"`
}

func NewRequestMail(subject string) *RequestEmail {
	return &RequestEmail{
		Subject:     subject,
		Attachments: make(map[string][]byte),
		Embeds:      make(map[string][]byte),
		ContentType: "text/html",
	}
}

func (input *RequestEmail) Send() error {
	var (
		all  []string
		body []byte
	)

	body, all = input.toByte()
	addr := fmt.Sprintf("%s:%s", config.ApplicationConfiguration.GetEmail().Host, strconv.Itoa(config.ApplicationConfiguration.GetEmail().Port))
	host := config.ApplicationConfiguration.GetEmail().Host
	auth := smtp.PlainAuth("", config.ApplicationConfiguration.GetEmail().Address, config.ApplicationConfiguration.GetEmail().Password, host)
	if err := smtp.SendMail(addr, auth, config.ApplicationConfiguration.GetEmail().Address, all, body); err != nil {
		return err
	}

	if config.ApplicationConfiguration.GetServerLogLevel() == constanta.LogLevelDebug {
		logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
		logModel.Status = 200
		logModel.Message = "Server Start in port : " + strconv.Itoa(config.ApplicationConfiguration.GetServerPort())
		util.LogInfo(logModel.ToLoggerObject())
	}

	return nil
}

func (input *RequestEmail) toByte() (output []byte, allEmail []string) {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(input.Attachments) > 0
	withEmbeds := len(input.Embeds) > 0
	var (
		to, cc, bcc, all []string
	)

	for _, mail := range input.To {
		to = append(to, mail.String())
		all = append(all, mail.Address)
	}

	for _, mail := range input.Cc {
		cc = append(cc, mail.String())
		all = append(all, mail.Address)
	}

	for _, mail := range input.Bcc {
		bcc = append(bcc, mail.String())
		all = append(all, mail.Address)
	}

	headers := make(map[string]string)
	headers["Subject"] = input.Subject
	headers["To"] = strings.Join(to, ", ")
	headers["Cc"] = strings.Join(cc, ", ")
	headers["Bcc"] = strings.Join(bcc, ", ")

	for key, value := range headers {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else if withEmbeds {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/related; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	}

	buf.WriteString("Content-Type: text/html; charset=\"UTF-8\"\n")

	buf.WriteString(input.Body)
	if withAttachments {
		for k, v := range input.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}
	}

	if withEmbeds {
		for k, v := range input.Embeds {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s; name=\"%s\"\n", http.DetectContentType(v), k))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-ID: <%s>\n", k))
			buf.WriteString(fmt.Sprintf("Content-Disposition: inline; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}
	}

	if withEmbeds || withAttachments {
		buf.WriteString("--")
	}

	return buf.Bytes(), all
}

func (input *RequestEmail) AttachFile(src string) error {
	file, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	input.Attachments[fileName] = file
	return nil
}

func (input *RequestEmail) EmbedFile(src string) error {
	file, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	_, fileName := filepath.Split(src)
	input.Embeds[fileName] = file
	return nil
}

func (input *RequestEmail) ParseTemplate(templateFileName, local string, data interface{}) error {
	tmpl := template.Must(template.New(local).ParseFiles(templateFileName))
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		return err
	}
	input.Body = buf.String()
	return nil
}

func (input *RequestEmail) GenerateEmailTemplate(templateSrc string, templateData interface{}, locale string) error {
	return input.ParseTemplate(templateSrc, locale, templateData)
}

func (input *RequestEmail) SendEmail() error {
	var (
		mime        = fmt.Sprintf("MIME-version: 1.0;\nContent-Type: %s; charset=\"UTF-8\";\n\n", input.ContentType)
		addr        = fmt.Sprintf("%s:%s", config.ApplicationConfiguration.GetEmail().Host, strconv.Itoa(config.ApplicationConfiguration.GetEmail().Port))
		host        = config.ApplicationConfiguration.GetEmail().Host
		to, cc, all []string
		content     string
	)

	fmt.Println("address : ", addr)

	for _, itemMailTo := range input.To {
		to = append(to, itemMailTo.String())
		all = append(all, itemMailTo.Address)
	}

	for _, itemMailCc := range input.Cc {
		cc = append(cc, itemMailCc.String())
		all = append(all, itemMailCc.Address)
	}

	headers := make(map[string]string)
	headers["Subject"] = input.Subject
	headers["To"] = strings.Join(to, ", ")
	headers["Cc"] = strings.Join(cc, ", ")

	if input.SenderName != "" {
		headers["From"] = fmt.Sprintf(`"%s" <%s>`, input.SenderName, config.ApplicationConfiguration.GetEmail().Address)
	}

	for key, value := range headers {
		content += fmt.Sprintf("%s: %s\r\n", key, value)
	}

	content += mime + "\n" + input.Body

	auth := smtp.PlainAuth("", config.ApplicationConfiguration.GetEmail().Address, config.ApplicationConfiguration.GetEmail().Password, host)
	if err := smtp.SendMail(addr, auth, config.ApplicationConfiguration.GetEmail().Address, all, []byte(content)); err != nil {
		return err
	}

	return nil
}

func SendMessageToEmail(toEmail string, message string, subject string, logModel applicationModel.LoggerModel) {
	receiver := make([]string, 1)
	receiver[0] = toEmail

	if strings.Contains(toEmail, ",") {
		receiver = strings.Split(toEmail, ",")
	}
	err := SendEmail(receiver, subject, message, nil)
	if err != nil {
		logModel.Status = 500
		logModel.Message = "Sending " + subject + " Failed" + ", with err " + err.Error()
	} else {
		logModel.Status = 200
		logModel.Message = "Sending " + subject + " Success"
	}
	InputLog(errorModel.GenerateNonErrorModel(), logModel)
}

func SendEmail(to []string, subject string, message string, attachments map[string][]byte) error {
	defer func() {
		if r := recover(); r != nil {
			//TopRecoverLog("SendEmail.go", "SendEmail", r)
		}
	}()
	host := config.ApplicationConfiguration.GetEmail().Host
	port := strconv.Itoa(config.ApplicationConfiguration.GetEmail().Port)
	servername := host + ":" + port
	witAttachment := len(attachments) > 0

	from := config.ApplicationConfiguration.GetEmail().Address
	subj := subject

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = strings.Join(to, ",")
	headers["Subject"] = subj
	if witAttachment {
		headers["MIME-Version"] = "1.0"
		headers["Content-Type"] = "multipart/mixed; boundary=f46d043c813270fc6b04c2d223da"
	}

	content := ""
	for k, v := range headers {
		content += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	if witAttachment {

		content += "\n--f46d043c813270fc6b04c2d223da\nContent-Type: text/plain; charset=utf-8\n"
	}
	content += "\r\n" + message

	if witAttachment {
		content += "\r\n\r\n"
		for fileName, fileContent := range attachments {
			content += "\n--f46d043c813270fc6b04c2d223da\n"
			headersFile := make(map[string]string)

			ext := filepath.Ext(fileName)
			mimetype := mime.TypeByExtension(ext)
			if mimetype != "" {
				headersFile["Content-Type"] = mimetype + "; charset=utf-8"
			}

			headersFile["Content-Transfer-Encoding"] = "base64"
			//handling email scheduler ada perubahan format nama file

			headersFile["Content-Disposition"] = `attachment; filename="` + fileName + `"`

			for k, v := range headersFile {
				content += fmt.Sprintf("%s: %s\r\n", k, v)
			}
			content += "\n"

			b := make([]byte, base64.StdEncoding.EncodedLen(len(fileContent)))
			base64.StdEncoding.Encode(b, fileContent)

			// write base64 content in lines of up to 76 chars
			for i, l := 0, len(b); i < l; i++ {
				content += string(b[i])
				if (i+1)%76 == 0 {
					content += "\r\n"
				}
			}

		}
	}

	fmt.Println("Email Address : ", config.ApplicationConfiguration.GetEmail().Address)
	fmt.Println("Email Password : ", config.ApplicationConfiguration.GetEmail().Password)
	fmt.Println("Email Host : ", config.ApplicationConfiguration.GetEmail().Host)
	fmt.Println("Email port : ", config.ApplicationConfiguration.GetEmail().Port)

	client, err := smtp.Dial(servername)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		return err
	}

	auth := smtp.PlainAuth("", config.ApplicationConfiguration.GetEmail().Address, config.ApplicationConfiguration.GetEmail().Password, host)
	if err = client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(config.ApplicationConfiguration.GetEmail().Address); err != nil {
		return err
	}
	//client := serverconfig.ServerAttribute.EmailClient

	for i := range to {
		if err := client.Rcpt(to[i]); err != nil {
			return err
		}

	}

	writeCloser, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writeCloser.Write([]byte(content))
	if err != nil {
		return err
	}

	defer func() {
		_ = writeCloser.Close()
		_ = client.Quit()
	}()

	return nil
}