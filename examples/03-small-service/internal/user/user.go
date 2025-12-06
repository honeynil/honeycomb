package user

import (
	"errors"
	"strings"
	"time"
)

// User represents a user in the system
// Идиома Go: структуры данных в основном файле пакета (user.go)
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Validate checks if user data is valid
// Идиома Go: методы валидации на структурах (не отдельный validator)
func (u *User) Validate() error {
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("name is required")
	}
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}
	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email format")
	}
	return nil
}

// IsRecent returns true if user was created in last 24 hours
func (u *User) IsRecent() bool {
	return time.Since(u.CreatedAt) < 24*time.Hour
}
