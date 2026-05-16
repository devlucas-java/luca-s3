package enums

// VideoType define se o vídeo é normal ou 360°.
type VideoType string

const (
	VideoTypeNormal VideoType = "normal" // Vídeo plano tradicional
	VideoType360    VideoType = "360"    // Vídeo esférico/equirectangular
)

func (v VideoType) String() string {
	return string(v)
}
