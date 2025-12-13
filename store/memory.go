package store

import (
	"errors"
	"sync"
)

// MemoryUserStore is an in-memory implementation of UserStore
type MemoryUserStore struct {
	users  map[int]User
	nextID int
	mutex  sync.RWMutex
}

// NewMemoryUserStore creates a new in-memory user store
func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users:  make(map[int]User),
		nextID: 1,
	}
}

// GetAll returns all users
func (m *MemoryUserStore) GetAll() ([]User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	users := make([]User, 0, len(m.users))
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

// GetByID returns a user by ID
func (m *MemoryUserStore) GetByID(id int) (*User, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// Create adds a new user and returns the created user with assigned ID
func (m *MemoryUserStore) Create(user User) (*User, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	user.ID = m.nextID
	m.nextID++
	m.users[user.ID] = user
	return &user, nil
}

// Update modifies an existing user
func (m *MemoryUserStore) Update(id int, user User) (*User, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.users[id]; !exists {
		return nil, errors.New("user not found")
	}

	user.ID = id // Ensure ID matches the parameter
	m.users[id] = user
	return &user, nil
}

// Delete removes a user by ID
func (m *MemoryUserStore) Delete(id int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if _, exists := m.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(m.users, id)
	return nil
}
