package enums

type Status string

const (
	FAILED      Status = "FAILED"
	PENDING     Status = "PENDING"
	UPLOADED    Status = "UPLOADED"
	NO_UPLOADED Status = "NO_UPLOADED"
)

func (s Status) ToString() string {
	return string(s)
}
