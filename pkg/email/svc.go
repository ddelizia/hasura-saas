package email

type Service interface {
	SendEmail(emails []string, subject, template string, data interface{}) error
}
