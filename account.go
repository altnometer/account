package account

// Account holds core user details.
type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RegRawData holds a user submitted account details.
type RegRawData struct {
	Account
	RawPassword string `json:"raw_password"`
}

// RegData holds ready for db data.
type RegData struct {
	Account
	EncodedPassword string `json:"encoded_password"`
}
