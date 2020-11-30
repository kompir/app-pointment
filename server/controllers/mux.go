package controllers

import (
	"context"
	"net/http"
	"regexp"
	"strings"

	"github.com/kompir/app-pointment/server/models"
	"github.com/kompir/app-pointment/server/transport"
)

const paramsKey = "ps"

type ctxKey string

type urlParam struct {
	name     string
	regEx    string
	value    string
	position int
}

type route struct {
	path    string
	method  string
	params  map[string]urlParam
	handler http.Handler
}

func (r *route) populate(req *http.Request) string {
	urlSlice := splitURL(req.URL.Path)
	pathSlice := splitURL(r.path)
	if len(pathSlice) != len(urlSlice) {
		return ""
	}
	for name, param := range r.params {
		regexParamVal := urlSlice[param.position]
		regex := regexp.MustCompile(param.regEx)
		if name != "" && regex.MatchString(regexParamVal) {
			param.value = regexParamVal
			r.params[name] = param
			pathSlice[param.position] = regexParamVal
		}
	}
	pathStr := "/" + strings.Join(pathSlice, "/")
	if req.URL.Path == pathStr {
		return r.method + pathStr
	}
	return ""
}

type RegexpMux struct {
	routes    []*route
	routesMap map[string]*route
}

func (h *RegexpMux) Get(pattern string, handler http.Handler) {
	h.Handle(http.MethodGet, pattern, handler)
}

func (h *RegexpMux) Post(pattern string, handler http.Handler) {
	h.Handle(http.MethodPost, pattern, handler)
}

func (h *RegexpMux) Patch(pattern string, handler http.Handler) {
	h.Handle(http.MethodPatch, pattern, handler)
}

func (h *RegexpMux) Put(pattern string, handler http.Handler) {
	h.Handle(http.MethodPut, pattern, handler)
}

func (h *RegexpMux) Delete(pattern string, handler http.Handler) {
	h.Handle(http.MethodDelete, pattern, handler)
}

func (h *RegexpMux) Handle(method, pattern string, handler http.Handler) {
	ps := h.params(pattern)
	r := &route{
		method:  method,
		path:    pattern,
		params:  ps,
		handler: handler,
	}
	h.routes = append(h.routes, r)
}

func (h RegexpMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.routesMap = map[string]*route{}
	for _, route := range h.routes {
		key := route.populate(r)
		h.routesMap[key] = route
	}
	key := r.Method + r.URL.Path
	route, ok := h.routesMap[key]
	if !ok {
		transport.SendError(w, models.NotFoundError{})
		return
	}
	ctx := r.Context()
	if len(route.params) != 0 {
		ctx = context.WithValue(ctx, ctxKey(paramsKey), route.params)
	}
	route.handler.ServeHTTP(w, r.WithContext(ctx))
}

func (h RegexpMux) params(url string) map[string]urlParam {
	ps := map[string]urlParam{}
	for _, v := range splitURL(url) {
		p := h.parseParam(url, v)
		if p.name != "" {
			ps[p.name] = p
		}
	}
	return ps
}

func (h RegexpMux) parseParam(url, regexParam string) urlParam {
	r := regexp.MustCompile(`({[a-z]+}:)(.+)`)
	matches := r.FindStringSubmatch(regexParam)
	if len(matches) < 3 {
		return urlParam{
			regEx: ".+",
		}
	}
	replacer := strings.NewReplacer(
		"{", "",
		"}", "",
		":", "",
	)
	name := replacer.Replace(matches[1])
	regEx := matches[2]
	var position int
	for i, v := range splitURL(url) {
		if v == matches[1]+matches[2] {
			position = i
		}
	}
	return urlParam{
		name:     name,
		regEx:    regEx,
		position: position,
	}
}

func splitURL(s string) []string {
	var res []string
	for _, p := range strings.Split(strings.TrimSpace(s), "/") {
		if strings.TrimSpace(p) != "" {
			res = append(res, p)
		}
	}
	return res
}
