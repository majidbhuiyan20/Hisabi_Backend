package services

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"hisabi.com/m/config"
)

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// OTP Email পাঠাও
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func SendOTPEmail(toEmail, username, otp string) error {

	from := config.Config.SMTPEmail
	password := config.Config.SMTPPassword
	host := config.Config.SMTPHost
	port := config.Config.SMTPPort

	// ── Email Content বানাও ───────────────────────────────
	subject := "Your Hisabi Verification Code"
	body := buildOTPEmailBody(username, otp)

	// ── MIME format এ email বানাও ─────────────────────────
	// Plain text + HTML দুইটাই দিচ্ছি
	// Email client যেটা support করে সেটা দেখাবে
	message := buildMIMEMessage(from, toEmail, subject, body)

	// ── SMTP Auth ─────────────────────────────────────────
	auth := smtp.PlainAuth("", from, password, host)

	// ── TLS Connection (secure) ───────────────────────────
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	// Port 587 এ STARTTLS ব্যবহার হয়
	conn, err := smtp.Dial(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return fmt.Errorf("SMTP connection failed: %w", err)
	}
	defer conn.Close()

	// STARTTLS শুরু করো
	if err = conn.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("STARTTLS failed: %w", err)
	}

	// Authenticate করো
	if err = conn.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}

	// From address set করো
	if err = conn.Mail(from); err != nil {
		return fmt.Errorf("SMTP mail from failed: %w", err)
	}

	// To address set করো
	if err = conn.Rcpt(toEmail); err != nil {
		return fmt.Errorf("SMTP rcpt failed: %w", err)
	}

	// Email body লিখো
	wc, err := conn.Data()
	if err != nil {
		return fmt.Errorf("SMTP data failed: %w", err)
	}
	defer wc.Close()

	_, err = fmt.Fprint(wc, message)
	return err
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Beautiful HTML Email Template
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func buildOTPEmailBody(username, otp string) string {

	// OTP কে আলাদা করে দেখাই: 1 2 3 4 5 6
	digits := strings.Join(strings.Split(otp, ""), " ")

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,sans-serif;">

  <table width="100%%" cellpadding="0" cellspacing="0" style="background-color:#f4f4f4;padding:40px 0;">
    <tr>
      <td align="center">
        <table width="500" cellpadding="0" cellspacing="0"
               style="background-color:#ffffff;border-radius:12px;
                      box-shadow:0 4px 20px rgba(0,0,0,0.08);overflow:hidden;">

          <!-- Header -->
          <tr>
            <td style="background:linear-gradient(135deg,#667eea 0%%,#764ba2 100%%);
                       padding:40px;text-align:center;">
              <h1 style="color:#ffffff;margin:0;font-size:28px;font-weight:700;
                         letter-spacing:1px;">Hisabi</h1>
              <p style="color:rgba(255,255,255,0.85);margin:8px 0 0;font-size:14px;">
                Business Management Platform
              </p>
            </td>
          </tr>

          <!-- Body -->
          <tr>
            <td style="padding:40px 40px 30px;">

              <h2 style="color:#1a1a2e;margin:0 0 12px;font-size:22px;">
                Verify Your Email
              </h2>
              <p style="color:#555;font-size:15px;line-height:1.6;margin:0 0 30px;">
                Hi <strong>%s</strong>, welcome to Hisabi! 👋<br>
                Use the code below to verify your email address.
              </p>

              <!-- OTP Box -->
              <div style="background:#f8f7ff;border:2px dashed #667eea;
                          border-radius:12px;padding:30px;text-align:center;
                          margin-bottom:30px;">
                <p style="color:#888;font-size:12px;margin:0 0 12px;
                           text-transform:uppercase;letter-spacing:2px;">
                  Your Verification Code
                </p>
                <div style="font-size:42px;font-weight:800;letter-spacing:14px;
                            color:#667eea;font-family:monospace;">
                  %s
                </div>
                <p style="color:#e74c3c;font-size:13px;margin:16px 0 0;">
                  ⏱ Expires in <strong>10 minutes</strong>
                </p>
              </div>

              <!-- Warning -->
              <div style="background:#fff8e1;border-left:4px solid #ffc107;
                          border-radius:4px;padding:16px;margin-bottom:24px;">
                <p style="color:#856404;font-size:13px;margin:0;line-height:1.5;">
                  🔒 <strong>Security tip:</strong> Never share this code with anyone.
                  Hisabi will never ask for your OTP via phone or email.
                </p>
              </div>

              <p style="color:#888;font-size:13px;line-height:1.6;margin:0;">
                If you didn't create a Hisabi account, you can safely ignore this email.
              </p>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="background:#f8f8f8;padding:24px 40px;text-align:center;
                       border-top:1px solid #eee;">
              <p style="color:#aaa;font-size:12px;margin:0;">
                © 2025 Hisabi. All rights reserved.<br>
                This is an automated email, please do not reply.
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

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// MIME Message format বানাও
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func buildMIMEMessage(from, to, subject, htmlBody string) string {
	return fmt.Sprintf(
		"From: Hisabi <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		from, to, subject, htmlBody,
	)
}
