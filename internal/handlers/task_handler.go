package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"tasks-app/internal/storage"
	"tasks-app/internal/utils"
)

type TaskHandler struct {
	storage *storage.MemoryStorage
}

func NewTaskHandler(storage *storage.MemoryStorage) *TaskHandler {
	return &TaskHandler{
		storage: storage,
	}
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (h *TaskHandler) Health(w http.ResponseWriter, r *http.Request) {
	utils.JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		tasks := h.storage.GetAll()
		utils.JSON(w, http.StatusOK, tasks)

	case http.MethodPost:
		var req CreateTaskRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.JSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid body",
			})
			return
		}

		task := h.storage.Create(req.Title, req.Description)

		utils.JSON(w, http.StatusCreated, task)

	default:
		utils.JSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
	}
}

func (h *TaskHandler) TaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		utils.JSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid id",
		})
		return
	}

	switch r.Method {

	case http.MethodGet:
		task, err := h.storage.GetByID(id)

		if err != nil {
			utils.JSON(w, http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}

		utils.JSON(w, http.StatusOK, task)

	case http.MethodPut:
		var req UpdateTaskRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.JSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid body",
			})
			return
		}

		task, err := h.storage.Update(
			id,
			req.Title,
			req.Description,
			req.Done,
		)

		if err != nil {
			utils.JSON(w, http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}

		utils.JSON(w, http.StatusOK, task)

	case http.MethodDelete:
		err := h.storage.Delete(id)

		if err != nil {
			utils.JSON(w, http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}

		utils.JSON(w, http.StatusOK, map[string]string{
			"message": "task deleted",
		})

	default:
		utils.JSON(w, http.StatusMethodNotAllowed, map[string]string{
			"error": "method not allowed",
		})
	}
}
