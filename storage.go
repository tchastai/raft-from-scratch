package main

import (
	"errors"
	"sync"
)

type Storage interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

type Store struct {
	mu sync.Mutex
	db map[string]string
}

func NewStore() Storage {
	return &Store{
		db: make(map[string]string),
	}
}

func (s *Store) Get(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.db[key]
	if !ok {
		return "", errors.New("Key does not exists")
	}
	return value, nil
}
func (s *Store) Set(key string, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, _ := s.Get(key)
	if v != "" {
		return errors.New("Key already exists")
	}
	s.db[key] = value
	return nil
}
