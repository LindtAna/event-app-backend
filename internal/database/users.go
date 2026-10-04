package database

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatarUrl"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO users(email, password, name, bio, avatar_url) VALUES(?, ?, ?, '', '')"

	result, err := m.DB.ExecContext(ctx, query, user.Email, user.Password, user.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.Id = int(id)
	return nil
}

func (m *UserModel) getUser(query string, args ...interface{}) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password, &user.Bio, &user.AvatarUrl,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (m *UserModel) Get(id int) (*User, error) {
	query := "SELECT id, name, email, password FROM users WHERE id = ?"
	return m.getUser(query, id)

}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := "SELECT id, name, email, password FROM users WHERE email = ?"
	return m.getUser(query, email)

}

// Aktualisierung des Profils
func (m *UserModel) UpdateProfile(id int, name, bio, avatarUrl string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE users SET name = ?, bio = ?, avatar_url = ? WHERE id = ?"
	_, err := m.DB.ExecContext(ctx, query, name, bio, avatarUrl, id)
	return err
}
