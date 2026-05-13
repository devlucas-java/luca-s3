package enums

type Role string

const (
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

func (r Role) String() string {
	return string(r)
}
