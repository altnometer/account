package mocks

// BoltClient is a mock of IBoltClient interface.
type BoltClient struct {
	OpenDBIsCalled bool
	GetCall        struct {
		Receives struct {
			Name string
		}
		Returns struct {
			ID    []byte
			Error error
		}
	}
	SetCall struct {
		Receives struct {
			Name string
		}
		Returns struct {
			Error error
		}
	}
	InitBucketCalled bool
}

// OpenDB opens DB connection.
func (bc *BoltClient) OpenDB() {
	bc.OpenDBIsCalled = true
}

// Get user data from DB.
func (bc *BoltClient) Get(name string) ([]byte, error) {
	bc.GetCall.Receives.Name = name
	return bc.GetCall.Returns.ID, bc.GetCall.Returns.Error
}

// Set user data into DB.
func (bc *BoltClient) Set(name string) error {
	bc.SetCall.Receives.Name = name
	return bc.SetCall.Returns.Error

}

// InitBucket initializes the bucket.
func (bc *BoltClient) InitBucket() {
	bc.InitBucketCalled = true
}
