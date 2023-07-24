package router

import (
	"net/http"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestHandler_GetProducts(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
		p httprouter.Params
	}
	tests := []struct {
		name   string
		router *Handler
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.router.GetProducts(tt.args.w, tt.args.r, tt.args.p)
		})
	}
}

func TestHandler_PostTableURL(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
		p httprouter.Params
	}
	tests := []struct {
		name   string
		router *Handler
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.router.PostTableURL(tt.args.w, tt.args.r, tt.args.p)
		})
	}
}
