package email

type Sender struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Recipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type MessageVersion struct {
	To          []Recipient `json:"to"`
	Subject     string      `json:"subject,omitempty"`
	HTMLContent string      `json:"htmlContent,omitempty"`
}

type EmailRequest struct {
	Sender          Sender           `json:"sender"`
	Subject         string           `json:"subject"`
	HTMLContent     string           `json:"htmlContent"`
	MessageVersions []MessageVersion `json:"messageVersions"`
}

type EmailSender interface {
	SendEmail(payload EmailRequest) error
}
