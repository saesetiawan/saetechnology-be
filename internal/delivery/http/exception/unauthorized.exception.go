package exception

type Unauthorized struct {
	Message string
}

func (Unauthorized Unauthorized) Error() string {
	return Unauthorized.Message
}

func NewUnauthorized(err string) Unauthorized {
	return Unauthorized{Message: err}
}
