package ports

type Storage interface {
	Save(key string, value []byte) (int, error)
	Close() error
}
