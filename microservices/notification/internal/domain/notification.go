package domain

type NotificationType string

const (
	NotificationPaymentCompleted NotificationType = "payment_completed"
)

type Notification struct {
	ID      uint
	UserID  uint
	Type    NotificationType
	Subject string
	Body    string
}