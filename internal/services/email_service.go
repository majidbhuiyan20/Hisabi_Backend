package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"hisabi.com/m/config"
)

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// OTP email send using Gmail SMTP
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func SendOTPEmail(toEmail, username, otp string) error {

	from := config.Config.SMTPEmail
	password := config.Config.SMTPPassword
	host := config.Config.SMTPHost
	port := config.Config.SMTPPort

	auth := smtp.PlainAuth("", from, password, host)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	conn, err := smtp.Dial(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return fmt.Errorf("SMTP connection failed: %w", err)
	}
	defer conn.Close()

	if err = conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("STARTTLS failed: %w", err)
	}

	if err = conn.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	if err = conn.Mail(from); err != nil {
		return fmt.Errorf("SMTP mail from failed: %w", err)
	}

	if err = conn.Rcpt(toEmail); err != nil {
		return fmt.Errorf("SMTP rcpt failed: %w", err)
	}

	wc, err := conn.Data()
	if err != nil {
		return fmt.Errorf("SMTP data failed: %w", err)
	}
	defer wc.Close()

	message := buildMessage(from, toEmail, username, otp)
	_, err = fmt.Fprint(wc, message)
	return err
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Create Email Message
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func buildMessage(from, to, username, otp string) string {
	return fmt.Sprintf(
		"From: Hisabi <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: Your Hisabi Verification Code\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		from, to, buildOTPEmailBody(username, otp),
	)
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Beautiful HTML Email — OTP
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func buildOTPEmailBody(username, otp string) string {
	// OTP digits Seperate: "847291" → "8 4 7 2 9 1"
	digits := strings.Join(strings.Split(otp, ""), "  ")

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1.0">
</head>
<body style="margin:0;padding:0;background:#f0f2f5;font-family:'Segoe UI',Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0" style="padding:50px 0;">
    <tr>
      <td align="center">
        <table width="480" cellpadding="0" cellspacing="0"
               style="background:#fff;border-radius:16px;
                      box-shadow:0 8px 30px rgba(0,0,0,0.10);
                      overflow:hidden;">

          <!-- ── Header ── -->
          <tr>
            <td style="background:linear-gradient(135deg,#4f46e5 0%%,#7c3aed 100%%);
                       padding:36px 40px;text-align:center;">
              <div style="display:inline-block;background:rgba(255,255,255,0.15);
                          border-radius:12px;padding:8px 20px;margin-bottom:12px;">
                <span style="color:#fff;font-size:22px;font-weight:800;
                             letter-spacing:2px;">HISABI</span>
              </div>
              <p style="color:rgba(255,255,255,0.80);margin:0;font-size:13px;">
                Business Management Platform
              </p>
            </td>
          </tr>

          <!-- ── Body ── -->
          <tr>
            <td style="padding:40px 40px 30px;">

              <h2 style="color:#1e1b4b;margin:0 0 10px;font-size:22px;font-weight:700;">
                Verify Your Email Address
              </h2>
              <p style="color:#6b7280;font-size:15px;line-height:1.7;margin:0 0 28px;">
                Hi <strong style="color:#1e1b4b;">%s</strong> 👋 — Welcome to Hisabi!<br>
                Enter the code below to complete your registration.
              </p>

              <!-- ── OTP One Line Box ── -->
              <div style="background:linear-gradient(135deg,#eef2ff,#f5f3ff);
                          border:2px solid #c7d2fe;border-radius:14px;
                          padding:28px 20px;text-align:center;margin-bottom:28px;">
                <p style="color:#6366f1;font-size:11px;font-weight:700;
                           letter-spacing:3px;text-transform:uppercase;margin:0 0 14px;">
                  One-Time Verification Code
                </p>
                <div style="font-size:46px;font-weight:900;
                            letter-spacing:10px;color:#4f46e5;
                            font-family:'Courier New',monospace;
                            background:#fff;border-radius:10px;
                            padding:16px 24px;display:inline-block;
                            box-shadow:0 2px 12px rgba(79,70,229,0.15);">
                  %s
                </div>
                <p style="color:#ef4444;font-size:13px;
                           font-weight:600;margin:16px 0 0;">
                  ⏱ &nbsp;Expires in <strong>10 minutes</strong>
                </p>
              </div>

              <!-- ── Warning ── -->
              <div style="background:#fffbeb;border:1px solid #fde68a;
                          border-radius:10px;padding:14px 18px;margin-bottom:24px;">
                <p style="color:#92400e;font-size:13px;margin:0;line-height:1.6;">
                  🔒 &nbsp;<strong>Security Notice:</strong> Never share this code with anyone.
                  Hisabi will never ask for your OTP via phone or chat.
                </p>
              </div>

              <p style="color:#9ca3af;font-size:13px;margin:0;">
                Didn't create a Hisabi account? You can safely ignore this email.
              </p>

            </td>
          </tr>

          <!-- ── Footer ── -->
          <tr>
            <td style="background:#f9fafb;padding:20px 40px;
                       text-align:center;border-top:1px solid #f3f4f6;">
              <p style="color:#d1d5db;font-size:12px;margin:0;">
                © 2025 Hisabi · All rights reserved<br>
                This is an automated message — please do not reply.
              </p>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>
</body>
</html>`, username, digits)
}
