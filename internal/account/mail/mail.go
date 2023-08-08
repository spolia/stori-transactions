package mail

import (
	"bytes"
	"context"
	"html/template"

	"gopkg.in/gomail.v2"
)

type Configuration struct {
	Sender     string
	Password   string
	SmtpServer string
	SmtpPort   int
}

type SMTPClient struct {
	dialer *gomail.Dialer
}

func New(conf Configuration) *SMTPClient {
	return &SMTPClient{dialer: gomail.NewDialer(conf.SmtpServer, conf.SmtpPort, conf.Sender, conf.Password)}
}

// SendEmail sends an email to the recipient with the movements summary
func (s *SMTPClient) SendEmail(ctx context.Context, recipient string, movementsByMonth map[string]int, totalBalance float64, avgDebitByMonth, avgCreditByMonth map[string]float64) error {

	// Generate the email content using the HTML template
	emailBody, err := s.generateEmailBody(recipient, movementsByMonth, totalBalance, avgDebitByMonth, avgCreditByMonth)
	if err != nil {
		return err
	}
	// Create the email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", s.dialer.Username)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", "Movements Summary")
	msg.SetBody("text/html", emailBody)

	// Embed company logo in the email
	msg.Embed("internal/account/mail/storilogo.png")

	// Send the email using SMTP server
	if err = s.dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

// generateEmailBody generates the email content using the HTML template
func (s *SMTPClient) generateEmailBody(recipient string, movementsByMonth map[string]int, totalBalance float64, avgDebitByMonth, avgCreditByMonth map[string]float64) (string, error) {
	tmpl := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>Transaction Summary</title>
		</head>
		<body>
			<img src="cid:companylogo" alt="Company Logo" />
			<h2>Transaction Summary</h2>
			<p>Total balance is {{printf "%.2f" .TotalBalance}}</p>
			{{range $month, $count := .MovementsByMonth}}
				<p>Number of movements in {{$month}}: {{$count}}</p>
				<p>Average debit amount in {{$month}}: {{printf "%.2f" (index $.AvgDebitByMonth $month)}}</p>
				<p>Average credit amount in {{$month}}: {{printf "%.2f" (index $.AvgCreditByMonth $month)}}</p>
			{{end}}
		</body>
		</html>
	`

	tmplData := struct {
		TotalBalance     float64
		MovementsByMonth map[string]int
		AvgCreditByMonth map[string]float64
		AvgDebitByMonth  map[string]float64
	}{
		TotalBalance:     totalBalance,
		MovementsByMonth: movementsByMonth,
		AvgCreditByMonth: avgCreditByMonth,
		AvgDebitByMonth:  avgDebitByMonth,
	}

	var emailBody bytes.Buffer
	tmplObj := template.Must(template.New("emailTemplate").Parse(tmpl))
	err := tmplObj.Execute(&emailBody, tmplData)
	if err != nil {
		return "", err
	}

	return emailBody.String(), nil
}
