package entity

import (
	// "time"
)

// Advert represents an album record.
type Advert struct {
  AdvertId    int64  `db:"pk,advert_id"`
  UserId      int64  `db:"user_id"`
  CategoryId  int64  `db:"category_id"`
  Title       string  `db:"title"`
  Nanopost    string  `db:"nanopost"`
  Price       int     `db:"price"`
  Currency    string  `db:"currency"`
  Picture1    string  `db:"picture1"`
  Picture2    string  `db:"picture2"`
  Picture3    string  `db:"picture3"`
  Picture4    string  `db:"picture4"`
  // Picture5    string  `db:"picture5"`
  // Picture6    string  `db:"picture6"`
  ModeratorId int64   `db:"moderator_id"`
  Created     string  `db:"created"`
  Active      int     `db:"active"`
}
