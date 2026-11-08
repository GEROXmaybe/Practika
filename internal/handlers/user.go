package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"clothing-store/internal/models"

	"github.com/google/uuid"
)

type UserHandler struct {
	DB *sql.DB
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if u.Name == "" || u.Email == "" {
		http.Error(w, "имя и email обязательны", http.StatusBadRequest)
		return
	}
	u.ID = uuid.New().String()
	u.RegistrationDate = time.Now().Format(time.RFC3339)

	_, err := h.DB.Exec(`INSERT INTO users (id, name, email, phone, role_id, registration_date) VALUES (?, ?, ?, ?, ?, ?)`,
		u.ID, u.Name, u.Email, u.Phone, u.RoleID, u.RegistrationDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit := 10
	offset := 0
	if l := q.Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			limit = v
		}
	}
	if o := q.Get("offset"); o != "" {
		if v, err := strconv.Atoi(o); err == nil && v >= 0 {
			offset = v
		}
	}

	rows, err := h.DB.Query("SELECT id, name, email, phone, role_id, registration_date FROM users LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.RoleID, &u.RegistrationDate); err != nil {
			continue
		}
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetByRole(w http.ResponseWriter, r *http.Request) {
	roleID := r.PathValue("id")
	rows, err := h.DB.Query("SELECT id, name, email, phone, role_id, registration_date FROM users WHERE role_id = ?", roleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	users := []models.User{}
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Phone, &u.RoleID, &u.RegistrationDate); err != nil {
			continue
		}
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := h.DB.Exec(`UPDATE users SET name=?, email=?, phone=?, role_id=? WHERE id=?`, u.Name, u.Email, u.Phone, u.RoleID, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "пользователь не найден", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, _ := h.DB.Exec("DELETE FROM users WHERE id = ?", id)
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "пользователь не найден", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
