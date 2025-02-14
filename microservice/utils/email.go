package utils

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"os"
)

func SendEmail(to, subject, body, attachment string) error {
	// Validate environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := "587"
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASS")
	if smtpHost == "" || smtpUser == "" || smtpPassword == "" {
		return fmt.Errorf("SMTP configuration is missing in environment variables")
	}

	// Create the email message
	m := gomail.NewMessage()
	m.SetHeader("From", smtpUser)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	// Attach the file if the path is provided
	if attachment != "" {
		m.Attach(attachment)
	}

	// Log SMTP settings for debugging
	log.Printf("SMTP Settings: Host=%s, Port=%s, User=%s", smtpHost, smtpPort, smtpUser)

	// Create a new dialer and send the email
	d := gomail.NewDialer(smtpHost, 587, smtpUser, smtpPassword)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v", err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("Email sent successfully")
	return nil
}
