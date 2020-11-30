package controllers

import (
	"net/http"

	"github.com/kompir/app-pointment/server/models"
	"github.com/kompir/app-pointment/server/transport"
)

type lister interface {
	List(ids []int) ([]models.Reminder, error)
}

func listReminders(service lister) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids, err := parseIDsParam(r.Context())
		if err != nil {
			transport.SendError(w, err)
			return
		}
		reminders, err := service.List(ids)
		if err != nil {
			transport.SendError(w, err)
			return
		}
		transport.SendJSON(w, reminders, http.StatusOK)
	})
}
