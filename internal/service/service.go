package service

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/nivanov045/calendar/internal"
)

type service struct {
	storage Storage
}

func New(storage Storage) *service {
	return &service{storage: storage}
}

func (s *service) CreateUser(body []byte) ([]byte, error) {
	var userInfo internal.CustomUserInfo
	err := json.Unmarshal(body, &userInfo)
	if err != nil {
		log.Error().Err(err).Stack()
		return nil, errors.New("wrong query")
	}
	id := uuid.New().String()
	newUser := internal.User{
		Info: internal.CustomUserInfo{},
		ID:   id,
	}
	err = s.storage.AddUser(newUser)
	if err != nil {
		log.Error().Err(err).Stack()
		return nil, err
	}
	type response struct {
		ID string `json:"id"`
	}
	currentResponse := response{ID: id}
	marshal, err := json.Marshal(currentResponse)
	if err != nil {
		log.Error().Err(err).Stack()
		return nil, err
	}
	return marshal, nil
}

func (s *service) CreateEventWithUsers(body []byte) ([]byte, error) {
	var curEvent internal.Event
	err := json.Unmarshal(body, &curEvent)
	if err != nil {
		log.Error().Err(err).Stack()
		return nil, errors.New("wrong query")
	}
	if curEvent.RepeatType < internal.MinRepeatType || internal.MaxRepeatType < curEvent.RepeatType {
		return nil, errors.New("wrong repeat type")
	}
	//TODO: add validation of begin earlier then end
	id := uuid.New().String()
	curEvent.ID = id
	err = s.storage.AddEvent(curEvent)
	if err != nil {
		log.Error().Err(err).Stack()
		if err.Error() == "event with this id already existed" {
			return nil, err
		}
		return nil, errors.New("unable to create event")
	}
	type response struct {
		ID string `json:"id"` //id
	}
	currentResponse := response{ID: id}
	marshal, err := json.Marshal(currentResponse)
	if err != nil {
		log.Error().Err(err).Stack()
		return nil, err
	}
	return marshal, nil
}

func (s *service) GetEventDetails(body []byte) ([]byte, error) {
	type request struct {
		Event string `json:"event"` //event id
	}
	var curRequest request
	err := json.Unmarshal(body, &curRequest)
	if err != nil {
		log.Error().Err(err).Stack()
		return nil, errors.New("wrong query")
	}
	myEvent, err := s.storage.GetEvent(curRequest.Event)
	if err != nil {
		log.Error().Err(err).Stack()
		if err.Error() == "unexisted event" {
			return nil, err
		}
		return nil, errors.New("unable to get event")
	}
	marshal, err := json.Marshal(myEvent)
	if err != nil {
		log.Error().Err(err).Stack()
		return nil, err
	}
	return marshal, nil
}

func (s *service) AcceptInvitation(body []byte) error {
	type request struct {
		User  string `json:"user"`  //user id
		Event string `json:"event"` //event id
	}
	var currentRequest request
	err := json.Unmarshal(body, &currentRequest)
	if err != nil {
		log.Error().Err(err).Stack()
		return errors.New("wrong query")
	}
	err = s.storage.Accept(currentRequest.User, currentRequest.Event)
	if err != nil {
		log.Error().Err(err).Stack()
		if err.Error() == "unexisted event" || err.Error() == "unexisted user in event" {
			return err
		}
		return errors.New("unable to accept invitation")
	}
	return nil
}

func (s *service) RejectInvitation(body []byte) error {
	type request struct {
		User  string `json:"user"`  //user id
		Event string `json:"event"` //event id
	}
	var currentRequest request
	err := json.Unmarshal(body, &currentRequest)
	if err != nil {
		log.Error().Err(err).Stack()
		return errors.New("wrong query")
	}
	err = s.storage.Reject(currentRequest.User, currentRequest.Event)
	if err != nil {
		log.Error().Err(err).Stack()
		if err.Error() == "unexisted event" || err.Error() == "unexisted user in event" {
			return err
		}
		return errors.New("unable to reject invitation")
	}
	return nil
}

func (s *service) GetEvents(body []byte) ([]byte, error) {
	type request struct {
		User string    `json:"user"` //user id
		From time.Time `json:"from"` //from what moment find events
		To   time.Time `json:"to"`   //to what moment find events
	}
	var currentRequest request
	err := json.Unmarshal(body, &currentRequest)
	if err != nil {
		log.Error().Err(err).Stack()
		return nil, errors.New("wrong query")
	}
	res, err := s.storage.GetEvents(currentRequest.User, currentRequest.From, currentRequest.To)
	if err != nil {
		log.Error().Err(err).Stack()
		if err.Error() == "unexisted user" {
			return nil, err
		}
		return nil, errors.New("unable to find events")
	}
	marshal, err := json.Marshal(res)
	if err != nil {
		log.Error().Err(err).Stack()
		return nil, err
	}
	return marshal, nil
}

func (s *service) FindSlot(body []byte) ([]byte, error) {
	type request struct {
		Users      []string      `json:"users"`       //users id
		Duration   time.Duration `json:"duration"`    //duration of event
		ValidUntil time.Time     `json:"valid_until"` //last interesting time
	}
	var currentRequest request
	//TODO: Add custom unmarshall of duration
	err := json.Unmarshal(body, &currentRequest)
	if err != nil {
		log.Error().Err(err).Stack()
		return nil, errors.New("wrong query")
	}
	begin, err := s.storage.FindFreeSlot(currentRequest.Users, time.Now(), currentRequest.Duration, currentRequest.ValidUntil)
	if err != nil {
		log.Error().Err(err).Stack()
		if err.Error() == "unexisted user" {
			return nil, err
		} else if err.Error() == "no such slot" {
			return nil, err
		}
		return nil, errors.New("unable to find space")
	}
	type response struct {
		From time.Time `json:"begin"` //begin of space
	}
	currentResponse := response{From: begin}
	marshal, err := json.Marshal(currentResponse)
	if err != nil {
		log.Error().Err(err).Stack()
		return nil, err
	}
	return marshal, nil
}

//TODO: Add tests.
