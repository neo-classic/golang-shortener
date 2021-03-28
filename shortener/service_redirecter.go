package shortener

type ServiceRedirecter interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
