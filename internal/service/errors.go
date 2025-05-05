package service

type ServiceError struct {
	Message string
	Code    uint
}

func (e ServiceError) Error() string {
	return e.Message
}

func NewError(message string, code uint) ServiceError {
	return ServiceError{
		Message: message,
		Code:    code,
	}
}
