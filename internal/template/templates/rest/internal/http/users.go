package http

import (
	"net/http"

	"example.com/rest/internal/domain"
)

type userService interface {
	Create(req *domain.UserCredentials) (*domain.User, error)
	Get(id int) (*domain.User, error)
	GetByUsername(username string) (*domain.User, error)
	Authenticate(req *domain.UserCredentials) (*domain.User, error)
	Update(id int, req *domain.UpdateUser) (*domain.User, error)
	Delete(id int) error
}

type userHandler struct {
	userService   userService
	generateToken func(userID int) (string, error)
}

func NewUserHandler(userService userService, generateToken func(int) (string, error)) *userHandler {
	return &userHandler{
		userService:   userService,
		generateToken: generateToken,
	}
}

func (h *userHandler) register(w http.ResponseWriter, r *http.Request) {
	var req domain.UserCredentials
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
	var req domain.UserCredentials
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

	var req domain.UpdateUser
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
