package auth

import (
	"errors"
	"sync"
)

type UserStore struct {
	mu     sync.RWMutex
	users  map[string]string
	tokens map[string]string
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:  make(map[string]string),
		tokens: make(map[string]string),
		mu:     sync.RWMutex{},
	}
}

func (store *UserStore) AddUser(username, password string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.users[username]; exists {
		return errors.New("user already exists")
	}

	store.users[username] = password
	return nil
}

func (store *UserStore) DeleteUser(username string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.users[username]; !exists {
		return errors.New("user does not exist")
	}

	delete(store.users, username)
	return nil
}

func (store *UserStore) UpdatePassword(username, newPassword string) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.users[username]; !exists {
		return errors.New("user does not exist")
	}

	store.users[username] = newPassword
	return nil
}

func (store *UserStore) ValidateCredentials(username, password string) bool {
	store.mu.RLock()
	defer store.mu.RUnlock()

	storedPassword, exists := store.users[username]
	return exists && storedPassword == password
}

func (store *UserStore) SaveToken(username, token string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.tokens[username] = token
}

func (store *UserStore) ValidateToken(username, token string) bool {
	store.mu.RLock()
	defer store.mu.RUnlock()
	storedToken, exists := store.tokens[username]
	return exists && storedToken == token
}
