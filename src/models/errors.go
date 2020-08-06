package models

type ErrNotValidRequest string

func (e ErrNotValidRequest) Error() string {
	return string(e)
}

type ErrInternalSrvErr struct{}

func (e ErrInternalSrvErr) Error() string {
	const msg = "internal service error"
	return msg
}
