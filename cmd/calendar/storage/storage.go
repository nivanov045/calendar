package storage

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/nivanov045/curly-waffle/internal/event"
	"github.com/nivanov045/curly-waffle/internal/user"
)

type storage struct {
	users       map[string]user.User //users by id
	usersMutex  sync.RWMutex
	events      map[string]event.Event //events by id
	eventsMutex sync.RWMutex
}

func New() *storage {
	return &storage{
		users:       map[string]user.User{},
		usersMutex:  sync.RWMutex{},
		events:      map[string]event.Event{},
		eventsMutex: sync.RWMutex{},
	}
}

func (s *storage) isUserExist(user string) bool {
	s.usersMutex.RLock()
	defer s.usersMutex.RUnlock()
	if _, ok := s.users[user]; ok {
		return true
	}
	return false
}

func (s *storage) AddUser(user user.User) error {
	s.usersMutex.Lock()
	defer s.usersMutex.Unlock()
	if s.isUserExist(user.ID) {
		return errors.New("user with this id already existed")
	}
	s.users[user.ID] = user
	return nil
}

func (s *storage) AddEvent(event event.Event) error {
	//TODO: Save in user map info about events to speedup several functions
	s.eventsMutex.Lock()
	defer s.eventsMutex.Unlock()
	if _, ok := s.events[event.ID]; ok {
		return errors.New("event with this id already existed")
	}
	s.events[event.ID] = event
	log.Println("storage::AddEvent::info: created event with id:", event.ID)
	return nil
}

func (s *storage) GetEvent(id string) (event.Event, error) {
	s.eventsMutex.RLock()
	defer s.eventsMutex.RUnlock()
	if _, ok := s.events[id]; !ok {
		log.Println("storage::GetEvent::warning: can't find event with id:", id)
		return event.Event{}, errors.New("unexisted event")
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

func (s *storage) GetEvents(user string, begin time.Time, end time.Time) ([]event.Event, error) {
	s.eventsMutex.RLock()
	defer s.eventsMutex.RUnlock()
	if !s.isUserExist(user) {
		return nil, errors.New("unexisted user")
	}
	var result []event.Event
	for _, curEvent := range s.events {
		for _, participant := range curEvent.Participants {
			if participant == user {
				//TODO: Speedup by use first date not from event start but after |begin|
				switch curEvent.RepeatType {
				case event.Daily:
					for first := curEvent.Start; first.Before(end); first = first.AddDate(0, 0, 1) {
						curEvent.Finish = curEvent.Finish.Add(first.Sub(curEvent.Start))
						curEvent.Start = first
						if (curEvent.Finish.After(begin) && curEvent.Finish.Before(end)) ||
							(curEvent.Start.After(begin) && curEvent.Start.Before(end)) {
							result = append(result, curEvent)
						}
					}
				case event.Weekly:
					for first := curEvent.Start; first.Before(end); first = first.AddDate(0, 0, 7) {
						curEvent.Finish = curEvent.Finish.Add(first.Sub(curEvent.Start))
						curEvent.Start = first
						if (curEvent.Finish.After(begin) && curEvent.Finish.Before(end)) ||
							(curEvent.Start.After(begin) && curEvent.Start.Before(end)) {
							result = append(result, curEvent)
						}
					}
				case event.Workdays:
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
				case event.Yearly:
					for first := curEvent.Start; first.Before(end); first = first.AddDate(1, 0, 0) {
						curEvent.Finish = curEvent.Finish.Add(first.Sub(curEvent.Start))
						curEvent.Start = first
						if (curEvent.Finish.After(begin) && curEvent.Finish.Before(end)) ||
							(curEvent.Start.After(begin) && curEvent.Start.Before(end)) {
							result = append(result, curEvent)
						}
					}
				case event.Once:
					if (curEvent.Finish.After(begin) && curEvent.Finish.Before(end)) ||
						(curEvent.Start.After(begin) && curEvent.Start.Before(end)) {
						result = append(result, curEvent)
					}
				}
				break
			}
		}
	}
	log.Println("storage::GetEvents::info: found", len(result), "events from", len(s.events))
	return result, nil
}

func (s *storage) FindFreeSlot(users []string, begin time.Time, duration time.Duration, validUntil time.Time) (from time.Time, err error) {
	s.eventsMutex.RLock()
	defer s.eventsMutex.RUnlock()
	events := map[string][]event.Event{}
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
