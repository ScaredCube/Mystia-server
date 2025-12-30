package db

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type User struct {
	ID           string
	Username     string
	PasswordHash string
	Nickname     string
	Role         string
	CreatedAt    time.Time
}

type Channel struct {
	ID        string
	Name      string
	CreatedAt time.Time
}

type DB struct {
	*sql.DB
}

func NewDB(path string) (*DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	d := &DB{db}
	if err := d.init(); err != nil {
		return nil, err
	}

	return d, nil
}

func (d *DB) init() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			nickname TEXT NOT NULL,
			role TEXT NOT NULL,
			created_at DATETIME NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS channels (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			created_at DATETIME NOT NULL
		)`,
	}

	for _, q := range queries {
		if _, err := d.Exec(q); err != nil {
			return err
		}
	}

	return nil
}

func (d *DB) CreateUser(user *User) error {
	_, err := d.Exec(
		"INSERT INTO users (id, username, password_hash, nickname, role, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		user.ID, user.Username, user.PasswordHash, user.Nickname, user.Role, user.CreatedAt,
	)
	return err
}

func (d *DB) GetUserByUsername(username string) (*User, error) {
	row := d.QueryRow("SELECT id, username, password_hash, nickname, role, created_at FROM users WHERE username = ?", username)
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Nickname, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *DB) GetUserByID(id string) (*User, error) {
	row := d.QueryRow("SELECT id, username, password_hash, nickname, role, created_at FROM users WHERE id = ?", id)
	var u User
	err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Nickname, &u.Role, &u.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (d *DB) GetUserCount() (int, error) {
	var count int
	err := d.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	return count, err
}

func (d *DB) UpdateUserRole(userID string, role string) error {
	_, err := d.Exec("UPDATE users SET role = ? WHERE id = ?", role, userID)
	return err
}

func (d *DB) ListUsers() ([]*User, error) {
	rows, err := d.Query("SELECT id, username, nickname, role, created_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.Nickname, &u.Role, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}

func (d *DB) CreateChannel(channel *Channel) error {
	_, err := d.Exec("INSERT INTO channels (id, name, created_at) VALUES (?, ?, ?)", channel.ID, channel.Name, channel.CreatedAt)
	return err
}

func (d *DB) ListChannels() ([]*Channel, error) {
	rows, err := d.Query("SELECT id, name, created_at FROM channels")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*Channel
	for rows.Next() {
		var c Channel
		if err := rows.Scan(&c.ID, &c.Name, &c.CreatedAt); err != nil {
			return nil, err
		}
		channels = append(channels, &c)
	}
	return channels, nil
}
