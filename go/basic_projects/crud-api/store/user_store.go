package store

import (
	"errors"
	"sync"

	"crud-api/models"
)


// mu protects shared state from concurrent HTTP requests
type UserStore struct {
	mu     sync.Mutex
	users  []models.User
	nextID int
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  []models.User{},
		nextID: 1,
	}
}

func (s *UserStore) Create(name, email string) models.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := models.User{
		ID:    s.nextID,
		Name:  name,
		Email: email,
	}

	s.nextID++
	s.users = append(s.users, user)

	return user
}

func (s *UserStore) GetAll() []models.User {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.users
}


func (s *UserStore) GetByID(id int) (models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, user := range s.users {
		if user.ID == id {
			return user, nil
		}
	}

	return models.User{}, errors.New("user not found")
}


func (s *UserStore) Update(id int, name, email string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, user := range s.users {
		if user.ID == id {
			s.users[i].Name = name
			s.users[i].Email = email
			return nil
		}
	}

	return errors.New("user not found")
}

func (s *UserStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, user := range s.users {
		if user.ID == id {
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}

	return errors.New("user not found")
}




