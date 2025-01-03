package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

type TaskStore struct {
	sync.RWMutex
	tasks map[string]Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: map[string]Task{
			"1": {
				ID:          "1",
				Description: "Сделать финальное задание темы REST API",
				Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
				Applications: []string{
					"VS Code",
					"Terminal",
					"git",
				},
			},
			"2": {
				ID:          "2",
				Description: "Протестировать финальное задание с помощью Postman",
				Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
				Applications: []string{
					"VS Code",
					"Terminal",
					"git",
					"Postman",
				},
			},
		},
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Ошибка при кодировании JSON: %v", err)
			http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
		}

	}
}

func (ts *TaskStore) getAllTasks(w http.ResponseWriter, r *http.Request) {
	ts.RLock()
	defer ts.RUnlock()

	tasksSlice := make([]Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		tasksSlice = append(tasksSlice, task)
	}
	sendJSONResponse(w, tasksSlice, http.StatusOK)
}

func (ts *TaskStore) createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		sendJSONResponse(w, ErrorResponse{Error: "Некорректный запрос"}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if newTask.ID == "" {
		sendJSONResponse(w, ErrorResponse{Error: "ID задачи не может быть пустым"}, http.StatusBadRequest)
		return
	}

	ts.Lock()
	defer ts.Unlock()

	if _, exists := ts.tasks[newTask.ID]; exists {
		sendJSONResponse(w, ErrorResponse{Error: "Задача с таким ID уже существует"}, http.StatusBadRequest)
		return
	}

	ts.tasks[newTask.ID] = newTask
	sendJSONResponse(w, newTask, http.StatusCreated)
}

func (ts *TaskStore) getTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ts.RLock()
	task, exists := ts.tasks[id]
	ts.RUnlock()

	if !exists {
		sendJSONResponse(w, ErrorResponse{Error: "Задача не найдена"}, http.StatusNotFound)
		return
	}
	sendJSONResponse(w, task, http.StatusOK)
}

func (ts *TaskStore) deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	ts.Lock()
	defer ts.Unlock()

	if _, exists := ts.tasks[id]; !exists {
		sendJSONResponse(w, ErrorResponse{Error: "Задача не найдена"}, http.StatusNotFound)
		return
	}

	delete(ts.tasks, id)
	sendJSONResponse(w, nil, http.StatusOK)
}

func main() {
	r := chi.NewRouter()
	taskStore := NewTaskStore()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", taskStore.getAllTasks)
		r.Post("/", taskStore.createTask)
		r.Get("/{id}", taskStore.getTaskByID)
		r.Delete("/{id}", taskStore.deleteTaskByID)
	})

	addr := ":8080"
	log.Printf("Сервер запущен на порту %s\n", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %s\n", err.Error())
	}
}
