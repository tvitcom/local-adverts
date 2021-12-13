package dto

import (
	// "time"
)

//sender_id | name      | email           | tel  | sended       
type MessageSender struct {
  SenderId     int64  `db:"sender_id"`
  Name        *string  `db:"name"`
  Email       *string  `db:"email"`
  Tel         *string  `db:"tel"`
  Sended      *string  `db:"sended"`
  ModeratorId  int64  `db:"moderator_id"`
}

/* FORMAT JSON-HTTP COMMUNICATION
{
  "ok": resultOfcurrentService,
  "data": "Данные или текст об ошибке работы сайта",
}
*/

/* MQ FORMAT:
{
"topic":"[sublect]",
  "channel":"todo",
  "channel":"done",
}
*/