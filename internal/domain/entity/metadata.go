package entity

import (
	"fmt"
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/enums"
	"github.com/gocql/gocql"
)

type MetaData struct {
	UserID gocql.UUID `db:"user_id"`
	ID     gocql.UUID `db:"id"`

	Filename    string `db:"filename"`
	ObjectKey   string `db:"object_key"`
	ManifestKey string `db:"manifest_key"`
	VideoType   string `db:"video_type"`
	Status      string `db:"status"`

	Size            int64   `db:"size"`
	MimeType        string  `db:"mime_type"`
	DurationSeconds float64 `db:"duration_seconds"`

	Width  int `db:"width"`
	Height int `db:"height"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewMetaData(
	userID gocql.UUID,
	filename string,
	objectKey string,
	manifestKey string,
	videoType enums.VideoType,
	size int64,
	mimeType string,
	durationSeconds float64,
	width int,
	height int,
) (*MetaData, error) {
	uuid, err := gocql.RandomUUID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate metadata uuid: %w", err)
	}

	now := time.Now()
	return &MetaData{
		ID:              uuid,
		UserID:          userID,
		Filename:        filename,
		ObjectKey:       objectKey,
		ManifestKey:     manifestKey,
		VideoType:       videoType.String(),
		Status:          enums.StatusNoUploaded.String(),
		Size:            size,
		MimeType:        mimeType,
		DurationSeconds: durationSeconds,
		Width:           width,
		Height:          height,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}
