package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dorianneto/bugfy/internal/api/model"
	service "github.com/dorianneto/bugfy/internal/service"
	"github.com/dorianneto/bugfy/util"
)

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req model.RequestCreateProject
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateProject - JSON decode error: %v", err)
		util.WriteError(w, http.StatusBadRequest, "invalid JSON payload")
		return
	}

	log.Printf("CreateProject - Request received: title=%s", req.Title)

	project, err := h.projectService.CreateProject(r.Context(), req)
	if err != nil {
		log.Printf("CreateProject - Service error: %v", err)
		util.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("CreateProject - Success: project created with ID=%s, title=%s", project.ID, project.Title)

	util.WriteJSON(w, http.StatusCreated, project)
}
