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

type OrderHandler struct {
	DB *sql.DB
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var o models.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if o.UserID == "" {
		http.Error(w, "user_id обязателен", http.StatusBadRequest)
		return
	}
	o.ID = uuid.New().String()
	o.DateTime = time.Now().Format(time.RFC3339)
	if o.Status == "" {
		o.Status = "Принят"
	}
	if !models.ValidOrderStatuses[o.Status] {
		http.Error(w, "недопустимый статус заказа", http.StatusBadRequest)
		return
	}

	_, err := h.DB.Exec(`INSERT INTO orders (id, user_id, datetime, status) VALUES (?, ?, ?, ?)`,
		o.ID, o.UserID, o.DateTime, o.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(o)
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
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

	if status := q.Get("status"); status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}
	if from := q.Get("from"); from != "" {
		conditions = append(conditions, "datetime >= ?")
		args = append(args, from)
	}
	if to := q.Get("to"); to != "" {
		conditions = append(conditions, "datetime <= ?")
		args = append(args, to)
	}

	query := "SELECT id, user_id, datetime, status FROM orders WHERE " +
		joinConditions(conditions) + " ORDER BY datetime DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.DateTime, &o.Status); err != nil {
			continue
		}
		orders = append(orders, o)
	}
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var o models.Order
	err := h.DB.QueryRow(`SELECT id, user_id, datetime, status FROM orders WHERE id = ?`, id).
		Scan(&o.ID, &o.UserID, &o.DateTime, &o.Status)
	if err != nil {
		http.Error(w, "заказ не найден", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(o)
}

func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var o models.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !models.ValidOrderStatuses[o.Status] {
		http.Error(w, "недопустимый статус заказа", http.StatusBadRequest)
		return
	}
	res, err := h.DB.Exec(`UPDATE orders SET status=? WHERE id=?`, o.Status, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if n, _ := res.RowsAffected(); n == 0 {
		http.Error(w, "заказ не найден", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

type OrderItemHandler struct {
	DB *sql.DB
}

func (h *OrderItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var oi models.OrderItem
	if err := json.NewDecoder(r.Body).Decode(&oi); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	oi.ID = uuid.New().String()
	_, err := h.DB.Exec(`INSERT INTO order_items (id, order_id, product_id, price) VALUES (?, ?, ?, ?)`,
		oi.ID, oi.OrderID, oi.ProductID, oi.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(oi)
}

func (h *OrderItemHandler) GetByOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.PathValue("order_id")
	rows, err := h.DB.Query("SELECT id, order_id, product_id, price FROM order_items WHERE order_id = ?", orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	items := []models.OrderItem{}
	for rows.Next() {
		var oi models.OrderItem
		if err := rows.Scan(&oi.ID, &oi.OrderID, &oi.ProductID, &oi.Price); err != nil {
			continue
		}
		items = append(items, oi)
	}
	json.NewEncoder(w).Encode(items)
}

func joinConditions(conds []string) string {
	result := ""
	for i, c := range conds {
		if i > 0 {
			result += " AND "
		}
		result += c
	}
	return result
}
