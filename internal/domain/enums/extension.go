package enums

type Extension string

const (
	EXTENSION_MP4  Extension = "mp4"
	EXTENSION_MOV  Extension = "mov"
	EXTENSION_WEBM Extension = "webm"
)

func (e Extension) String() string {
	return string(e)
}

func (e Extension) IsValid() bool {
	switch e {
	case EXTENSION_MP4, EXTENSION_MOV, EXTENSION_WEBM:
		return true
	default:
		return false
	}
}
