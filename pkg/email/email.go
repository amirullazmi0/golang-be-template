package email

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/amirullazmi0/kratify-backend/config"
)

type EmailService struct {
	config *config.SMTPConfig
}

type EmailData struct {
	To      string
	Subject string
	Body    string
}

func NewEmailService(cfg *config.SMTPConfig) *EmailService {
	return &EmailService{config: cfg}
}

// SendEmail sends an email using SMTP
func (s *EmailService) SendEmail(to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", s.config.Email, s.config.Password, s.config.Host)

	// Email headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.config.FromName, s.config.FromEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + htmlBody

	// Send email
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	err := smtp.SendMail(addr, auth, s.config.FromEmail, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendVerificationEmail sends email verification
func (s *EmailService) SendVerificationEmail(to, name, verificationToken, baseURL string) error {
	verificationLink := fmt.Sprintf("%s/auth/verify-email?token=%s", baseURL, verificationToken)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
</head>
<body style="margin: 0; padding: 0; font-family: Arial, sans-serif; background-color: #f4f4f4;">
    <table role="presentation" style="width: 100%; border-collapse: collapse;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 30px; text-align: center; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 8px 8px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: bold;">Welcome to {{.AppName}}! 🎉</h1>
                        </td>
                    </tr>
                    
                    <!-- Body -->
                    <tr>
                        <td style="padding: 40px;">
                            <h2 style="margin: 0 0 20px; color: #333333; font-size: 24px;">Hi {{.Name}},</h2>
                            <p style="margin: 0 0 20px; color: #666666; font-size: 16px; line-height: 1.6;">
                                Thank you for registering! We're excited to have you on board. 
                            </p>
                            <p style="margin: 0 0 30px; color: #666666; font-size: 16px; line-height: 1.6;">
                                To complete your registration and activate your account, please verify your email address by clicking the button below:
                            </p>
                            
                            <!-- Button -->
                            <table role="presentation" style="margin: 0 auto;">
                                <tr>
                                    <td style="border-radius: 6px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
                                        <a href="{{.VerificationLink}}" target="_blank" style="display: inline-block; padding: 16px 48px; color: #ffffff; text-decoration: none; font-size: 16px; font-weight: bold; border-radius: 6px;">
                                            Verify Email Address
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="margin: 30px 0 0; color: #999999; font-size: 14px; line-height: 1.6;">
                                Or copy and paste this link in your browser:
                            </p>
                            <p style="margin: 10px 0 0; color: #667eea; font-size: 14px; word-break: break-all;">
                                {{.VerificationLink}}
                            </p>
                            
                            <div style="margin-top: 40px; padding-top: 30px; border-top: 1px solid #eeeeee;">
                                <p style="margin: 0 0 10px; color: #999999; font-size: 14px;">
                                    <strong>⏱️ Important:</strong> This verification link will expire in <strong>24 hours</strong>.
                                </p>
                                <p style="margin: 0; color: #999999; font-size: 14px;">
                                    If you didn't create this account, please ignore this email.
                                </p>
                            </div>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px; text-align: center; background-color: #f9f9f9; border-radius: 0 0 8px 8px;">
                            <p style="margin: 0 0 10px; color: #999999; font-size: 14px;">
                                Best regards,<br>
                                <strong>{{.AppName}} Team</strong>
                            </p>
                            <p style="margin: 0; color: #cccccc; font-size: 12px;">
                                © 2025 {{.AppName}}. All rights reserved.
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`

	t, err := template.New("verification").Parse(tmpl)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	data := map[string]string{
		"Name":             name,
		"VerificationLink": verificationLink,
		"AppName":          "Kratify Backend",
	}

	if err := t.Execute(&body, data); err != nil {
		return err
	}

	return s.SendEmail(to, "Verify Your Email Address", body.String())
}

// GenerateVerificationToken generates a random verification token
func GenerateVerificationToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
