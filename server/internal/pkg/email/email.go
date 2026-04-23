package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/hellomyheart/zero-life/server/internal/config"
	"go.uber.org/zap"
)

// Send sends an HTML email using the configured SMTP server.
func Send(to, subject, htmlBody string) error {
	cfg := config.C.SMTP
	if cfg.Host == "" {
		zap.L().Warn("SMTP not configured, skipping email send")
		return nil
	}

	from := cfg.From
	if from == "" {
		from = cfg.User
	}

	// Build email message
	msg := strings.Builder{}
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", to))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	// Use TLS for port 465, STARTTLS for port 587
	if cfg.Port == 465 {
		return sendWithTLS(addr, from, to, cfg, msg.String())
	}
	return sendWithSTARTTLS(addr, from, to, cfg, msg.String())
}

func sendWithSTARTTLS(addr, from, to string, cfg config.SMTPConfig, msg string) error {
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("smtp dial: %w", err)
	}
	defer client.Close()

	// Send HELO/EHLO
	if err := client.Hello("localhost"); err != nil {
		return fmt.Errorf("smtp hello: %w", err)
	}

	// Start TLS
	tlsConfig := &tls.Config{ServerName: cfg.Host}
	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("smtp starttls: %w", err)
	}

	// Authenticate
	if cfg.User != "" {
		auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	// Send mail
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close: %w", err)
	}

	return client.Quit()
}

func sendWithTLS(addr, from, to string, cfg config.SMTPConfig, msg string) error {
	tlsConfig := &tls.Config{ServerName: cfg.Host}
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		return fmt.Errorf("smtp new client: %w", err)
	}
	defer client.Close()

	// Authenticate
	if cfg.User != "" {
		auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}

	// Send mail
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close: %w", err)
	}

	return client.Quit()
}

// SendPasswordReset sends a password reset email.
func SendPasswordReset(to, resetURL string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf(`
	<html>
	<body>
		<h2>Password Reset</h2>
		<p>You have requested to reset your password.</p>
		<p>Click the link below to reset your password:</p>
		<p><a href="%s">Reset Password</a></p>
		<p>This link will expire in 24 hours.</p>
		<p>If you did not request this, please ignore this email.</p>
	</body>
	</html>`, resetURL)
	return Send(to, subject, body)
}
