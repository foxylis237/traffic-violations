package storage

import (
	"github.com/kvolis/tesgode/dog"
)

type DogStorage struct {
	client *dog.Dog
}

func New() *DogStorage {
	return &DogStorage{
		client: dog.New(),
	}
}

func (s *DogStorage) Connect(conn string) error {
	return s.client.Connect(conn)
}

func (s *DogStorage) Save(key string, value []byte) (int, error) {
	return s.client.Insert(key, value)
}

func (s *DogStorage) Close() error {
	return s.client.Close()
}
