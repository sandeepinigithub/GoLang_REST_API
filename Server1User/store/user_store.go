package store

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/sandeepshukla/golangrestproject1/models"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyExist = errors.New("email already exists")
)

type UserStore struct {
	mu      sync.RWMutex
	users   map[string]*models.User
	counter int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]*models.User),
	}
}

func (s *UserStore) generateID() string {
	s.counter++
	return fmt.Sprintf("%d", s.counter)
}

func (s *UserStore) emailExists(email, excludeID string) bool {
	for _, u := range s.users {
		if u.Email == email && u.ID != excludeID {
			return true
		}
	}
	return false
}

func (s *UserStore) Create(req models.CreateUserRequest) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.emailExists(req.Email, "") {
		return nil, ErrEmailAlreadyExist
	}

	now := time.Now()
	user := &models.User{
		ID:        s.generateID(),
		Name:      req.Name,
		Email:     req.Email,
		Age:       req.Age,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.users[user.ID] = user
	return user, nil
}

func (s *UserStore) GetByID(id string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserStore) GetAll() []*models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]*models.User, 0, len(s.users))
	for _, u := range s.users {
		users = append(users, u)
	}
	return users
}

func (s *UserStore) Update(id string, req models.UpdateUserRequest) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.users[id]
	if !ok {
		return nil, ErrUserNotFound
	}

	if req.Email != user.Email && s.emailExists(req.Email, id) {
		return nil, ErrEmailAlreadyExist
	}

	user.Name = req.Name
	user.Email = req.Email
	user.Age = req.Age
	user.UpdatedAt = time.Now()

	return user, nil
}

func (s *UserStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.users[id]; !ok {
		return ErrUserNotFound
	}

	delete(s.users, id)
	return nil
}
