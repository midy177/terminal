package privilege

type Privilege interface {
	IsAdmin() bool
	Elevate() error
}
type privilege struct {
}

func New() Privilege {
	return &privilege{}
}
