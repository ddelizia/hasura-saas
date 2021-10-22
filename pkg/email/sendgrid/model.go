package sendgrid

type SendgridEmailData struct {
	From             *From               `json:"from"`
	Personalizations []*Personalizations `json:"personalizations"`
	TemplateID       string              `json:"template_id"`
	MailSettings     *MailSettings       `json:"mail_settings"`
}
type From struct {
	Email string `json:"email"`
}
type To struct {
	Email string `json:"email"`
}
type MailSettings struct {
	SandboxMode *SandboxMode `json:"sandbox_mode"`
}
type SandboxMode struct {
	Enable bool `json:"enable"`
}
type Personalizations struct {
	To                  []*To       `json:"to"`
	DynamicTemplateData interface{} `json:"dynamic_template_data"`
}
