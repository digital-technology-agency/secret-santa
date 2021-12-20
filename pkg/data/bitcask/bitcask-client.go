package bitcask

import (
	"fmt"
	prolog "git.mills.io/prologic/bitcask"
	"os"
)

var (
	MAX_PART_CASH_SIZE = 1024 * 1024 * 100
)

type Data struct {
	db *prolog.Bitcask
}

//Connect init data
func Connect(path string) (*Data, error) {
	tmpPath := fmt.Sprintf("%s/%s", os.TempDir(), path)
	base, err := prolog.Open(tmpPath, prolog.WithMaxDatafileSize(MAX_PART_CASH_SIZE), prolog.WithAutoRecovery(true))
	if err != nil {
		return nil, err
	}
	return &Data{
		db: base,
	}, nil
}

//Add add value by key
func (data *Data) Add(key, value []byte) error {
	return data.db.Put(key, value)
}

//Get get value by key
func (data *Data) Get(key []byte) ([]byte, error) {
	return data.db.Get(key)
}

//Remove remove value by key
func (data *Data) Remove(key []byte) error {
	return data.db.Delete(key)
}

//GetAll get all values
func (data *Data) GetAll() (map[string][]byte, error) {
	result := map[string][]byte{}
	err := data.db.Fold(func(key []byte) error {
		get, err := data.db.Get(key)
		if err != nil {
			return err
		}
		result[string(key)] = get
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
