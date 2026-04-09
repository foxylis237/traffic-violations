package ports

type Storage interface {
	Connect(conn string) error
	Save(key string, value []byte) (int, error)
	Close() error
}
