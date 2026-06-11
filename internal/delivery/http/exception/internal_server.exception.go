package exception

type InternalServiceException struct {
	Message string
}

func (internalServiceException InternalServiceException) Error() string {
	return internalServiceException.Message
}

func NewInternalServiceException(err string) InternalServiceException {
	return InternalServiceException{Message: err}
}
