package entity

import (
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/devlucas-java/luca-s3/pkg/id"
)

type MetaData struct {
	UserID id.UUID `db:"user_id"`
	ID     id.UUID `db:"id"`

	Filename string `db:"filename"`

	ObjectKey   string `db:"object_key"`
	ManifestKey string `db:"manifest_key"`

	VideoType string `db:"video_type"`

	Size            int64   `db:"size"`
	MimeType        string  `db:"mime_type"`
	DurationSeconds float64 `db:"duration_seconds"`

	Width  int `db:"width"`
	Height int `db:"height"`

	CreatedAt time.Time `db:"created_at"`
}

func NewMetaData(
	userID id.UUID,
	filename string,
	objectKey string,
	manifestKey string,
	videoType enums.VideoType,
	size int64,
	mimeType string,
	durationSeconds float64,
	width int,
	height int,
) *MetaData {

	return &MetaData{
		ID:              id.NewUUID(),
		UserID:          userID,
		Filename:        filename,
		ObjectKey:       objectKey,
		ManifestKey:     manifestKey,
		VideoType:       videoType.ToString(),
		Size:            size,
		MimeType:        mimeType,
		DurationSeconds: durationSeconds,
		Width:           width,
		Height:          height,
		CreatedAt:       time.Now(),
	}
}
