package storage

type KeyValueRegister interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
	Remove(key string) error

	RetrieveURLForPutOperation(key string) (string, error)
	RetrieveURLForGetOperation(key string) (string, error)
}
