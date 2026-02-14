package helpers

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"strconv"
	"time"
	"tubexxi/video-api/config"
	"tubexxi/video-api/internal/dto"

	"go.uber.org/zap"
)

type MailHelper struct {
	userHelper    *UserHelper
	sessionHelper *SessionHelper
	appConfig     *config.AppConfig
	mailConfig    *config.EmailConfig
	logger        *zap.Logger
}

func NewMailHelper(userHelper *UserHelper, sessionHelper *SessionHelper, appConfig *config.AppConfig, mailConfig *config.EmailConfig, logger *zap.Logger) *MailHelper {
	return &MailHelper{
		userHelper:    userHelper,
		sessionHelper: sessionHelper,
		appConfig:     appConfig,
		mailConfig:    mailConfig,
		logger:        logger,
	}
}
func (h *MailHelper) SendEmail(payload *dto.SendMailMetaData, clientOrigin string) error {
	switch payload.Type {
	case dto.ResetPassword:
		return h.SendResetPasswordEmail(payload, clientOrigin)
	case dto.EmailVerification:
		return h.SendVerificationEmail(payload, clientOrigin)
	case dto.RegistrationInfo:
		return h.SendRegistrationInfo(payload, clientOrigin)
	default:
		return fmt.Errorf("unknown email type: %s", payload.Type)
	}
}

