package database

import (
	"fmt"
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/repository"
	"github.com/devlucas-java/luca-s3/pkg/pagination"
	"github.com/gocql/gocql"
)

type MetaDataDB struct {
	session *gocql.Session
}

func NewMetaDataDB(session *gocql.Session) repository.MetaDataRepository {
	return &MetaDataDB{session: session}
}

func (m *MetaDataDB) Create(metadata *entity.MetaData) (*entity.MetaData, error) {
	if metadata == nil {
		return nil, fmt.Errorf("metadata is nil")
	}

	if metadata.ID == (gocql.UUID{}) {
		uuid, err := gocql.RandomUUID()
		if err != nil {
			return nil, fmt.Errorf("failed to generate metadata uuid: %w", err)
		}
		metadata.ID = uuid
	}

	now := time.Now()
	metadata.CreatedAt = now
	metadata.UpdatedAt = now

	const query = `
		INSERT INTO metadata (
			id, user_id, filename, object_key, manifest_key,
			video_type, status, size, mime_type, duration_seconds,
			width, height, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	if err := m.session.Query(query,
		metadata.ID,
		metadata.UserID,
		metadata.Filename,
		metadata.ObjectKey,
		metadata.ManifestKey,
		metadata.VideoType,
		metadata.Status,
		metadata.Size,
		metadata.MimeType,
		metadata.DurationSeconds,
		metadata.Width,
		metadata.Height,
		metadata.CreatedAt,
		metadata.UpdatedAt,
	).Exec(); err != nil {
		return nil, fmt.Errorf("failed to create metadata: %w", err)
	}

	return metadata, nil
}

func (m *MetaDataDB) Updates(metadata *entity.MetaData) (*entity.MetaData, error) {
	if metadata == nil {
		return nil, fmt.Errorf("metadata is nil")
	}

	if metadata.ID == (gocql.UUID{}) {
		return nil, fmt.Errorf("metadata id is required for updates")
	}

	current, err := m.FindByID(metadata.UserID, metadata.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}
	if current == nil {
		return nil, fmt.Errorf("metadata not found")
	}

	if metadata.Filename != "" {
		current.Filename = metadata.Filename
	}
	if metadata.ObjectKey != "" {
		current.ObjectKey = metadata.ObjectKey
	}
	if metadata.ManifestKey != "" {
		current.ManifestKey = metadata.ManifestKey
	}
	if metadata.VideoType != "" {
		current.VideoType = metadata.VideoType
	}
	if metadata.Status != "" {
		current.Status = metadata.Status
	}
	if metadata.Size != 0 {
		current.Size = metadata.Size
	}
	if metadata.MimeType != "" {
		current.MimeType = metadata.MimeType
	}
	if metadata.DurationSeconds != 0 {
		current.DurationSeconds = metadata.DurationSeconds
	}
	if metadata.Width != 0 {
		current.Width = metadata.Width
	}
	if metadata.Height != 0 {
		current.Height = metadata.Height
	}
	current.UpdatedAt = time.Now()

	const query = `
		UPDATE metadata
		SET filename = ?, object_key = ?, manifest_key = ?, video_type = ?, status = ?,
		    size = ?, mime_type = ?, duration_seconds = ?, width = ?, height = ?, updated_at = ?
		WHERE user_id = ? AND id = ?`

	if err := m.session.Query(query,
		current.Filename,
		current.ObjectKey,
		current.ManifestKey,
		current.VideoType,
		current.Status,
		current.Size,
		current.MimeType,
		current.DurationSeconds,
		current.Width,
		current.Height,
		current.UpdatedAt,
		current.UserID,
		current.ID,
	).Exec(); err != nil {
		return nil, fmt.Errorf("failed to update metadata: %w", err)
	}

	return current, nil
}

func (m *MetaDataDB) FindByID(userID gocql.UUID, id gocql.UUID) (*entity.MetaData, error) {
	const query = `
		SELECT id, user_id, filename, object_key, manifest_key,
		       video_type, status, size, mime_type, duration_seconds,
		       width, height, created_at, updated_at
		FROM metadata
		WHERE user_id = ? AND id = ?`

	var meta entity.MetaData
	if err := m.session.Query(query, userID, id).Consistency(gocql.One).Scan(
		&meta.ID,
		&meta.UserID,
		&meta.Filename,
		&meta.ObjectKey,
		&meta.ManifestKey,
		&meta.VideoType,
		&meta.Status,
		&meta.Size,
		&meta.MimeType,
		&meta.DurationSeconds,
		&meta.Width,
		&meta.Height,
		&meta.CreatedAt,
		&meta.UpdatedAt,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find metadata by id: %w", err)
	}

	return &meta, nil
}

func (m *MetaDataDB) FindAllByUserID(userID gocql.UUID, page pagination.Page) (*pagination.PagedResult[*entity.MetaData], error) {
	const query = `
		SELECT id, user_id, filename, object_key, manifest_key,
		       video_type, status, size, mime_type, duration_seconds,
		       width, height, created_at, updated_at
		FROM metadata
		WHERE user_id = ?`

	iter := m.session.Query(query, userID).
		Consistency(gocql.One).
		PageSize(page.Size).
		PageState(page.PageState).
		Iter()

	var items []*entity.MetaData
	for {
		var meta entity.MetaData
		if !iter.Scan(
			&meta.ID,
			&meta.UserID,
			&meta.Filename,
			&meta.ObjectKey,
			&meta.ManifestKey,
			&meta.VideoType,
			&meta.Status,
			&meta.Size,
			&meta.MimeType,
			&meta.DurationSeconds,
			&meta.Width,
			&meta.Height,
			&meta.CreatedAt,
			&meta.UpdatedAt,
		) {
			break
		}
		cp := meta
		items = append(items, &cp)
	}

	nextPageState := iter.PageState()
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to list metadata by user id: %w", err)
	}

	return &pagination.PagedResult[*entity.MetaData]{Items: items, NextPageState: nextPageState}, nil
}

func (m *MetaDataDB) FindByObjectKey(objectKey string, page pagination.Page) (*pagination.PagedResult[*entity.MetaData], error) {
	const query = `
		SELECT id, user_id, filename, object_key, manifest_key,
		       video_type, status, size, mime_type, duration_seconds,
		       width, height, created_at, updated_at
		FROM metadata
		WHERE object_key = ?
		ALLOW FILTERING`

	iter := m.session.Query(query, objectKey).
		Consistency(gocql.One).
		PageSize(page.Size).
		PageState(page.PageState).
		Iter()

	var items []*entity.MetaData
	for {
		var meta entity.MetaData
		if !iter.Scan(
			&meta.ID,
			&meta.UserID,
			&meta.Filename,
			&meta.ObjectKey,
			&meta.ManifestKey,
			&meta.VideoType,
			&meta.Status,
			&meta.Size,
			&meta.MimeType,
			&meta.DurationSeconds,
			&meta.Width,
			&meta.Height,
			&meta.CreatedAt,
			&meta.UpdatedAt,
		) {
			break
		}
		cp := meta
		items = append(items, &cp)
	}

	nextPageState := iter.PageState()
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to find metadata by object key: %w", err)
	}

	return &pagination.PagedResult[*entity.MetaData]{Items: items, NextPageState: nextPageState}, nil
}

func (m *MetaDataDB) FindAllByStatus(userID gocql.UUID, status string, page pagination.Page) (*pagination.PagedResult[*entity.MetaData], error) {
	const query = `
		SELECT id, user_id, filename, object_key, manifest_key,
		       video_type, status, size, mime_type, duration_seconds,
		       width, height, created_at, updated_at
		FROM metadata
		WHERE user_id = ? AND status = ?
		ALLOW FILTERING`

	iter := m.session.Query(query, userID, status).
		Consistency(gocql.One).
		PageSize(page.Size).
		PageState(page.PageState).
		Iter()

	var items []*entity.MetaData
	for {
		var meta entity.MetaData
		if !iter.Scan(
			&meta.ID,
			&meta.UserID,
			&meta.Filename,
			&meta.ObjectKey,
			&meta.ManifestKey,
			&meta.VideoType,
			&meta.Status,
			&meta.Size,
			&meta.MimeType,
			&meta.DurationSeconds,
			&meta.Width,
			&meta.Height,
			&meta.CreatedAt,
			&meta.UpdatedAt,
		) {
			break
		}
		cp := meta
		items = append(items, &cp)
	}

	nextPageState := iter.PageState()
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to list metadata by status: %w", err)
	}

	return &pagination.PagedResult[*entity.MetaData]{Items: items, NextPageState: nextPageState}, nil
}

func (m *MetaDataDB) DeleteByID(userID gocql.UUID, id gocql.UUID) error {
	if userID == (gocql.UUID{}) || id == (gocql.UUID{}) {
		return fmt.Errorf("user_id and id are required for delete")
	}

	const query = `DELETE FROM metadata WHERE user_id = ? AND id = ?`

	if err := m.session.Query(query, userID, id).Exec(); err != nil {
		return fmt.Errorf("failed to delete metadata: %w", err)
	}

	return nil
}
