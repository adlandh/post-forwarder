package domain

type RequestIDKey string

func (r RequestIDKey) String() string {
	return string(r)
}
