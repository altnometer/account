package dbclient

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/satori/uuid"
)

// declare where it is used
// var DBClient dbclient.IBoltClient

// IBoltClient to make mocking easier, interact with db via this interface.
type IBoltClient interface {
	OpenDB()
	InitBucket()
	Get(name string) ([]byte, error)
	Set(name string) error
}

// BoltClient a real Bolt DB Client implementation.
type BoltClient struct {
	boltDB     *bolt.DB
	FileName   string
	BucketName string
}

// OpenDB opens DB connection.
func (bc *BoltClient) OpenDB() {
	var err error
	bc.boltDB, err = bolt.Open(bc.FileName, 0600, nil)
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
}

// Get user data from DB.
func (bc *BoltClient) Get(byname string) ([]byte, error) {
	// acc := account.Account{}
	var id []byte
	err := bc.boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bc.BucketName))
		// accountBytes := b.Get([]byte(byname))
		// accountBytes := b.Get([]byte(byname))
		id = b.Get([]byte(byname))
		// err := json.Unmarshal(accountBytes, &acc)
		// if err != nil {
		// return errors.New("Error when unmarshal db []bytes to acc struct")
		// }
		return nil
	})
	if err != nil {
		// return account.Account{}, err
		return nil, err
	}
	// return acc, nil
	return id, nil
}

// Set user data into DB.
func (bc *BoltClient) Set(name string) error {
	// id := uuid.NewV4().String()
	idUUIDObj, err := uuid.NewV4()
	if err != nil {
		return err
	}
	id := idUUIDObj.Bytes()

	// acc := Account{
	// 	ID:   id,
	// 	Name: name,
	// }
	// jsonBytes, _ := json.Marshal(acc)
	err = bc.boltDB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bc.BucketName))
		// err := b.Put([]byte(name), jsonBytes)
		err := b.Put([]byte(name), id)
		return err
	})
	return err
}

// InitBucket initializes the bucket.
func (bc *BoltClient) InitBucket() {
	err := bc.boltDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bc.BucketName))
		if err != nil {
			return fmt.Errorf("create bucket failed: %s", err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}
