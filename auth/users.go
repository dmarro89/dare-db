package auth

import (
	"errors"
	"sync"
)

type UserStore struct {
	usersMu sync.RWMutex
	tokenMu sync.RWMutex
	users   map[string]string
	tokens  map[string]string
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:   make(map[string]string),
		tokens:  make(map[string]string),
		usersMu: sync.RWMutex{},
		tokenMu: sync.RWMutex{},
	}
}

func (store *UserStore) AddUser(username, password string) error {
	store.usersMu.Lock()
	defer store.usersMu.Unlock()

	if _, exists := store.users[username]; exists {
		return errors.New("user already exists")
	}

	store.users[username] = password
	return nil
}

func (store *UserStore) DeleteUser(username string) error {
	store.usersMu.Lock()
	defer store.usersMu.Unlock()

	if _, exists := store.users[username]; !exists {
		return errors.New("user does not exist")
	}

	delete(store.users, username)
	store.DeleteToken(username)
	return nil
}

func (store *UserStore) UpdatePassword(username, newPassword string) error {
	store.usersMu.Lock()
	defer store.usersMu.Unlock()

	if _, exists := store.users[username]; !exists {
		return errors.New("user does not exist")
	}

	store.users[username] = newPassword
	return nil
}

func (store *UserStore) ValidateCredentials(username, password string) bool {
	store.usersMu.RLock()
	defer store.usersMu.RUnlock()

	storedPassword, exists := store.users[username]
	return exists && storedPassword == password
}

func (store *UserStore) SaveToken(username, token string) {
	store.tokenMu.Lock()
	defer store.tokenMu.Unlock()
	store.tokens[username] = token
}

func (store *UserStore) DeleteToken(username string) error {
	store.tokenMu.Lock()
	defer store.tokenMu.Unlock()

	if _, exists := store.tokens[username]; !exists {
		return errors.New("user does not exist")
	}

	delete(store.tokens, username)
	return nil
}

func (store *UserStore) ValidateToken(username, token string) bool {
	store.usersMu.RLock()
	defer store.usersMu.RUnlock()
	storedToken, exists := store.tokens[username]
	return exists && storedToken == token
}
