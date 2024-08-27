package emailservice

import (
	"Loan-Tracker-API/Config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	YOUR_SERVICE_ID  = Config.SERVICE_ID
	YOUR_TEMPLATE_ID = Config.TEMPLATE_ID
	YOUR_PUBLIC_KEY  = Config.PUBLIC_KEY
)

type MailService struct {
	apiToken string
}

func NewMailService() *MailService {
	Config.Envinit()
	return &MailService{}
}
func (s *MailService) SendEmail(toEmail string, subject string, text string, category string) error {
	url := "https://api.emailjs.com/api/v1.0/email/send"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	data := map[string]interface{}{
		"service_id":  Config.SERVICE_ID,
		"template_id": Config.TEMPLATE_ID,
		"user_id":     Config.PUBLIC_KEY,
		"template_params": map[string]string{
			"message":  text,
			"subject":  subject,
			"category": category,
			"to_email": toEmail,
		},
	}

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling json: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var responseBody map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&responseBody)
		if err != nil {
			return fmt.Errorf("error sending email verification email: %s", resp.Status)
		}
		return errors.New("error sending email verification email: " + responseBody["message"].(string))
		// fmt.Errorf("error sending email: %s", responseBody["message"])
	}

	return nil
}

func (s *MailService) SendActivationEmail(email string, activationToken string) error {
	// send activation email
	return s.SendEmail(email, "Verify Email", `Click "`+Config.BASE_URL+`/auth/activate/`+activationToken+`"here to verify email.
	`, "reset")
}

func (s *MailService) SendPasswordResetEmail(email string, resetToken string) error {
	// send password reset email
	return s.SendEmail(email, "Reset Password", `Click "`+Config.BASE_URL+`/auth/password-reset/`+resetToken+`">hereto reset your password.
`, "reset")
}
