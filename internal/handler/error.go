package handler

import (
	"edgex-club/internal/core"
	"net/http"
)

type ErrorMsg struct {
	Body string
}

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	t := core.TemplateStore["errorTpl"]
	t.Execute(w, nil)
}
