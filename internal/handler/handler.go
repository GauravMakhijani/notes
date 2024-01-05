package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GauravMakhijani/notes/internal/domain"
	"github.com/GauravMakhijani/notes/internal/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Response struct {
	Data         interface{} `json:"data,omitempty"`
	ErrorMessage string      `json:"error,omitempty"`
	ErrorCode    int64       `json:"error_code,omitempty"`
}

// SuccessResponse encodes the provided response in JSON format
func SuccessResponse(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	w.WriteHeader(statusCode)

	// For 204 status code, no response body should be sent
	if statusCode == http.StatusNoContent {
		return
	}

	response := Response{
		Data: data,
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		WriteServerErrorResponse(ctx, w)
		return
	}
}

// WriteServerErrorResponse ...
func WriteServerErrorResponse(ctx context.Context, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	_, err := w.Write([]byte(fmt.Sprintf("{\"message\":%s}", "internal server error")))
	if err != nil {
		logrus.Error("Error writing server error response", err)
	}
}

func SignUpHanler(service service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var signupReq domain.SignupRequest
		if err := json.NewDecoder(r.Body).Decode(&signupReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := service.CreateNewUser(r.Context(), signupReq)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		SuccessResponse(r.Context(), w, http.StatusCreated, map[string]interface{}{"message": "User created successfully"})
		w.WriteHeader(http.StatusOK)
	}

}

func LoginHandler(service service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var loginReq domain.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		loginResponse, err := service.LoginUser(r.Context(), loginReq)
		if err != nil {
			http.Error(w, "Failed to login", http.StatusInternalServerError)
			return
		}
		SuccessResponse(r.Context(), w, http.StatusOK, loginResponse)
		w.WriteHeader(http.StatusOK)
	}

}

func CreateNoteHandler(service service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var noteReq domain.NoteRequest
		if err := json.NewDecoder(r.Body).Decode(&noteReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		note, err := service.CreateNote(r.Context(), noteReq)
		if err != nil {
			http.Error(w, "Failed to create note", http.StatusInternalServerError)
			return
		}
		SuccessResponse(r.Context(), w, http.StatusCreated, note)
		w.WriteHeader(http.StatusOK)
	}

}

func GetNoteByIDHandler(service service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the id from the url
		noteID := mux.Vars(r)["note_id"]

		note, err := service.GetNoteByID(r.Context(), noteID)
		if err != nil {
			http.Error(w, "Failed to get note", http.StatusInternalServerError)
			return
		}
		SuccessResponse(r.Context(), w, http.StatusOK, note)
		w.WriteHeader(http.StatusOK)
	}

}

func ListNotesHandler(service service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		notes, err := service.ListNotes(r.Context())
		if err != nil {
			http.Error(w, "Failed to list notes", http.StatusInternalServerError)
			return
		}
		SuccessResponse(r.Context(), w, http.StatusOK, notes)
		w.WriteHeader(http.StatusOK)
	}

}

func DeleteNoteHandler(service service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the id from the url
		noteID := mux.Vars(r)["note_id"]
		err := service.DeleteNoteByID(r.Context(), noteID)
		if err != nil {
			http.Error(w, "Failed to delete note", http.StatusInternalServerError)
			return
		}
		SuccessResponse(r.Context(), w, http.StatusOK, map[string]interface{}{"message": "Note deleted successfully"})
		w.WriteHeader(http.StatusOK)
	}

}

func UpdateNoteHandler(service service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the id from the url
		noteID := mux.Vars(r)["note_id"]

		// Parse the request body
		var noteReq domain.NoteRequest
		if err := json.NewDecoder(r.Body).Decode(&noteReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if noteReq.Title == "" || noteReq.Body == "" {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		note, err := service.UpdateNoteByID(r.Context(), noteID, noteReq)
		if err != nil {
			http.Error(w, "Failed to update note", http.StatusInternalServerError)
			return
		}
		SuccessResponse(r.Context(), w, http.StatusOK, note)
		w.WriteHeader(http.StatusOK)
	}

}

func ShareNoteHandler(service service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the id from the url
		noteID := mux.Vars(r)["note_id"]

		// Parse the request body
		var shareReq domain.SharedNoteRequest
		if err := json.NewDecoder(r.Body).Decode(&shareReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err := service.ShareNoteWithUser(r.Context(), noteID, shareReq)
		if err != nil {
			http.Error(w, "Failed to share note", http.StatusInternalServerError)
			return
		}
		SuccessResponse(r.Context(), w, http.StatusOK, map[string]interface{}{"message": "Note shared successfully"})
		w.WriteHeader(http.StatusOK)
	}

}

func SearchNotesHandler(service service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the id from the url
		searchTerm := r.URL.Query().Get("q")
		notes, err := service.SearchNotes(r.Context(), searchTerm)
		if err != nil {
			http.Error(w, "Failed to search notes", http.StatusInternalServerError)
			return
		}
		SuccessResponse(r.Context(), w, http.StatusOK, notes)
		w.WriteHeader(http.StatusOK)
	}
}
