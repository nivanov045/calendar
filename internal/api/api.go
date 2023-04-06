package api

import (
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

type api struct {
	service Service
}

func New(service Service) *api {
	return &api{service: service}
}

func (a *api) Run(address string) error {
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
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Stack()
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	resp, err := a.service.CreateUser(requestBody)
	if err != nil {
		log.Error().Err(err).Stack()
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

func (a *api) createEventWithUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Stack()
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	resp, err := a.service.CreateEventWithUsers(requestBody)
	if err != nil {
		log.Error().Err(err).Stack()
		if err.Error() == "wrong query" || err.Error() == "wrong repeat type" {
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

func (a *api) getEventDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Stack()
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	resp, err := a.service.GetEventDetails(requestBody)
	if err != nil {
		log.Error().Err(err).Stack()
		if err.Error() == "wrong query" || err.Error() == "unable to get event" {
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

func (a *api) acceptInvitationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Stack()
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	if err := a.service.AcceptInvitation(requestBody); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		log.Error().Err(err).Stack()
		if err.Error() == "wrong query" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	w.Write([]byte("{}"))
}

func (a *api) rejectInvitationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Stack()
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{}"))
		return
	}
	if err := a.service.RejectInvitation(requestBody); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		log.Error().Err(err).Stack()
		if err.Error() == "wrong query" {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	w.Write([]byte("{}"))
}

func (a *api) getEventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Stack()
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
	w.Header().Set("content-type", "application/json")
	defer r.Body.Close()
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error().Err(err).Stack()
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

//TODO: Add tests.
