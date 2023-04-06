package service

import (
	"time"

	"github.com/nivanov045/calendar/internal"
)

type Storage interface {
	AddUser(user internal.User) error
	AddEvent(event internal.Event) error
	GetEvent(id string) (internal.Event, error)
	Accept(user string, event string) error
	Reject(user string, event string) error
	GetEvents(user string, begin time.Time, end time.Time) ([]internal.Event, error)
	FindFreeSlot(users []string, begin time.Time, duration time.Duration, validUntil time.Time) (from time.Time, err error)
}
