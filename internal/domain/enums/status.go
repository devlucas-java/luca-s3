package enums

type Status string

const (
	StatusFailed     Status = "FAILED"
	StatusPending    Status = "PENDING"
	StatusUploaded   Status = "UPLOADED"
	StatusNoUploaded Status = "NO_UPLOADED"
)

func (s Status) String() string {
	return string(s)
}
