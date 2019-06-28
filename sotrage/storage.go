package storage

import (
	. "github.com/tobira-shoe/event-models"
)

type Storage interface {
	Open() error
	Close() error

	SaveEventsList(events []DouEvent) error
	SaveTagsList(tags []string) error

	GetEventsList(limit, offset int, tags []string) ([]DouEvent, error)
	GetEvent(id int) (DouEvent, error)
	GetTagsList() ([]string, error)
}
