package entity

import (
	// "time"
)

// Message represents an message table record.
type Message struct {
  MessageId    int64  `db:"pk,message_id"`
  SenderId     int64  `db:"sender_id"`
  ReceiverId   int64  `db:"receiver_id"`
  Content      string  `db:"content"`
  Sended       string  `db:"sended"`
  Readed       string  `db:"readed"`
  ModeratorId  int64  `db:"moderator_id"`
}
