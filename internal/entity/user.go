package entity

// User represents a user.
type User struct {
  UserId       int64  `db:"user_id"`
  Name         string  `db:"name"`
  Email        string  `db:"email"`
  Tel          string  `db:"tel"`
  Impp         string  `db:"impp"`
  Authkey      string  `db:"authkey"`
  Passhash     string  `db:"passhash"`
  Approvetoken string  `db:"approvetoken"`
  Picture      string  `db:"picture"`
  Created      string  `db:"created"`
  Lastlogin    string  `db:"lastlogin"`
  Roles        string  `db:"roles"`
  Notes        string  `db:"notes"`
}

// GetID returns the user ID.
func (u User) GetID() int64 {
	return u.UserId
}

// GetName returns the user name.
func (u User) GetName() string {
	return u.Name
}
