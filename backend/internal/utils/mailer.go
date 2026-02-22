package utils

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/resend/resend-go/v3"
)

type EmailService struct {
	client      *resend.Client
	fromAddress string
	frontendURL string
}

func NewEmailService() *EmailService {
	apiKey := os.Getenv("RESEND_API_KEY")
	if apiKey == "" {
		Log.Fatal("RESEND_API_KEY not set")
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		Log.Fatal("FRONTEND_URL not set")
	}

	fromAddress := os.Getenv("RESEND_FROM_EMAIL")
	if fromAddress == "" {
		fromAddress = "FrameRate <onboarding@resend.dev>"
	}

	return &EmailService{
		client:      resend.NewClient(apiKey),
		fromAddress: fromAddress,
		frontendURL: frontendURL,
	}
}

func (s *EmailService) SendVerificationEmail(to, username, token string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", s.frontendURL, token)

	html := fmt.Sprintf(`
        <!DOCTYPE html>
        <html>
        <head>
            <meta charset="UTF-8">
            <style>
                body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
                .container { max-width: 600px; margin: 0 auto; padding: 20px; }
                .header { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; padding: 30px; text-align: center; border-radius: 10px 10px 0 0; }
                .content { background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; }
                .button { display: inline-block; background: #667eea; color: white; padding: 15px 30px; text-decoration: none; border-radius: 5px; margin: 20px 0; }
                .footer { text-align: center; margin-top: 20px; color: #888; font-size: 12px; }
            </style>
        </head>
        <body>
            <div class="container">
                <div class="header">
                    <h1>🎬 Welcome to FrameRate!</h1>
                </div>
                <div class="content">
                    <p>Hi <strong>%s</strong>,</p>
                    <p>Thanks for signing up! Please verify your email to start tracking your movies.</p>
                    <p style="text-align: center;">
                        <a href="%s" class="button">Verify Email</a>
                    </p>
                    <p>Or copy this link:</p>
                    <p style="background: white; padding: 10px; border-left: 3px solid #667eea; word-break: break-all;">
                        %s
                    </p>
                    <p><small>This link expires in 24 hours.</small></p>
                </div>
                <div class="footer">
                    <p>If you didn't create an account, you can safely ignore this email.</p>
                </div>
            </div>
        </body>
        </html>
    `, username, verifyURL, verifyURL)

	params := &resend.SendEmailRequest{
		From:    s.fromAddress,
		To:      []string{to},
		Subject: "Verify your FrameRate account",
		Html:    html,
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		Log.Error("Failed to send email",
			zap.String("to", to),
			zap.Error(err),
		)
		return fmt.Errorf("failed to send email: %w", err)
	}

	Log.Info("Verification email sent", zap.String("to", to))
	return nil
}
