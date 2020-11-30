package controllers

import (
	"net/http"

	"github.com/kompir/app-pointment/server/middleware"
)

const (
	idParamName  = "id"
	idsParamName = "ids"
	idParam      = `{` + idParamName + `}:^[0-9]+$`
	idsParam     = `{` + idsParamName + `}:[0-9]+(,[0-9]+)*`
)

type RemindersService interface {
	creator
	editor
	lister
	deleter
}

type RouterConfig struct {
	Service RemindersService
}

func NewRouter(cfg RouterConfig) http.Handler {
	r := RegexpMux{}
	m := middleware.New(
		middleware.HTTPLogger,
	)
	r.Get("/health", m.Then(health()))
	r.Post("/reminders", m.Then(createReminder(cfg.Service)))
	r.Get("/reminders/"+idsParam, m.Then(listReminders(cfg.Service)))
	r.Delete("/reminders/"+idsParam, m.Then(deleteReminders(cfg.Service)))
	r.Patch("/reminders/"+idParam, m.Then(editReminder(cfg.Service)))
	return r
}
