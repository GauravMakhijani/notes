package main

import (
	"net/http"

	"github.com/GauravMakhijani/notes/internal/handler"
	"github.com/GauravMakhijani/notes/internal/middleware"
	"github.com/GauravMakhijani/notes/internal/service"
	"github.com/gorilla/mux"
)

func initRouter(service service.Service) *mux.Router {

	router := mux.NewRouter()
	router.Use(middleware.RateLimiter)

	//Auth router
	authRouter := router.PathPrefix("/api/auth").Subrouter()

	authRouter.HandleFunc("/signup", handler.SignUpHanler(service)).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", handler.LoginHandler(service)).Methods(http.MethodPost)
	authRouter.HandleFunc("/ping", middleware.SetMiddleWareAuthentication(PingHandler())).Methods(http.MethodGet)

	//Notes router
	notesRouter := router.PathPrefix("/api/notes").Subrouter()
	notesRouter.HandleFunc("", middleware.SetMiddleWareAuthentication(handler.CreateNoteHandler(service))).Methods(http.MethodPost)
	notesRouter.HandleFunc("", middleware.SetMiddleWareAuthentication(handler.ListNotesHandler(service))).Methods(http.MethodGet)
	notesRouter.HandleFunc("/{note_id}", middleware.SetMiddleWareAuthentication(handler.GetNoteByIDHandler(service))).Methods(http.MethodGet)
	notesRouter.HandleFunc("/{note_id}", middleware.SetMiddleWareAuthentication(handler.DeleteNoteHandler(service))).Methods(http.MethodDelete)
	notesRouter.HandleFunc("/{note_id}", middleware.SetMiddleWareAuthentication(handler.UpdateNoteHandler(service))).Methods(http.MethodPut)
	notesRouter.HandleFunc("/{note_id}/share", middleware.SetMiddleWareAuthentication(handler.ShareNoteHandler(service))).Methods(http.MethodPost)

	//Search router
	router.HandleFunc("/api/search", middleware.SetMiddleWareAuthentication(handler.SearchNotesHandler(service))).Methods(http.MethodGet)
	return router
}

func PingHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	}
}
