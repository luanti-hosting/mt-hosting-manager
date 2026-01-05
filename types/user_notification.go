package types

type UserNotification struct {
	ID        string `json:"id" gorm:"primarykey;column:id"`
	UserID    string `json:"user_id" gorm:"column:user_id"`
	Timestamp int64  `json:"timestamp" gorm:"column:timestamp"`
	Seen      bool   `json:"seen" gorm:"column:seen"`
	Title     string `json:"title" gorm:"column:title"`
	Message   string `json:"message" gorm:"column:message"`
}

func (*UserNotification) TableName() string {
	return "user_notification"
}
