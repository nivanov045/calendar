package internal

import "time"

type User struct {
	Info CustomUserInfo `json:"info"`         // info about user
	ID   string         `json:"id,omitempty"` //id
}

type CustomUserInfo struct {
	Name string `json:"name"` // user's name
}

type Event struct {
	ID           string          `json:"id,omitempty"`          //id
	Candidates   []string        `json:"candidates,omitempty"`  //list of candidates
	Participants []string        `json:"participants"`          //list of participants, at list one required
	Start        time.Time       `json:"start"`                 //start time, required
	Finish       time.Time       `json:"finish"`                //finish time, required
	RepeatType   RepeatType      `json:"repeat_type,omitempty"` //type of repeating
	Info         CustomEventInfo `json:"info,omitempty"`        //info about event
}

type RepeatType int

const (
	Once RepeatType = iota
	Daily
	Weekly
	Yearly
	Workdays
	MinRepeatType = Once
	MaxRepeatType = Workdays
)

type CustomEventInfo struct {
	Description string `json:"description,omitempty"` //description of event
	Name        string `json:"name,omitempty"`        //event's name
}
