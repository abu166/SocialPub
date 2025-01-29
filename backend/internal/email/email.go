package email

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

// Load environment variables
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: Could not load .env file")
	}
}

// SendEmail sends a confirmation email with a code
func SendEmail(toEmail, code string) error {
	// Get SMTP details from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_USER")
	senderPassword := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpPort == "" || senderEmail == "" || senderPassword == "" {
		return fmt.Errorf("SMTP configuration is missing in environment variables")
	}

	// Message content
	subject := "Account Confirmation Code"
	message := fmt.Sprintf("Subject: %s\r\n\r\nYour confirmation code is: %s", subject, code)

	// SMTP authentication
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// Send email
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{toEmail}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Println("Confirmation email sent to", toEmail)
	return nil
}
