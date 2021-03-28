package shortener

type RepositoryRedirecter interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
