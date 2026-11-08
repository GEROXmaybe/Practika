package models

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Gender      string  `json:"gender"`
	Size        string  `json:"size"`
	Color       string  `json:"color"`
	Description string  `json:"description"`
	Brand       string  `json:"brand"`
	Season      string  `json:"season"`
	Price       float64 `json:"price"`
}

type Photo struct {
	ID        string `json:"id"`
	Filename  string `json:"filename"`
	ProductID string `json:"product_id"`
}

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	Phone            string `json:"phone"`
	RoleID           string `json:"role_id"`
	RegistrationDate string `json:"registration_date"`
}

type Cart struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	Count     int    `json:"count"`
}

type Order struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	DateTime string `json:"datetime"`
	Status   string `json:"status"`
}

type OrderItem struct {
	ID        string  `json:"id"`
	OrderID   string  `json:"order_id"`
	ProductID string  `json:"product_id"`
	Price     float64 `json:"price"`
}

var ValidOrderStatuses = map[string]bool{
	"Принят":         true,
	"В работе":       true,
	"Собран":         true,
	"Готов к выдаче": true,
	"Отменён":        true,
}
