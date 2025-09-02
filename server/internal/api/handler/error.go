package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dorianneto/bugfy/internal/api/model"
	service "github.com/dorianneto/bugfy/internal/service"
	"github.com/dorianneto/bugfy/util"
)

type ErrorHandler struct {
	errorService *service.ErrorService
}

func NewErrorHandler(errorService *service.ErrorService) *ErrorHandler {
	return &ErrorHandler{
		errorService: errorService,
	}
}

func (h *ErrorHandler) CreateError(w http.ResponseWriter, r *http.Request) {
	var req model.RequestCreateError
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateError - JSON decode error: %v", err)
		util.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	log.Printf("CreateError - Request received: projectID=%s", req.ProjectID)

	e, err := h.errorService.CreateError(r.Context(), req)
	if err != nil {
		log.Printf("CreateError - Service error: %v", err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("CreateError - Success: error created with ID=%s, projectID=%s", e.ID, e.ProjectID)

	util.WriteJSON(w, http.StatusCreated, e)
}
