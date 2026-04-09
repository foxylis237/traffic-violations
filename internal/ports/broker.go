package ports

type Broker interface {
	Connect(conn string) error
	Subscribe() (<-chan []byte, error)
	Close() error
}
