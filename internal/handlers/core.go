package handlers

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"clothing-store/internal/models"

	"github.com/google/uuid"
)

type RoleHandler struct {
	DB *sql.DB
}

func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	role.ID = uuid.New().String()
	if _, err := h.DB.Exec("INSERT INTO roles (id, name) VALUES (?, ?)", role.ID, role.Name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func (h *RoleHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query("SELECT id, name FROM roles")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	roles := []models.Role{}
	for rows.Next() {
		var r models.Role
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			continue
		}
		roles = append(roles, r)
	}
	json.NewEncoder(w).Encode(roles)
}

func (h *RoleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var role models.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := h.DB.Exec("UPDATE roles SET name=? WHERE id=?", role.Name, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "роль не найдена", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, _ := h.DB.Exec("DELETE FROM roles WHERE id=?", id)
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "роль не найдена", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type CartHandler struct {
	DB *sql.DB
}

func (h *CartHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c models.Cart
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.ID = uuid.New().String()
	if c.Count <= 0 {
		c.Count = 1
	}
	if _, err := h.DB.Exec("INSERT INTO carts (id, user_id, product_id, count) VALUES (?, ?, ?, ?)", c.ID, c.UserID, c.ProductID, c.Count); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func (h *CartHandler) GetByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("user_id")
	rows, err := h.DB.Query("SELECT id, user_id, product_id, count FROM carts WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	carts := []models.Cart{}
	for rows.Next() {
		var c models.Cart
		if err := rows.Scan(&c.ID, &c.UserID, &c.ProductID, &c.Count); err != nil {
			continue
		}
		carts = append(carts, c)
	}
	json.NewEncoder(w).Encode(carts)
}

func (h *CartHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var c models.Cart
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := h.DB.Exec("UPDATE carts SET count=? WHERE id=?", c.Count, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "запись корзины не найдена", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *CartHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, _ := h.DB.Exec("DELETE FROM carts WHERE id=?", id)
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "запись корзины не найдена", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type PhotoHandler struct {
	DB        *sql.DB
	UploadDir string
}

func (h *PhotoHandler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "ошибка чтения формы", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "файл photo обязателен", http.StatusBadRequest)
		return
	}
	defer file.Close()
	productID := r.FormValue("product_id")
	if productID == "" {
		http.Error(w, "product_id обязателен", http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := uuid.New().String() + ext
	dst, err := os.Create(filepath.Join(h.UploadDir, filename))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	photo := models.Photo{
		ID:        uuid.New().String(),
		Filename:  filename,
		ProductID: productID,
	}
	if _, err := h.DB.Exec("INSERT INTO photos (id, filename, product_id) VALUES (?, ?, ?)", photo.ID, photo.Filename, photo.ProductID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photo)
}

func (h *PhotoHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var filename string
	err := h.DB.QueryRow("SELECT filename FROM photos WHERE id = ?", id).Scan(&filename)
	if err != nil {
		http.Error(w, "фото не найдено", http.StatusNotFound)
		return
	}
	http.ServeFile(w, r, filepath.Join(h.UploadDir, filename))
}

func (h *PhotoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var filename string
	err := h.DB.QueryRow("SELECT filename FROM photos WHERE id = ?", id).Scan(&filename)
	if err != nil {
		http.Error(w, "фото не найдено", http.StatusNotFound)
		return
	}
	os.Remove(filepath.Join(h.UploadDir, filename))
	h.DB.Exec("DELETE FROM photos WHERE id=?", id)
	w.WriteHeader(http.StatusNoContent)
}
