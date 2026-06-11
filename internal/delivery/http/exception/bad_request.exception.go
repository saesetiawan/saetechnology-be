package exception

type BadRequestException struct {
	Message string
}

func (badRequestException BadRequestException) Error() string {
	return badRequestException.Message
}

func NewBadRequestException(err string) BadRequestException {
	return BadRequestException{Message: err}
}
