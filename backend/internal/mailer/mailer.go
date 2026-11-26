package mailer

import (
	"bytes"
	htmltemplate "html/template"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	texttemplate "text/template"
	"time"

	"karsa/internal/config"
	"karsa/internal/models"
)

const sendTimeout = 10 * time.Second

var textBody = texttemplate.Must(texttemplate.New("text").Parse(`{{.Heading}}

{{range .Lines}}{{.}}

{{end}}{{.Button}}: {{.URL}}

Proyek demo untuk pembelajaran. Perusahaan dan lowongan bersifat fiktif.
`))

var htmlBody = htmltemplate.Must(htmltemplate.New("html").Parse(`<!doctype html>
<html lang="id">
<body style="margin:0;padding:24px;background:#FAFAF7;font-family:Inter,Arial,sans-serif;color:#14171F">
  <div style="max-width:560px;margin:0 auto;padding:32px;background:#FFFFFF;border-radius:12px">
    <p style="margin:0 0 24px;font-weight:800;color:#0F766E">Karsa</p>
    <h1 style="margin:0 0 16px;font-size:22px">{{.Heading}}</h1>
    {{range .Lines}}<p style="margin:0 0 12px;line-height:1.6">{{.}}</p>{{end}}
    <p style="margin:24px 0"><a href="{{.URL}}" style="padding:12px 20px;background:#0F766E;color:#FFFFFF;border-radius:8px;text-decoration:none">{{.Button}}</a></p>
    <p style="margin:0;font-size:12px;color:#5F6672">Proyek demo untuk pembelajaran. Perusahaan dan lowongan bersifat fiktif.</p>
  </div>
</body>
</html>
`))

type Mailer struct {
	Config config.Config
}

func (mailer Mailer) Send(job models.EmailJob) error {
	message, err := Build(mailer.Config.MailFrom, job)
	if err != nil {
		return err
	}
	connection, err := net.DialTimeout("tcp", net.JoinHostPort(mailer.Config.SMTPHost, mailer.Config.SMTPPort), sendTimeout)
	if err != nil {
		return err
	}
	defer connection.Close()
	if err := connection.SetDeadline(time.Now().Add(sendTimeout)); err != nil {
		return err
	}
	client, err := smtp.NewClient(connection, mailer.Config.SMTPHost)
	if err != nil {
		return err
	}
	if err := client.Mail(mailer.Config.MailFrom); err != nil {
		return err
	}
	if err := client.Rcpt(job.ToEmail); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write(message); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	_ = client.Quit()
	return nil
}

func Build(from string, job models.EmailJob) ([]byte, error) {
	var body bytes.Buffer
	parts := multipart.NewWriter(&body)
	for _, part := range []struct {
		contentType string
		render      func(*quotedprintable.Writer) error
	}{
		{"text/plain", func(writer *quotedprintable.Writer) error { return textBody.Execute(writer, job.Payload) }},
		{"text/html", func(writer *quotedprintable.Writer) error { return htmlBody.Execute(writer, job.Payload) }},
	} {
		section, err := parts.CreatePart(textproto.MIMEHeader{
			"Content-Type":              {part.contentType + "; charset=utf-8"},
			"Content-Transfer-Encoding": {"quoted-printable"},
		})
		if err != nil {
			return nil, err
		}
		writer := quotedprintable.NewWriter(section)
		if err := part.render(writer); err != nil {
			return nil, err
		}
		if err := writer.Close(); err != nil {
			return nil, err
		}
	}
	if err := parts.Close(); err != nil {
		return nil, err
	}

	headers := strings.Join([]string{
		"From: " + (&mail.Address{Name: "Karsa", Address: from}).String(),
		"To: " + (&mail.Address{Address: job.ToEmail}).String(),
		"Subject: " + mime.QEncoding.Encode("utf-8", job.Subject),
		"Date: " + time.Now().Format(time.RFC1123Z),
		"MIME-Version: 1.0",
		"Content-Type: multipart/alternative; boundary=" + parts.Boundary(),
	}, "\r\n")
	return append([]byte(headers+"\r\n\r\n"), body.Bytes()...), nil
}
