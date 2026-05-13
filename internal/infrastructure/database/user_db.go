package database

import (
	"fmt"
	"time"

	"github.com/devlucas-java/luca-s3/internal/domain/entity"
	"github.com/devlucas-java/luca-s3/internal/infrastructure/repository"
	"github.com/devlucas-java/luca-s3/pkg/pagination"
	"github.com/gocql/gocql"
)

type UserDB struct {
	session *gocql.Session
}

func NewUserDB(session *gocql.Session) repository.UserRepository {
	return &UserDB{session: session}
}

func (u *UserDB) Create(user *entity.User) (*entity.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	if user.ID == (gocql.UUID{}) {
		uuid, err := gocql.RandomUUID()
		if err != nil {
			return nil, fmt.Errorf("failed to generate user uuid: %w", err)
		}
		user.ID = uuid
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	const query = `
		INSERT INTO users (id, name, email, username, password, roles, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	if err := u.session.Query(query,
		user.ID,
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

func (u *UserDB) Save(user *entity.User) (*entity.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	if user.ID == (gocql.UUID{}) {
		return nil, fmt.Errorf("user id is required for save")
	}

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	const query = `
		INSERT INTO users (id, name, email, username, password, roles, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	if err := u.session.Query(query,
		user.ID,
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

func (u *UserDB) Updates(user *entity.User) (*entity.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user is nil")
	}

	if user.ID == (gocql.UUID{}) {
		return nil, fmt.Errorf("user id is required for updates")
	}

	current, err := u.FindByID(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
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

	const query = `
		UPDATE users
		SET name = ?, email = ?, username = ?, password = ?, roles = ?, updated_at = ?
		WHERE id = ?`

	if err := u.session.Query(query,
		current.Name,
		current.Email,
		current.Username,
		current.Password,
		current.Roles,
		current.UpdatedAt,
		current.ID,
	).Exec(); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return current, nil
}

func (u *UserDB) FindByID(userID gocql.UUID) (*entity.User, error) {
	const query = `
		SELECT id, name, email, username, password, roles, created_at, updated_at
		FROM users
		WHERE id = ?`

	var user entity.User
	if err := u.session.Query(query, userID).Consistency(gocql.One).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Roles,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return &user, nil
}

func (u *UserDB) FindByEmail(email string, page pagination.Page) (*pagination.PagedResult[*entity.User], error) {
	const query = `
		SELECT id, name, email, username, password, roles, created_at, updated_at
		FROM users
		WHERE email = ?
		ALLOW FILTERING`

	iter := u.session.Query(query, email).
		Consistency(gocql.One).
		PageSize(page.Size).
		PageState(page.PageState).
		Iter()

	var items []*entity.User
	for {
		var user entity.User
		if !iter.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.Roles,
			&user.CreatedAt,
			&user.UpdatedAt,
		) {
			break
		}
		cp := user
		items = append(items, &cp)
	}

	nextPageState := iter.PageState()
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to find users by email: %w", err)
	}

	return &pagination.PagedResult[*entity.User]{Items: items, NextPageState: nextPageState}, nil
}

func (u *UserDB) FindByUsername(username string, page pagination.Page) (*pagination.PagedResult[*entity.User], error) {
	const query = `
		SELECT id, name, email, username, password, roles, created_at, updated_at
		FROM users
		WHERE username = ?
		ALLOW FILTERING`

	iter := u.session.Query(query, username).
		Consistency(gocql.One).
		PageSize(page.Size).
		PageState(page.PageState).
		Iter()

	var items []*entity.User
	for {
		var user entity.User
		if !iter.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Username,
			&user.Password,
			&user.Roles,
			&user.CreatedAt,
			&user.UpdatedAt,
		) {
			break
		}
		cp := user
		items = append(items, &cp)
	}

	nextPageState := iter.PageState()
	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf("failed to find users by username: %w", err)
	}

	return &pagination.PagedResult[*entity.User]{Items: items, NextPageState: nextPageState}, nil
}

func (u *UserDB) FindByEmailOrUsername(login string) (*entity.User, error) {
	page := pagination.Page{Size: 1}

	emailResult, err := u.FindByEmail(login, page)
	if err != nil {
		return nil, err
	}
	if len(emailResult.Items) > 0 {
		return emailResult.Items[0], nil
	}

	usernameResult, err := u.FindByUsername(login, page)
	if err != nil {
		return nil, err
	}
	if len(usernameResult.Items) > 0 {
		return usernameResult.Items[0], nil
	}

	return nil, nil
}

func (u *UserDB) ExistsByEmailOrUsername(str string) (bool, error) {
	user, err := u.FindByEmailOrUsername(str)
	if err != nil && user == nil {
		return false, err
	}
	return true, err
}

func (u *UserDB) DeleteByID(userID gocql.UUID) error {
	if userID == (gocql.UUID{}) {
		return fmt.Errorf("user id is required for delete")
	}

	const query = `DELETE FROM users WHERE id = ?`

	if err := u.session.Query(query, userID).Exec(); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
