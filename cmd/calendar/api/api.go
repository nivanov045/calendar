package api

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Service interface {
	CreateUser(body []byte) ([]byte, error)
	CreateEventWithUsers(body []byte) ([]byte, error)
	GetEventDetails(body []byte) ([]byte, error)
	AcceptInvitation(body []byte) error
	RejectInvitation(body []byte) error
	GetEvents(body []byte) ([]byte, error)
	FindSlot(body []byte) ([]byte, error)
}

type api struct {
	service Service
}

func New(service Service) *api {
	return &api{service: service}
}

func (a *api) Run(address string) error {
	log.Println("api::Run::info: started with addr:", address)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/create-user/", a.createUserHandler)
	r.Post("/create-event-with-users/", a.createEventWithUsersHandler)
	r.Get("/event-details/", a.getEventDetailsHandler)
	r.Post("/accept-invitation/", a.acceptInvitationHandler)
	r.Post("/reject-invitation/", a.rejectInvitationHandler)
	r.Get("/events/", a.getEventsHandler)
	r.Get("/find-slot/", a.findSlotHandler)

	return http.ListenAndServe(address, r)
}

func (a *api) createUserHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("api::createUserHandler::info: started")
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("api::createUserHandler::warning: can't read request body with:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	resp, err := a.service.CreateUser(requestBody)
	if err != nil {
		log.Println("api::createUserHandler::warning: in user creation:", err)
		if err.Error() == "wrong query" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write([]byte("{}"))
		return
	}
	log.Println("api::createUserHandler::info: user created")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (a *api) createEventWithUsersHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("api::createEventWithUsersHandler::info: started")
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("api::createEventWithUsersHandler::warning: can't read request body with:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	resp, err := a.service.CreateEventWithUsers(requestBody)
	if err != nil {
		log.Println("api::createEventWithUsersHandler::warning: in event creation:", err)
		if err.Error() == "wrong query" || err.Error() == "wrong repeat type" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write([]byte("{}"))
		return
	}
	log.Println("api::createEventWithUsersHandler::info: event created")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (a *api) getEventDetailsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("api::getEventDetailsHandler::info: started")
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("api::getEventDetailsHandler::warning: can't read request body with:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	resp, err := a.service.GetEventDetails(requestBody)
	if err != nil {
		log.Println("api::getEventDetailsHandler::warning:", err)
		if err.Error() == "wrong query" || err.Error() == "unable to get event" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write([]byte("{}"))
		return
	}
	log.Println("api::getEventDetailsHandler::info: get details")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (a *api) acceptInvitationHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("api::acceptInvitationHandler::info: started")
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("api::acceptInvitationHandler::warning: can't read request body with:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	if err := a.service.AcceptInvitation(requestBody); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("api::acceptInvitationHandler::warning:", err)
		if err.Error() == "wrong query" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	w.Write([]byte("{}"))
}

func (a *api) rejectInvitationHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("api::rejectInvitationHandler::info: started")
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("api::rejectInvitationHandler::warning: can't read request body with:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	if err := a.service.RejectInvitation(requestBody); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		log.Println("api::rejectInvitationHandler::warning:", err)
		if err.Error() == "wrong query" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	w.Write([]byte("{}"))
}

func (a *api) getEventsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("api::getEventsHandler::info: started")
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("api::getEventsHandler::warning: can't read request body with:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	resp, err := a.service.GetEvents(requestBody)
	if err != nil {
		if err.Error() == "wrong query" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write([]byte("{}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (a *api) findSlotHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("api::findSlotHandler::info: started")
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("api::findSlotHandler::warning: can't read request body with:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	resp, err := a.service.FindSlot(requestBody)
	if err != nil {
		if err.Error() == "wrong query" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write([]byte("{}"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
