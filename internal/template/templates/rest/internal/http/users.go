package http

import (
	"net/http"

	"example.com/rest/internal/domain"
	"example.com/rest/internal/services"
)

type UserHandler struct {
	*baseHandler
	userService   *services.UserService
	generateToken func(userID int) (string, error)
}

func NewUserHandler(baseHandler *baseHandler, userService *services.UserService, generateToken func(int) (string, error)) *UserHandler {
	return &UserHandler{
		baseHandler:   baseHandler,
		userService:   userService,
		generateToken: generateToken,
	}
}

func (h *UserHandler) register(w http.ResponseWriter, r *http.Request) {
	var req domain.UserCredentials
	if err := h.json.Read(r, &req); err != nil {
		h.json.WriteError(w, r, err)
		return
	}

	user, err := h.userService.Create(r.Context(), &req)
	if err != nil {
		h.json.WriteError(w, r, err)
		return
	}

	h.json.Write(w, http.StatusCreated, map[string]any{"user": user})
}

func (h *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	var req domain.UserCredentials
	if err := h.json.Read(r, &req); err != nil {
		h.json.WriteError(w, r, err)
		return
	}

	user, err := h.userService.Authenticate(r.Context(), &req)
	if err != nil {
		h.json.WriteError(w, r, err)
		return
	}

	token, err := h.generateToken(user.ID)
	if err != nil {
		h.json.WriteError(w, r, err)
		return
	}
	h.json.Write(w, http.StatusOK, map[string]any{"token": token, "user": user})
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.json.WriteError(w, r, err)
		return
	}

	user, err := h.userService.GetByID(r.Context(), userID)
	if err != nil {
		h.json.WriteError(w, r, err)
		return
	}

	h.json.Write(w, http.StatusOK, map[string]any{"user": user})
}

func (h *UserHandler) updateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.json.WriteError(w, r, err)
		return
	}

	var req domain.UserPatch
	if err := h.json.Read(r, &req); err != nil {
		h.json.WriteError(w, r, err)
		return
	}

	user, err := h.userService.Update(r.Context(), userID, &req)
	if err != nil {
		h.json.WriteError(w, r, err)
		return
	}

	h.json.Write(w, http.StatusOK, map[string]any{"user": user})
}

func (h *UserHandler) deleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserID(r)
	if err != nil {
		h.json.WriteError(w, r, err)
		return
	}

	if err := h.userService.Delete(r.Context(), userID); err != nil {
		h.json.WriteError(w, r, err)
		return
	}

	h.json.Write(w, http.StatusNoContent, nil)
}