func (h *MailHelper) SendVerificationEmail(payload *dto.SendMailMetaData, clientOrigin string) error {
	url := payload.GetURL(clientOrigin)
	if url == "" {
		url = fmt.Sprintf("%s/auth/verify-email?token=%s", clientOrigin, payload.Token)
	}

	username := payload.To
	if payload.User != nil && payload.User.FullName != "" {
		username = payload.User.FullName
	}

	data := struct {
		Username        string
		VerificationURL string
		Year            int
	}{
		Username:        username,
		VerificationURL: url,
		Year:            time.Now().Year(),
	}

	subject := "Verify Your Email Address"
	body, err := h.renderTemplate(emailVerificationTemplate, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return h.sendHTMLEmail(payload.To, subject, body)
}

func (h *MailHelper) SendResetPasswordEmail(payload *dto.SendMailMetaData, clientOrigin string) error {
	url := payload.GetURL(clientOrigin)
	if url == "" {
		url = fmt.Sprintf("%s/auth/reset-password?token=%s", clientOrigin, payload.Token)
	}

	username := payload.To
	if payload.User != nil && payload.User.FullName != "" {
		username = payload.User.FullName
	}

	data := struct {
		Username string
		ResetURL string
		Year     int
	}{
		Username: username,
		ResetURL: url,
		Year:     time.Now().Year(),
	}

	subject := "Reset Your Password"
	body, err := h.renderTemplate(resetPasswordTemplate, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return h.sendHTMLEmail(payload.To, subject, body)
}

func (h *MailHelper) SendRegistrationInfo(payload *dto.SendMailMetaData, clientOrigin string) error {
	username := payload.To
	if payload.User != nil && payload.User.FullName != "" {
		username = payload.User.FullName
	}

	data := struct {
		Username string
		Email    string
		Password string
		LoginURL string
		Year     int
	}{
		Username: username,
		Email:    payload.To,
		Password: payload.Password,
		LoginURL: fmt.Sprintf("%s/auth/login", clientOrigin),
		Year:     time.Now().Year(),
	}

	subject := "Welcome to AGC Forge - Account Information"
	body, err := h.renderTemplate(registrationInfoTemplate, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	return h.sendHTMLEmail(payload.To, subject, body)
}

func (m *MailHelper) SendContactEmail(ctx context.Context, payload *dto.ContactRequest, clientOrigin string) error {

	settingEmail := m.mailConfig

	// Validate required SMTP settings
	if settingEmail.SmtpHost == "" || settingEmail.SmtpPort == "0" || settingEmail.SmtpUsername == "" || settingEmail.SmtpPassword == "" {
		return fmt.Errorf("incomplete SMTP configuration")
	}

	siteName := m.appConfig.AppName
	if siteName == "" {
		siteName = "Simontok"
	}

	subject := fmt.Sprintf("Contact Us - %s", siteName)
	body := m.getContactHTML(siteName, clientOrigin, payload)

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", settingEmail.SmtpFromName, settingEmail.SmtpFromEmail)
	headers["To"] = m.appConfig.AdminEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Setup SMTP Auth
	auth := smtp.PlainAuth("", settingEmail.SmtpUsername, settingEmail.SmtpPassword, settingEmail.SmtpHost)

	// SMTP Server Address
	port, err := strconv.Atoi(settingEmail.SmtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %w", err)
	}
	addr := fmt.Sprintf("%s:%d", settingEmail.SmtpHost, port)

	// If port is 465, use implicit TLS
	if settingEmail.SmtpPort == "465" {
		// Handle TLS config for non-localhost
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			ServerName:         settingEmail.SmtpHost,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return fmt.Errorf("failed to dial TLS: %w", err)
		}

		c, err := smtp.NewClient(conn, settingEmail.SmtpHost)
		if err != nil {
			return fmt.Errorf("failed to create SMTP client: %w", err)
		}
		defer c.Quit()

		if err = c.Auth(auth); err != nil {
			return fmt.Errorf("failed to auth: %w", err)
		}

		if err = c.Mail(settingEmail.SmtpFromEmail); err != nil {
			return fmt.Errorf("failed to set sender: %w", err)
		}

		if err = c.Rcpt(m.appConfig.AdminEmail); err != nil {
			return fmt.Errorf("failed to set recipient: %w", err)
		}

		w, err := c.Data()
		if err != nil {
			return fmt.Errorf("failed to create data writer: %w", err)
		}

		_, err = w.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("failed to write body: %w", err)
		}

		err = w.Close()
		if err != nil {
			return fmt.Errorf("failed to close data writer: %w", err)
		}

		return nil

	} else {
		// Standard smtp.SendMail for port 587 (STARTTLS) or 25 (Plain)
		err := smtp.SendMail(addr, auth, settingEmail.SmtpFromEmail, []string{m.appConfig.AdminEmail}, []byte(message))
		if err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}

	return nil
}

func (h *MailHelper) renderTemplate(tmpl string, data interface{}) (string, error) {
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (h *MailHelper) sendHTMLEmail(to, subject, htmlBody string) error {
	message := []byte("Subject: " + subject + "\r\n" +
		"To: " + to + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n\r\n" +
		htmlBody + "\r\n")

	auth := smtp.PlainAuth(
		"",
		h.mailConfig.SmtpUsername,
		h.mailConfig.SmtpPassword,
		h.mailConfig.SmtpHost,
	)

	addr := fmt.Sprintf("%s:%s", h.mailConfig.SmtpHost, h.mailConfig.SmtpPort)
	err := smtp.SendMail(addr, auth, h.mailConfig.SmtpUsername, []string{to}, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (h *MailHelper) ValidateSMTPConfig() error {
	if h.mailConfig == nil {
		return fmt.Errorf("SMTP configuration is nil")
	}

	requiredFields := []struct {
		value *string
		field string
	}{
		{&h.mailConfig.SmtpHost, "SMTP host"},
		{&h.mailConfig.SmtpUsername, "SMTP username"},
		{&h.mailConfig.SmtpPassword, "SMTP password"},
		{&h.mailConfig.SmtpPort, "SMTP port"},
	}

	for _, field := range requiredFields {
		if field.value == nil || *field.value == "" {
			return fmt.Errorf("%s is not configured", field.field)
		}
	}

	if _, err := strconv.Atoi(h.mailConfig.SmtpPort); err != nil {
		return fmt.Errorf("invalid SMTP port format: %v", err)
	}

	return nil
}

// Email Verification Template
const emailVerificationTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Verification</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f7fa;">
    <table role="presentation" style="width: 100%; border-collapse: collapse; background-color: #f5f7fa;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 30px; text-align: center; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">Verify Your Email</h1>
                        </td>
                    </tr>
                    
                    <!-- Body -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                Hi <strong>{{.Username}}</strong>,
                            </p>
                            <p style="margin: 0 0 30px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                Thank you for signing up! To complete your registration and start using your account, please verify your email address by clicking the button below.
                            </p>
                            
                            <!-- CTA Button -->
                            <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                <tr>
                                    <td align="center" style="padding: 20px 0;">
                                        <a href="{{.VerificationURL}}" style="display: inline-block; padding: 16px 40px; background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 4px 6px rgba(102, 126, 234, 0.4);">
                                            Verify Email Address
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="margin: 30px 0 20px; color: #718096; font-size: 14px; line-height: 1.6;">
                                Or copy and paste this link into your browser:
                            </p>
                            <p style="margin: 0; padding: 15px; background-color: #f7fafc; border-radius: 6px; word-break: break-all; font-size: 13px; color: #4a5568; border-left: 4px solid #667eea;">
                                {{.VerificationURL}}
                            </p>
                            
                            <p style="margin: 30px 0 0; color: #718096; font-size: 14px; line-height: 1.6;">
                                This link will expire in 24 hours. If you didn't create an account, please ignore this email.
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px; background-color: #f7fafc; border-radius: 0 0 12px 12px; text-align: center;">
                            <p style="margin: 0 0 10px; color: #a0aec0; font-size: 13px;">
                                © {{.Year}} AGC Forge. All rights reserved.
                            </p>
                            <p style="margin: 0; color: #a0aec0; font-size: 13px;">
                                Need help? Contact us at support@socialforge.io
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

// Reset Password Template
const resetPasswordTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Password</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f7fa;">
    <table role="presentation" style="width: 100%; border-collapse: collapse; background-color: #f5f7fa;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 30px; text-align: center; background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 28px; font-weight: 600;">Reset Your Password</h1>
                        </td>
                    </tr>
                    
                    <!-- Body -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                Hi <strong>{{.Username}}</strong>,
                            </p>
                            <p style="margin: 0 0 30px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                We received a request to reset your password. Click the button below to create a new password for your account.
                            </p>
                            
                            <!-- CTA Button -->
                            <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                <tr>
                                    <td align="center" style="padding: 20px 0;">
                                        <a href="{{.ResetURL}}" style="display: inline-block; padding: 16px 40px; background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 4px 6px rgba(245, 87, 108, 0.4);">
                                            Reset Password
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <p style="margin: 30px 0 20px; color: #718096; font-size: 14px; line-height: 1.6;">
                                Or copy and paste this link into your browser:
                            </p>
                            <p style="margin: 0; padding: 15px; background-color: #f7fafc; border-radius: 6px; word-break: break-all; font-size: 13px; color: #4a5568; border-left: 4px solid #f5576c;">
                                {{.ResetURL}}
                            </p>
                            
                            <!-- Security Notice -->
                            <div style="margin-top: 30px; padding: 20px; background-color: #fff5f5; border-left: 4px solid #f5576c; border-radius: 6px;">
                                <p style="margin: 0 0 10px; color: #c53030; font-weight: 600; font-size: 14px;">
                                    🔒 Security Notice
                                </p>
                                <p style="margin: 0; color: #742a2a; font-size: 14px; line-height: 1.6;">
                                    This link will expire in 1 hour. If you didn't request a password reset, please ignore this email and your password will remain unchanged.
                                </p>
                            </div>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px; background-color: #f7fafc; border-radius: 0 0 12px 12px; text-align: center;">
                            <p style="margin: 0 0 10px; color: #a0aec0; font-size: 13px;">
                                © {{.Year}} AGC Forge. All rights reserved.
                            </p>
                            <p style="margin: 0; color: #a0aec0; font-size: 13px;">
                                Need help? Contact us at support@socialforge.io
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

// Registration Info Template
const registrationInfoTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to AGC Forge</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f7fa;">
    <table role="presentation" style="width: 100%; border-collapse: collapse; background-color: #f5f7fa;">
        <tr>
            <td align="center" style="padding: 40px 0;">
                <table role="presentation" style="width: 600px; border-collapse: collapse; background-color: #ffffff; border-radius: 12px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 30px; text-align: center; background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); border-radius: 12px 12px 0 0;">
                            <h1 style="margin: 0 0 10px; color: #ffffff; font-size: 32px; font-weight: 600;">🎉 Welcome!</h1>
                            <p style="margin: 0; color: #ffffff; font-size: 16px; opacity: 0.95;">Your account has been created</p>
                        </td>
                    </tr>
                    
                    <!-- Body -->
                    <tr>
                        <td style="padding: 40px;">
                            <p style="margin: 0 0 20px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                Hi <strong>{{.Username}}</strong>,
                            </p>
                            <p style="margin: 0 0 30px; color: #4a5568; font-size: 16px; line-height: 1.6;">
                                Welcome to <strong>AGC Forge</strong>! Your account has been successfully created. Below are your login credentials:
                            </p>
                            
                            <!-- Credentials Box -->
                            <div style="background: linear-gradient(135deg, #f6f8fb 0%, #e9ecef 100%); border-radius: 8px; padding: 25px; margin-bottom: 30px; border: 1px solid #dee2e6;">
                                <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                    <tr>
                                        <td style="padding: 12px 0; border-bottom: 1px solid #cbd5e0;">
                                            <p style="margin: 0; color: #718096; font-size: 13px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px;">Email Address</p>
                                            <p style="margin: 8px 0 0; color: #2d3748; font-size: 16px; font-weight: 500;">{{.Email}}</p>
                                        </td>
                                    </tr>
                                    <tr>
                                        <td style="padding: 12px 0;">
                                            <p style="margin: 0; color: #718096; font-size: 13px; font-weight: 600; text-transform: uppercase; letter-spacing: 0.5px;">Temporary Password</p>
                                            <p style="margin: 8px 0 0; color: #2d3748; font-size: 16px; font-weight: 500; font-family: 'Courier New', monospace; background-color: #ffffff; padding: 10px; border-radius: 4px; display: inline-block;">{{.Password}}</p>
                                        </td>
                                    </tr>
                                </table>
                            </div>
                            
                            <!-- CTA Button -->
                            <table role="presentation" style="width: 100%; border-collapse: collapse;">
                                <tr>
                                    <td align="center" style="padding: 10px 0;">
                                        <a href="{{.LoginURL}}" style="display: inline-block; padding: 16px 40px; background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); color: #ffffff; text-decoration: none; border-radius: 8px; font-weight: 600; font-size: 16px; box-shadow: 0 4px 6px rgba(79, 172, 254, 0.4);">
                                            Login to Your Account
                                        </a>
                                    </td>
                                </tr>
                            </table>
                            
                            <!-- Security Notice -->
                            <div style="margin-top: 30px; padding: 20px; background-color: #fffaf0; border-left: 4px solid #ed8936; border-radius: 6px;">
                                <p style="margin: 0 0 10px; color: #c05621; font-weight: 600; font-size: 14px;">
                                    🔐 Important Security Notice
                                </p>
                                <p style="margin: 0; color: #7c2d12; font-size: 14px; line-height: 1.6;">
                                    This is a temporary password. For your security, please change it immediately after your first login by going to Account Settings > Security.
                                </p>
                            </div>
                            
                            <p style="margin: 30px 0 0; color: #718096; font-size: 14px; line-height: 1.6;">
                                If you have any questions or need assistance, our support team is here to help!
                            </p>
                        </td>
                    </tr>
                    
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px; background-color: #f7fafc; border-radius: 0 0 12px 12px; text-align: center;">
                            <p style="margin: 0 0 10px; color: #a0aec0; font-size: 13px;">
                                © {{.Year}} AGC Forge. All rights reserved.
                            </p>
                            <p style="margin: 0; color: #a0aec0; font-size: 13px;">
                                Need help? Contact us at support@socialforge.io
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>`

func (m *MailHelper) getContactHTML(siteName, siteURL string, payload *dto.ContactRequest) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; background-color: #f9f9f9; }
        .header { background-color: #007bff; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { background-color: white; padding: 30px; border-radius: 0 0 5px 5px; box-shadow: 0 2px 5px rgba(0,0,0,0.1); }
        .button { display: inline-block; padding: 12px 24px; background-color: #007bff; color: white; text-decoration: none; border-radius: 4px; margin-top: 20px; font-weight: bold; }
        .footer { text-align: center; margin-top: 20px; font-size: 12px; color: #666; }
        p { margin-bottom: 15px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1 style="margin:0;">%s</h1>
        </div>
        <div class="content">
            <h2>New Contact Request</h2>
            <p>You received a new contact request from %s.</p>
            <p>Email: %s</p>
            <p>Name: %s</p>
            <p>Message: %s</p>
            <p style="margin-top: 30px; font-size: 14px;">Or copy and paste this link into your browser:</p>
            <p style="font-size: 13px; color: #007bff; word-break: break-all;"><a href="%s">%s</a></p>
            <p>This link will expire in 1 hour.</p>
        </div>
        <div class="footer">
            <p>&copy; %d %s. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`, siteName, payload.Name, payload.Email, payload.Name, payload.Message, siteURL, siteURL, 2025, siteName)
}
