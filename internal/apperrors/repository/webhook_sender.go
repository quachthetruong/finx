package repository

type NotifyWebhookRepository interface {
	Send(channel, message string) error
}
