package server

import (
	"errors"
	"github.com/go-chi/httplog/v2"
	"github.com/google/uuid"
	"github.com/mrityunjaygr8/guzei/store"
	"net/http"
	"strconv"
)

type UserInsertParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

func (a *Server) UserInsert(w http.ResponseWriter, r *http.Request) {
	var req UserInsertParams
	ok := a.ReadJSON(w, r, &req)
	if !ok {
		return
	}

	user, err := a.store.UserInsert(req.Email, req.Password, uuid.New().String(), req.Admin)
	if err != nil {
		if errors.Is(err, store.ErrUserExists) {
			a.JsonError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
		if errors.Is(err, store.ErrStoreError) {
			a.logger.Logger.Info("woo woo")
			//slog.Info("asdasd")
		}
		a.JsonError(w, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	a.WriteJSON(w, http.StatusOK, user)
}

func (a *Server) UserList(w http.ResponseWriter, r *http.Request) {
	pageSize := 10
	pageNumber := 1
	var err error
	query := r.URL.Query()
	oplog := httplog.LogEntry(r.Context())
	oplog.Info("woo woo", "poo poo", "doo doo")
	if query.Has("pageSize") {
		pageSize, err = strconv.Atoi(query.Get("pageSize"))
		if err != nil {
			a.JsonError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
	}
	if query.Has("pageNumber") {
		pageNumber, err = strconv.Atoi(query.Get("pageNumber"))
		if err != nil {
			a.JsonError(w, http.StatusBadRequest, err.Error(), nil)
			return
		}
	}
	users, err := a.store.UserList(pageNumber, pageSize)
	if err != nil {
		a.JsonError(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	a.WriteJSON(w, http.StatusOK, users)
}
