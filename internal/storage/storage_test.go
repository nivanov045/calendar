package storage

import (
	"reflect"
	"testing"
	"time"

	"github.com/nivanov045/calendar/internal"
)

func Test_storage_isUserExist(t *testing.T) {
	type fields struct {
		users map[string]internal.User
	}
	type args struct {
		user string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "Empty users storage",
			fields: fields{users: map[string]internal.User{}},
			args:   args{"qwerty-12345"},
			want:   false,
		},
		{
			name: "User exists in storage",
			fields: fields{users: map[string]internal.User{
				"qwerty-12345": {
					Info: internal.CustomUserInfo{},
					ID:   "qwerty-12345",
				},
			}},
			args: args{"qwerty-12345"},
			want: true,
		},
		{
			name: "User doesn't exists in storage",
			fields: fields{users: map[string]internal.User{
				"qwerty-123456": {
					Info: internal.CustomUserInfo{},
					ID:   "qwerty-123456",
				},
			}},
			args: args{"qwerty-12345"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				users: tt.fields.users,
			}
			if got := s.isUserExist(tt.args.user); got != tt.want {
				t.Errorf("isUserExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_AddUser(t *testing.T) {
	type fields struct {
		users map[string]internal.User
	}
	type args struct {
		user internal.User
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantStorage map[string]internal.User
	}{
		{
			name:   "Add to empty users storage",
			fields: fields{users: map[string]internal.User{}},
			args: args{user: internal.User{
				Info: internal.CustomUserInfo{
					Name: "Ivan",
				},
				ID: "qwerty-12345",
			}},
			wantErr: false,
			wantStorage: map[string]internal.User{
				"qwerty-12345": {
					Info: internal.CustomUserInfo{
						Name: "Ivan",
					},
					ID: "qwerty-12345",
				},
			},
		},
		{
			name: "Add to filled users storage",
			fields: fields{users: map[string]internal.User{
				"qwerty-123456": {
					Info: internal.CustomUserInfo{
						Name: "Petr",
					},
					ID: "qwerty-123456",
				}}},
			args: args{user: internal.User{
				Info: internal.CustomUserInfo{
					Name: "Ivan",
				},
				ID: "qwerty-12345",
			}},
			wantErr: false,
			wantStorage: map[string]internal.User{
				"qwerty-12345": {
					Info: internal.CustomUserInfo{
						Name: "Ivan",
					},
					ID: "qwerty-12345",
				},
				"qwerty-123456": {
					Info: internal.CustomUserInfo{
						Name: "Petr",
					},
					ID: "qwerty-123456",
				},
			},
		},
		{
			name: "Add user with existing id",
			fields: fields{users: map[string]internal.User{
				"qwerty-12345": {
					Info: internal.CustomUserInfo{
						Name: "Nikolay",
					},
					ID: "qwerty-12345",
				}}},
			args: args{user: internal.User{
				Info: internal.CustomUserInfo{
					Name: "Ivan",
				},
				ID: "qwerty-12345",
			}},
			wantErr: true,
			wantStorage: map[string]internal.User{
				"qwerty-12345": {
					Info: internal.CustomUserInfo{
						Name: "Nikolay",
					},
					ID: "qwerty-12345",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				users: tt.fields.users,
			}
			if err := s.AddUser(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
				eq := reflect.DeepEqual(s.users, tt.wantStorage)
				t.Errorf("AddUser() storage= %v, want %v", eq, true)
			}
		})
	}
}

func first(t time.Time, _ error) time.Time {
	return t
}

func Test_storage_AddEvent(t *testing.T) {
	type fields struct {
		events map[string]internal.Event
	}
	type args struct {
		event internal.Event
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantStorage map[string]internal.Event
	}{
		{
			name:   "Add to empty storage",
			fields: fields{map[string]internal.Event{}},
			args: args{
				event: internal.Event{
					ID:           "qwerty-12345",
					Candidates:   []string{"c-1"},
					Participants: []string{"p-1"},
					Start:        first(time.Parse(time.RFC3339, "2022-09-02T10:00:00Z")),
					Finish:       first(time.Parse(time.RFC3339, "2022-09-02T12:00:00Z")),
					RepeatType:   0,
					Info:         internal.CustomEventInfo{},
				},
			},
			wantErr: false,
			wantStorage: map[string]internal.Event{
				"qwerty-12345": {
					ID:           "qwerty-12345",
					Candidates:   []string{"c-1"},
					Participants: []string{"p-1"},
					Start:        first(time.Parse(time.RFC3339, "2022-09-02T10:00:00Z")),
					Finish:       first(time.Parse(time.RFC3339, "2022-09-02T12:00:00Z")),
					RepeatType:   0,
					Info:         internal.CustomEventInfo{},
				}},
		},
		{
			name: "Add to not empty storage",
			fields: fields{map[string]internal.Event{
				"qwerty-123456": {
					ID:           "qwerty-123456",
					Candidates:   []string{"c-1"},
					Participants: []string{"p-1"},
					Start:        first(time.Parse(time.RFC3339, "2022-09-02T10:00:00Z")),
					Finish:       first(time.Parse(time.RFC3339, "2022-09-02T12:00:00Z")),
					RepeatType:   0,
					Info:         internal.CustomEventInfo{},
				}}},
			args: args{
				event: internal.Event{
					ID:           "qwerty-12345",
					Candidates:   []string{"c-1"},
					Participants: []string{"p-1"},
					Start:        first(time.Parse(time.RFC3339, "2022-09-02T10:00:00Z")),
					Finish:       first(time.Parse(time.RFC3339, "2022-09-02T12:00:00Z")),
					RepeatType:   0,
					Info:         internal.CustomEventInfo{},
				},
			},
			wantErr: false,
			wantStorage: map[string]internal.Event{
				"qwerty-12345": {
					ID:           "qwerty-12345",
					Candidates:   []string{"c-1"},
					Participants: []string{"p-1"},
					Start:        first(time.Parse(time.RFC3339, "2022-09-02T10:00:00Z")),
					Finish:       first(time.Parse(time.RFC3339, "2022-09-02T12:00:00Z")),
					RepeatType:   0,
					Info:         internal.CustomEventInfo{},
				}},
		},
		{
			name: "Add existing id storage",
			fields: fields{map[string]internal.Event{
				"qwerty-12345": {
					ID:           "qwerty-12345",
					Candidates:   []string{"c-1"},
					Participants: []string{"p-1"},
					Start:        first(time.Parse(time.RFC3339, "2022-09-02T11:00:00Z")),
					Finish:       first(time.Parse(time.RFC3339, "2022-09-02T15:00:00Z")),
					RepeatType:   0,
					Info:         internal.CustomEventInfo{},
				}}},
			args: args{
				event: internal.Event{
					ID:           "qwerty-12345",
					Candidates:   []string{"c-1"},
					Participants: []string{"p-1"},
					Start:        first(time.Parse(time.RFC3339, "2022-09-02T10:00:00Z")),
					Finish:       first(time.Parse(time.RFC3339, "2022-09-02T12:00:00Z")),
					RepeatType:   0,
					Info:         internal.CustomEventInfo{},
				},
			},
			wantErr: true,
			wantStorage: map[string]internal.Event{
				"qwerty-12345": {
					ID:           "qwerty-12345",
					Candidates:   []string{"c-1"},
					Participants: []string{"p-1"},
					Start:        first(time.Parse(time.RFC3339, "2022-09-02T11:00:00Z")),
					Finish:       first(time.Parse(time.RFC3339, "2022-09-02T15:00:00Z")),
					RepeatType:   0,
					Info:         internal.CustomEventInfo{},
				}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				events: tt.fields.events,
			}
			if err := s.AddEvent(tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("AddEvent() error = %v, wantErr %v", err, tt.wantErr)
				eq := reflect.DeepEqual(s.events, tt.wantStorage)
				t.Errorf("AddEvent() storage= %v, want %v", eq, true)
			}
		})
	}
}

//TODO: Add tests.
