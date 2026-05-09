package enums

type Role string

const (
	ADMIN Role = "ADMIN"
	USER  Role = "USER"
)

func (r Role) ToString() string {
	return string(r)
}
