// Package user
package user

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// Storage - Интерфейсы определяются на стороне ПОТРЕБИТЕЛЯ (consumer-side interfaces)
type Storage struct {
	mu     sync.RWMutex
	users  map[string]*User
	nextID int
}

func NewStorage() *Storage {
	return &Storage{
		users:  make(map[string]*User),
		nextID: 1,
	}
}

func (s *Storage) Create(name, email string) (*User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := &User{
		ID:        fmt.Sprintf("user_%d", s.nextID),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	s.users[user.ID] = user
	s.nextID++

	return user, nil
}

func (s *Storage) List() []*User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].CreatedAt.After(users[j].CreatedAt)
	})

	return users
}

func (s *Storage) GetByID(id string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %s", id)
	}

	return user, nil
}

func (s *Storage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.users)
}
