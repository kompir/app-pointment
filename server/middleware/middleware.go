package middleware

import (
	"net/http"
)

func New(ms ...func(h http.Handler) http.Handler) *Middleware {
	return &Middleware{
		functions: ms,
	}
}

type Middleware struct {
	functions []func(h http.Handler) http.Handler
}

func (m *Middleware) Then(h http.Handler) http.Handler {
	if h == nil {
		h = http.DefaultServeMux
	}
	for i := range m.functions {
		h = m.functions[len(m.functions)-1-i](h)
	}
	return h
}
