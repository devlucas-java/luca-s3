package enums

import "fmt"

type Resolution string

const (
	RESOLUTION_360P  Resolution = "360p"
	RESOLUTION_480P  Resolution = "480p"
	RESOLUTION_720P  Resolution = "720p"
	RESOLUTION_1080P Resolution = "1080p"
	RESOLUTION_1440P Resolution = "1440p"
	RESOLUTION_4K    Resolution = "4k"
)

func (r Resolution) String() string {
	return string(r)
}

func (r Resolution) Width() int {
	switch r {
	case RESOLUTION_360P:
		return 640
	case RESOLUTION_480P:
		return 854
	case RESOLUTION_720P:
		return 1280
	case RESOLUTION_1080P:
		return 1920
	case RESOLUTION_1440P:
		return 2560
	case RESOLUTION_4K:
		return 3840
	default:
		return 0
	}
}

func (r Resolution) Height() int {
	switch r {
	case RESOLUTION_360P:
		return 360
	case RESOLUTION_480P:
		return 480
	case RESOLUTION_720P:
		return 720
	case RESOLUTION_1080P:
		return 1080
	case RESOLUTION_1440P:
		return 1440
	case RESOLUTION_4K:
		return 2160
	default:
		return 0
	}
}

func (r Resolution) AspectScale() string {
	return fmt.Sprintf("%d:-2", r.Width())
}
