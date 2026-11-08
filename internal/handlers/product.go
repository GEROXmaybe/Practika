package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"clothing-store/internal/models"

	"github.com/google/uuid"
)

type ProductHandler struct {
	DB *sql.DB
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if p.Name == "" {
		http.Error(w, "название товара обязательно", http.StatusBadRequest)
		return
	}
	p.ID = uuid.New().String()

	_, err := h.DB.Exec(`INSERT INTO products (id, name, type, gender, size, color, description, brand, season, price) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.ID, p.Name, p.Type, p.Gender, p.Size, p.Color, p.Description, p.Brand, p.Season, p.Price)
	if err != nil {
		slog.Error("ошибка добавления товара", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
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

	conditions := []string{"1=1"}
	args := []interface{}{}

	filters := map[string]string{
		"type": "type", "gender": "gender", "size": "size",
		"color": "color", "brand": "brand", "season": "season",
	}

	for param, col := range filters {
		if val := q.Get(param); val != "" {
			conditions = append(conditions, col+" = ?")
			args = append(args, val)
		}
	}

	orderBy := "id ASC"
	if q.Get("sort_by") == "price" {
		order := "ASC"
		if strings.ToLower(q.Get("order")) == "desc" {
			order = "DESC"
		}
		orderBy = "price " + order
	}

	query := "SELECT id, name, type, gender, size, color, description, brand, season, price FROM products WHERE " +
		strings.Join(conditions, " AND ") + " ORDER BY " + orderBy + " LIMIT ? OFFSET ?"

	args = append(args, limit, offset)

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Type, &p.Gender, &p.Size, &p.Color, &p.Description, &p.Brand, &p.Season, &p.Price); err != nil {
			continue
		}
		products = append(products, p)
	}
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var p models.Product
	err := h.DB.QueryRow(`SELECT id, name, type, gender, size, color, description, brand, season, price FROM products WHERE id = ?`, id).
		Scan(&p.ID, &p.Name, &p.Type, &p.Gender, &p.Size, &p.Color, &p.Description, &p.Brand, &p.Season, &p.Price)
	if err != nil {
		http.Error(w, "товар не найден", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.DB.Exec(`UPDATE products SET name=?, type=?, gender=?, size=?, color=?, description=?, brand=?, season=?, price=? WHERE id=?`,
		p.Name, p.Type, p.Gender, p.Size, p.Color, p.Description, p.Brand, p.Season, p.Price, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "товар не найден", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, _ := h.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "товар не найден", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
