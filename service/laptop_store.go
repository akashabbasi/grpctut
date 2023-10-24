package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/akashabbasi/pcbook/pb"
	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("record already exists")

// LaptopStore is an interface to store laptop
type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	Find(id string) (*pb.Laptop, error)
}

// InMemoryLaptopStore
type InMemoryLaptopStore struct {
	mu   sync.RWMutex
	data map[string]*pb.Laptop
}

// NewInMemoryLaptopStore returns new  in memory laptop store
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)

	if err != nil {
		return fmt.Errorf("cannot copy laptop data: %w", err)
	}

	store.data[other.Id] = other
	return nil
}

func (store *InMemoryLaptopStore) Find(id string) (*pb.Laptop, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %v", err)
	}

	return other, nil
}
