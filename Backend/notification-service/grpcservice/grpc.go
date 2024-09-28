package grpcservice

import (
	"context"
	"log"
	"notification-service/mailer"
)

// MailService реализует сгенерированный интерфейс
type MailService struct {
	mailer.UnimplementedMailerServer // Встраивание необработанного сервиса
}

func (s *MailService) SendEmail(ctx context.Context, req *mailer.SendEmailRequest) (*mailer.SendEmailResponse, error) {
	mailSender := mailer.NewMailSender(
		"smtp.yandex.ru",
		465,
		"system-project-management@yandex.ru",
		"pjolfmygivbginka",
	)

	err := mailSender.SendEmail(req.To, req.Subject, req.Body, req.Attachments)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return nil, err
	}

	return &mailer.SendEmailResponse{Message: "Email sent successfully!"}, nil
}
