package http

import (
	"example.com/app"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type noteHandler struct {
	noteService app.NoteService
}

func NewNoteHandler(noteService app.NoteService) *noteHandler {
	return &noteHandler{noteService: noteService}
}

func (h *noteHandler) createNote(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	var req app.CreateNote
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	note, err := h.noteService.Create(userID, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, note)
}

func (h *noteHandler) getUserNotes(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	notes, err := h.noteService.GetAll(userID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, notes)
}

func (h *noteHandler) updateNote(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	noteID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, app.Errorf(app.INVALID_ERROR, "invalid note ID"))
		return
	}

	var req app.UpdateNote
	if err := decodeJSON(r, &req); err != nil {
		writeError(w, err)
		return
	}

	note, err := h.noteService.Update(userID, noteID, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, note)
}

func (h *noteHandler) deleteNote(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		writeError(w, err)
		return
	}

	noteID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeError(w, app.Errorf(app.INVALID_ERROR, "invalid note ID"))
		return
	}

	if err := h.noteService.Delete(userID, noteID); err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, struct{}{})
}
