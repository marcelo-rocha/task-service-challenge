package entities

type UserKind string

const (
	Technician UserKind = "technician"
	Manager    UserKind = "manager"
)

type User struct {
	Id        int64
	Login     string
	Name      string
	Kind      UserKind
	Active    bool
	ManagerId int64
}

func ToUserKind(str string) UserKind {
	if str == "manager" {
		return Manager
	}
	return Technician
}
