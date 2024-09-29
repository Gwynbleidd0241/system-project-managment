package rest

import (
	"encoding/json"
	"net/http"
	"notification-service/mailer"
)

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func HandleSendEmail(mailSender *mailer.MailSender) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req EmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err := mailSender.SendEmail(req.To, req.Subject, req.Body, []string{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Email sent successfully!")
	}
}

func StartServer(addr string) {
	mailSender := mailer.NewMailSender("smtp.yandex.ru", 465, "your-email@yandex.ru", "your-password")

	http.HandleFunc("/send-email", HandleSendEmail(mailSender))
	http.ListenAndServe(addr, nil)
}
