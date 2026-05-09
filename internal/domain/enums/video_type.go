package enums

type VideoType string

const (
	VideoTypeNormal VideoType = "normal"
	VideoType360    VideoType = "360"
)

func (vt VideoType) ToString() string {
	return string(vt)
}
