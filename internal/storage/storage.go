package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/nivanov045/calendar/internal"
)

type storage struct {
	users       map[string]internal.User //users by id
	usersMutex  sync.RWMutex
	events      map[string]internal.Event //events by id
	eventsMutex sync.RWMutex
}

func New() *storage {
	return &storage{
		users:       map[string]internal.User{},
		usersMutex:  sync.RWMutex{},
		events:      map[string]internal.Event{},
		eventsMutex: sync.RWMutex{},
	}
}

func (s *storage) isUserExist(user string) bool {
	if _, ok := s.users[user]; ok {
		return true
	}
	return false
}

func (s *storage) AddUser(user internal.User) error {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()
	if s.isUserExist(user.ID) {
		return errors.New("user with this id already existed")
	}
	s.users[user.ID] = user
	return nil
}

func (s *storage) AddEvent(event internal.Event) error {
	//TODO: Save in user map info about events to speedup several functions
	s.eventsMutex.Lock()
	defer s.eventsMutex.Unlock()
	if _, ok := s.events[event.ID]; ok {
		return errors.New("event with this id already existed")
	}
	s.events[event.ID] = event
	return nil
}

func (s *storage) GetEvent(id string) (internal.Event, error) {
	s.eventsMutex.RLock()
	defer s.eventsMutex.RUnlock()
	if _, ok := s.events[id]; !ok {
		return internal.Event{}, errors.New("unexisted event")
	}
	return s.events[id], nil
}

func (s *storage) Accept(user string, event string) error {
	s.eventsMutex.Lock()
	defer s.eventsMutex.Unlock()
	if _, ok := s.events[event]; !ok {
		return errors.New("unexisted event")
	}
	for idx, value := range s.events[event].Candidates {
		if value != user {
			continue
		}
		eventTmp := s.events[event]
		eventTmp.Candidates = append(eventTmp.Candidates[:idx], eventTmp.Candidates[idx+1:]...)
		eventTmp.Participants = append(eventTmp.Participants, user)
		s.events[event] = eventTmp
		return nil
	}
	return errors.New("unexisted user in event")
}

func (s *storage) Reject(user string, event string) error {
	s.eventsMutex.Lock()
	defer s.eventsMutex.Unlock()
	if _, ok := s.events[event]; !ok {
		return errors.New("unexisted event")
	}
	for idx, value := range s.events[event].Candidates {
		if value != user {
			continue
		}
		eventTmp := s.events[event]
		eventTmp.Candidates = append(eventTmp.Candidates[:idx], eventTmp.Candidates[idx+1:]...)
		s.events[event] = eventTmp
		return nil
	}
	return errors.New("unexisted user in event")
}

func (s *storage) GetEvents(user string, begin time.Time, end time.Time) ([]internal.Event, error) {
	s.eventsMutex.RLock()
	defer s.eventsMutex.RUnlock()
	if !s.isUserExist(user) {
		return nil, errors.New("unexisted user")
	}
	var result []internal.Event
	for _, curEvent := range s.events {
		for _, participant := range curEvent.Participants {
			if participant == user {
				//TODO: Speedup by use first date not from event start but after |begin|
				switch curEvent.RepeatType {
				case internal.Daily:
					for first := curEvent.Start; first.Before(end); first = first.AddDate(0, 0, 1) {
						curEvent.Finish = curEvent.Finish.Add(first.Sub(curEvent.Start))
						curEvent.Start = first
						if (curEvent.Finish.After(begin) && curEvent.Finish.Before(end)) ||
							(curEvent.Start.After(begin) && curEvent.Start.Before(end)) {
							result = append(result, curEvent)
						}
					}
				case internal.Weekly:
					for first := curEvent.Start; first.Before(end); first = first.AddDate(0, 0, 7) {
						curEvent.Finish = curEvent.Finish.Add(first.Sub(curEvent.Start))
						curEvent.Start = first
						if (curEvent.Finish.After(begin) && curEvent.Finish.Before(end)) ||
							(curEvent.Start.After(begin) && curEvent.Start.Before(end)) {
							result = append(result, curEvent)
						}
					}
				case internal.Workdays:
					for first := curEvent.Start; first.Before(end); first = first.AddDate(0, 0, 1) {
						if first.Weekday() == time.Saturday || first.Weekday() == time.Sunday {
							continue
						}
						curEvent.Finish = curEvent.Finish.Add(first.Sub(curEvent.Start))
						curEvent.Start = first
						if (curEvent.Finish.After(begin) && curEvent.Finish.Before(end)) ||
							(curEvent.Start.After(begin) && curEvent.Start.Before(end)) {
							result = append(result, curEvent)
						}
					}
				case internal.Yearly:
					for first := curEvent.Start; first.Before(end); first = first.AddDate(1, 0, 0) {
						curEvent.Finish = curEvent.Finish.Add(first.Sub(curEvent.Start))
						curEvent.Start = first
						if (curEvent.Finish.After(begin) && curEvent.Finish.Before(end)) ||
							(curEvent.Start.After(begin) && curEvent.Start.Before(end)) {
							result = append(result, curEvent)
						}
					}
				case internal.Once:
					if (curEvent.Finish.After(begin) && curEvent.Finish.Before(end)) ||
						(curEvent.Start.After(begin) && curEvent.Start.Before(end)) {
						result = append(result, curEvent)
					}
				}
				break
			}
		}
	}
	return result, nil
}

func (s *storage) FindFreeSlot(users []string, begin time.Time, duration time.Duration, validUntil time.Time) (from time.Time, err error) {
	s.eventsMutex.RLock()
	defer s.eventsMutex.RUnlock()
	events := map[string][]internal.Event{}
	for _, myUser := range users {
		eventsInner, err := s.GetEvents(myUser, begin, validUntil)
		if err != nil {
			return time.Time{}, err
		}
		events[myUser] = eventsInner
	}
	for begin.Before(validUntil) {
		changed := false
		for _, userEvents := range events {
			for _, curEvent := range userEvents {
				finish := begin.Add(duration)
				if !((curEvent.Start.After(finish) && curEvent.Start.After(begin)) ||
					(curEvent.Finish.Before(finish) && curEvent.Finish.Before(begin))) {
					begin = curEvent.Finish
					changed = true
					break
				}
			}
		}
		if !changed {
			return begin, nil
		}
	}
	return time.Time{}, errors.New("no such slot")
}
