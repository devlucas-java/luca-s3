package enums

// Extension define os formatos de entrada suportados.
// Todos serão convertidos para HLS (.m3u8 + .ts) e DASH (.mpd + .m4s).
type Extension string

const (
	EXTENSION_MP4  Extension = "mp4"  // Container mais comum
	EXTENSION_MOV  Extension = "mov"  // Apple QuickTime
	EXTENSION_WEBM Extension = "webm" // VP9/VP8
)

func (e Extension) String() string {
	return string(e)
}

// IsValid verifica se a extensão é suportada.
func (e Extension) IsValid() bool {
	switch e {
	case EXTENSION_MP4, EXTENSION_MOV, EXTENSION_WEBM:
		return true
	default:
		return false
	}
}
