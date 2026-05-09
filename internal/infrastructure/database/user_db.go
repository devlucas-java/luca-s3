package database

import (
	"fmt"
	"time"

	"github.com/gocql/gocql"

	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/repository"
	"github.com/devlucas-java/luca-s3/pkg/id"
)

type UserDB struct {
	session *gocql.Session
}

func NewUserDB(session *gocql.Session) repository.UserRepository {
	return &UserDB{session: session}
}

// scanUser reads a full user row from a gocql.Scanner into an entity.User.
func scanUser(scan func(...interface{}) error) (*entity.User, error) {
	var (
		cassandraID gocql.UUID
		name        string
		email       string
		username    string
		password    string
		roles       []string
		createdAt   time.Time
		updatedAt   time.Time
	)

	if err := scan(&cassandraID, &name, &email, &username, &password, &roles, &createdAt, &updatedAt); err != nil {
		return nil, err
	}

	userID, err := id.Parse(cassandraID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse user id: %w", err)
	}

	return &entity.User{
		ID:        userID,
		Name:      name,
		Email:     email,
		Username:  username,
		Password:  password,
		Roles:     roles,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// Create inserts a new user and returns the persisted entity.
func (u *UserDB) Create(user *entity.User) (*entity.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	if id.IsNil(user.ID) {
		user.ID = id.NewUUID()
	}

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	const query = `INSERT INTO users (id, name, email, username, password, roles, created_at, updated_at)
	               VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	if err := u.session.Query(query,
		gocql.UUID(user.ID),
		user.Name,
		user.Email,
		user.Username,
		user.Password,
		user.Roles,
		user.CreatedAt,
		user.UpdatedAt,
	).Exec(); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Save is an upsert — inserts or fully replaces the user row.
func (u *UserDB) Save(user *entity.User) (*entity.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	if id.IsNil(user.ID) {
		user.ID = id.NewUUID()
	}

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	// INSERT in Cassandra is naturally an upsert (LWT not needed here)
	const query = `INSERT INTO users (id, name, email, username, password, roles, created_at, updated_at)
	               VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	if err := u.session.Query(query,
		gocql.UUID(user.ID),
		user.Name,
		user.Email,
		user.Username,
		user.Password,
		user.Roles,
		user.CreatedAt,
		user.UpdatedAt,
	).Exec(); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

// Updates applies a partial update — only non-zero fields are written.
// Fields left at their zero value are kept as-is in the database.
func (u *UserDB) Updates(user *entity.User) (*entity.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	if id.IsNil(user.ID) {
		return nil, fmt.Errorf("user id is required for updates")
	}

	// Fetch current state so we can merge
	current, err := u.FindByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user for partial update: %w", err)
	}
	if current == nil {
		return nil, fmt.Errorf("user not found")
	}

	if user.Name != "" {
		current.Name = user.Name
	}
	if user.Email != "" {
		current.Email = user.Email
	}
	if user.Username != "" {
		current.Username = user.Username
	}
	if user.Password != "" {
		current.Password = user.Password
	}
	if len(user.Roles) > 0 {
		current.Roles = user.Roles
	}
	current.UpdatedAt = time.Now()

	const query = `UPDATE users
	               SET name = ?, email = ?, username = ?, password = ?, roles = ?, updated_at = ?
	               WHERE id = ?`

	if err := u.session.Query(query,
		current.Name,
		current.Email,
		current.Username,
		current.Password,
		current.Roles,
		current.UpdatedAt,
		gocql.UUID(current.ID),
	).Exec(); err != nil {
		return nil, fmt.Errorf("failed to partially update user: %w", err)
	}

	return current, nil
}

// FindByID returns the user with the given UUID, or nil if not found.
func (u *UserDB) FindByID(userID id.UUID) (*entity.User, error) {
	const query = `SELECT id, name, email, username, password, roles, created_at, updated_at
	               FROM users WHERE id = ?`

	user, err := scanUser(
		u.session.Query(query, gocql.UUID(userID)).Consistency(gocql.One).Scan,
	)
	if err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return user, nil
}

// FindByEmailOrUsername searches by email first, then by username.
// Returns nil (no error) when the user does not exist.
func (u *UserDB) FindByEmailOrUsername(str string) (*entity.User, error) {
	// Try email index first
	user, err := u.findByField("email", str)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}

	// Fall back to username index
	return u.findByField("username", str)
}

// findByField is a helper that queries a secondary-indexed column.
func (u *UserDB) findByField(field, value string) (*entity.User, error) {
	query := fmt.Sprintf(
		`SELECT id, name, email, username, password, roles, created_at, updated_at
		 FROM users WHERE %s = ? LIMIT 1 ALLOW FILTERING`, field,
	)

	iter := u.session.Query(query, value).Consistency(gocql.One).Iter()
	defer iter.Close()

	var (
		cassandraID gocql.UUID
		name        string
		email       string
		username    string
		password    string
		roles       []string
		createdAt   time.Time
		updatedAt   time.Time
	)

	if iter.Scan(&cassandraID, &name, &email, &username, &password, &roles, &createdAt, &updatedAt) {
		userID, err := id.Parse(cassandraID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to parse user id: %w", err)
		}
		return &entity.User{
			ID:        userID,
			Name:      name,
			Email:     email,
			Username:  username,
			Password:  password,
			Roles:     roles,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}, nil
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to find user by %s: %w", field, err)
	}

	return nil, nil
}

// DeleteByID removes the user with the given UUID.
func (u *UserDB) DeleteByID(userID id.UUID) error {
	if id.IsNil(userID) {
		return fmt.Errorf("user id is required for delete")
	}

	const query = `DELETE FROM users WHERE id = ?`
	if err := u.session.Query(query, gocql.UUID(userID)).Exec(); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
