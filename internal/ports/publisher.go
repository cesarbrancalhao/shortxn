package ports

import "Shortxn/internal/domain"

type EventPublisher interface {
	PublishClickEvent(event *domain.ClickEvent) error
	Close() error
}
