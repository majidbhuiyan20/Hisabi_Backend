package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"hisabi.com/m/config"
)

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Brevo API দিয়ে OTP Email পাঠাও
// Domain ছাড়াই যেকোনো email এ পাঠাতে পারবে
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func SendOTPEmail(toEmail, username, otp string) error {

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  "Hisabi",
			"email": config.Config.SenderEmail, // তোমার Gmail
		},
		"to": []map[string]string{
			{"email": toEmail},
		},
		"subject":     "Your Hisabi Verification Code",
		"htmlContent": buildOTPEmailBody(username, otp),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to build payload: %w", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.brevo.com/v3/smtp/email",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("api-key", config.Config.BrevoAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("SMTP connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("brevo API error %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
// Beautiful HTML Email Template
// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
func buildOTPEmailBody(username, otp string) string {
	digits := strings.Join(strings.Split(otp, ""), " ")

	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin:0;padding:0;background-color:#f4f4f4;font-family:Arial,sans-serif;">
  <table width="100%%" cellpadding="0" cellspacing="0"
         style="background-color:#f4f4f4;padding:40px 0;">
    <tr>
      <td align="center">
        <table width="500" cellpadding="0" cellspacing="0"
               style="background:#fff;border-radius:12px;
                      box-shadow:0 4px 20px rgba(0,0,0,0.08);overflow:hidden;">

          <!-- Header -->
          <tr>
            <td style="background:linear-gradient(135deg,#667eea,#764ba2);
                       padding:40px;text-align:center;">
              <h1 style="color:#fff;margin:0;font-size:28px;">Hisabi</h1>
              <p style="color:rgba(255,255,255,0.85);margin:8px 0 0;font-size:14px;">
                Business Management Platform
              </p>
            </td>
          </tr>

          <!-- Body -->
          <tr>
            <td style="padding:40px;">
              <h2 style="color:#1a1a2e;margin:0 0 12px;">Verify Your Email</h2>
              <p style="color:#555;font-size:15px;line-height:1.6;margin:0 0 30px;">
                Hi <strong>%s</strong> 👋<br>
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
                <p style="color:#856404;font-size:13px;margin:0;">
                  🔒 Never share this code with anyone.
                </p>
              </div>

              <p style="color:#888;font-size:13px;">
                If you didn't create a Hisabi account, ignore this email.
              </p>
            </td>
          </tr>

          <!-- Footer -->
          <tr>
            <td style="background:#f8f8f8;padding:24px;text-align:center;
                       border-top:1px solid #eee;">
              <p style="color:#aaa;font-size:12px;margin:0;">
                © 2025 Hisabi. All rights reserved.
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
