package dto 

import (
	// "time"
)

// Advert represents an album record.
type AdvertDisplay struct {
  AdvertId      int64  `db:"pk,advert_id"`
  CategoryName  string  `db:"category_name"`
  Title         string  `db:"title"`
  Price         int     `db:"price"`
  Currency      string  `db:"currency"`
  ModeratorId   int64   `db:"moderator_id"`
  Created       string  `db:"created"`
  Active        int     `db:"active"`
}
/*
SELECT a.advert_id, c.name, a.title, a.price, a.currency, a.moderator_id, a.created, a.active
FROM advert a, category c
WHERE a.category_id = c.category_id
ORDER BY created ASC
GROUP BY a.category_id
LIMIT 100
*/
