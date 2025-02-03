package http

import (
	"example.com/app"
	"net/http"
)

type userHandler struct {
	userService   app.UserService
	generateToken func(userID int) (string, error)
}

func NewUserHandler(userService app.UserService, generateToken func(int) (string, error)) *userHandler {
	return &userHandler{
		userService:   userService,
		generateToken: generateToken,
	}
}

func (h *userHandler) register(w http.ResponseWriter, r *http.Request) {
	var req app.UserCredentials
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	user, err := h.userService.Create(&req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *userHandler) login(w http.ResponseWriter, r *http.Request) {
	var req app.UserCredentials
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	user, err := h.userService.Authenticate(&req)
	if err != nil {
		writeError(w, err)
		return
	}

	token, err := h.generateToken(user.ID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user})
}

func (h *userHandler) getUser(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	user, err := h.userService.Get(userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *userHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var req app.UpdateUser
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	user, err := h.userService.Update(userID, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *userHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	if err := h.userService.Delete(userID); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusNoContent, struct{}{})
}
