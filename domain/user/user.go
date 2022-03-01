package user

type Kind string

const (
	Technician Kind = "technician"
	Manager    Kind = "manager"
)

type User struct {
	Id   string
	Name string
	Kind Kind
}
