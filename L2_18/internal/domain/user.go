package domain

type User struct {
	ID     string
	Events map[string]struct{}
}
