package entity

import (
	// "time"
)

// Category represents an category record.
type Category struct {
  CategoryId int64  `db:"pk,category_id"`
  ParentId   int64  `db:"parent_id"`
  Name       string  `db:"name"`
  Slug       string  `db:"slug"`
  IsActive   int  `db:"is_active"`
}

type CategoryPath struct {
  CategoryId int64  `db:"pk,category_id"`
  Path       string  `db:path`
  // IsActive   int  `db:"is_active"`
}
