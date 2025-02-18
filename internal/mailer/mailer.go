package mailer

import (
	"bytes"
	"embed"
	"html/template"
	"time"

	"github.com/go-mail/mail/v2"
)

//go:embed "templates"
var templateFS embed.FS

// Mailer represents a mailer instance with methods for sending emails.
type Mailer struct {
	dialer *mail.Dialer
	sender string
}

// New creates a new Mailer instance.
//
// Parameters:
//
//	host - The SMTP server host
//	port - The SMTP server port
//	username - The SMTP server username
//	password - The SMTP server password
//	sender - The email address to use as the sender
//
// Returns:
//
//	*Mailer - A pointer to the newly created Mailer instance
func New(host string, port int, username, password, sender string) Mailer {

	dialer := mail.NewDialer(host, port, username, password)
	dialer.Timeout = 5 * time.Second

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}

// Send sends an email using the specified template and data.
//
// Parameters:
//
//	recipient - The recipient's email address
//	templateFile - The name of the email template file (e.g., "welcome.tmpl")
//	data - The data to pass to the template for rendering
//
// Returns:
//
//	error - If any error occurs during the process
func (m Mailer) Send(recipient, templateFile string, data any) error {

	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	plainBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(plainBody, "plainBody", data)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", plainBody.String())
	msg.AddAlternative("text/html", htmlBody.String())

	for range 3 {
		err = m.dialer.DialAndSend(msg)
		if nil == err {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return err
}
