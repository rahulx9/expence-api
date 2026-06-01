package main

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"
	"sync"
)

type Item struct {
	ID     string `json:"id"`
	Source string `json:"source" binding:"required"`
	Name   string `json:"name" binding:"required"`
}

type Store struct {
	mu    sync.Mutex
	Path  string
	Items []Item
}

var ErrDuplicateID = errors.New("duplicate id")

var nextID int

func NewStore(path string) *Store {
	return &Store{Path: path}
}

func (s *Store) Load() error {
	if _, err := os.Stat(s.Path); os.IsNotExist(err) {
		s.Items = []Item{}
		nextID = 1
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
		nextID = 1
		return nil
	}

	if err := json.Unmarshal(data, &s.Items); err != nil {
		return err
	}

	maxID := 0
	for _, item := range s.Items {
		id, err := strconv.Atoi(item.ID)
		if err == nil && id > maxID {
			maxID = id
		}
	}
	nextID = maxID + 1
	return nil
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

	item.ID = strconv.Itoa(nextID)
	nextID++

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
