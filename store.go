package main

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Item struct {
	ID    string `json:"id" binding:"required"`
	Value string `json:"value" binding:"required"`
}

type Store struct {
	mu    sync.Mutex
	Path  string
	Items []Item
}

var ErrDuplicateID = errors.New("duplicate id")

func NewStore(path string) *Store {
	return &Store{Path: path}
}

func (s *Store) Load() error {
	if _, err := os.Stat(s.Path); os.IsNotExist(err) {
		s.Items = []Item{}
		return nil
	} else if err != nil {
		return err
	}

	data, err := os.ReadFile(s.Path)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		s.Items = []Item{}
		return nil
	}

	return json.Unmarshal(data, &s.Items)
}

func (s *Store) Save() error {
	data, err := json.MarshalIndent(s.Items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.Path, data, 0644)
}

func (s *Store) AddItem(item Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, existing := range s.Items {
		if existing.ID == item.ID {
			return ErrDuplicateID
		}
	}

	s.Items = append(s.Items, item)
	return s.Save()
}

func (s *Store) GetAll() []Item {
	s.mu.Lock()
	defer s.mu.Unlock()

	copyItems := make([]Item, len(s.Items))
	copy(copyItems, s.Items)
	return copyItems
}

func (s *Store) DeleteItem(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, item := range s.Items {
		if item.ID == id {
			s.Items = append(s.Items[:i], s.Items[i+1:]...)
			s.Save()
			return true
		}
	}
	return false
}
