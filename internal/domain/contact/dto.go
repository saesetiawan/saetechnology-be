package contact

type CreateContactMessageDto struct {
	Name          string `json:"name" validate:"required,max=160"`
	Email         string `json:"email" validate:"required,email,max=180"`
	Phone         string `json:"phone" validate:"max=60"`
	Company       string `json:"company" validate:"max=160"`
	Subject       string `json:"subject" validate:"required,max=220"`
	Message       string `json:"message" validate:"required"`
	CaptchaID     string `json:"captcha_id" validate:"required"`
	CaptchaAnswer string `json:"captcha_answer" validate:"required"`
	Website       string `json:"website"`
}

type UpdateContactStatusDto struct {
	ID     string `json:"id"`
	Status string `json:"status" validate:"required,oneof=new read replied archived"`
}

type ListContactMessageQuery struct {
	Page      int
	Limit     int
	Search    string
	Status    string
	OrderBy   string
	OrderType string
}

type ContactCaptcha struct {
	ID       string `json:"id"`
	Question string `json:"question"`
}
