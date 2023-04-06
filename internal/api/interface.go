package api

type Service interface {
	CreateUser(body []byte) ([]byte, error)
	CreateEventWithUsers(body []byte) ([]byte, error)
	GetEventDetails(body []byte) ([]byte, error)
	AcceptInvitation(body []byte) error
	RejectInvitation(body []byte) error
	GetEvents(body []byte) ([]byte, error)
	FindSlot(body []byte) ([]byte, error)
}
