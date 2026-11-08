package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"clothing-store/internal/config"
	"clothing-store/internal/database"
	"clothing-store/internal/handlers"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg := config.Load()

	if err := os.MkdirAll(cfg.UploadDir, os.ModePerm); err != nil {
		slog.Error("не удалось создать папку загрузок", "err", err)
		return
	}

	db, err := database.Init(cfg.DBPath)
	if err != nil {
		slog.Error("не удалось инициализировать БД", "err", err)
		return
	}
	defer db.Close()

	prodH := &handlers.ProductHandler{DB: db}
	userH := &handlers.UserHandler{DB: db}
	orderH := &handlers.OrderHandler{DB: db}
	orderItemH := &handlers.OrderItemHandler{DB: db}
	roleH := &handlers.RoleHandler{DB: db}
	cartH := &handlers.CartHandler{DB: db}
	photoH := &handlers.PhotoHandler{DB: db, UploadDir: cfg.UploadDir}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /product", prodH.Create)
	mux.HandleFunc("GET /product", prodH.List)
	mux.HandleFunc("GET /product/{id}", prodH.GetByID)
	mux.HandleFunc("PUT /product/{id}", prodH.Update)
	mux.HandleFunc("DELETE /product/{id}", prodH.Delete)

	mux.HandleFunc("POST /photo", photoH.Upload)
	mux.HandleFunc("GET /photo/{id}", photoH.Get)
	mux.HandleFunc("DELETE /photo/{id}", photoH.Delete)

	mux.HandleFunc("POST /role", roleH.Create)
	mux.HandleFunc("GET /role", roleH.List)
	mux.HandleFunc("PUT /role/{id}", roleH.Update)
	mux.HandleFunc("DELETE /role/{id}", roleH.Delete)

	mux.HandleFunc("POST /user", userH.Create)
	mux.HandleFunc("GET /user", userH.List)
	mux.HandleFunc("GET /user/role/{id}", userH.GetByRole)
	mux.HandleFunc("PUT /user/{id}", userH.Update)
	mux.HandleFunc("DELETE /user/{id}", userH.Delete)

	mux.HandleFunc("POST /cart", cartH.Create)
	mux.HandleFunc("GET /cart/{user_id}", cartH.GetByUser)
	mux.HandleFunc("PUT /cart/{id}", cartH.Update)
	mux.HandleFunc("DELETE /cart/{id}", cartH.Delete)

	mux.HandleFunc("POST /order", orderH.Create)
	mux.HandleFunc("GET /order", orderH.List)
	mux.HandleFunc("GET /order/{id}", orderH.GetByID)
	mux.HandleFunc("PUT /order/{id}", orderH.Update)

	mux.HandleFunc("POST /order_item", orderItemH.Create)
	mux.HandleFunc("GET /order_item/{order_id}", orderItemH.GetByOrder)

	slog.Info("сервер запускается", "port", cfg.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), mux); err != nil {
		slog.Error("ошибка сервера", "err", err)
	}
}
