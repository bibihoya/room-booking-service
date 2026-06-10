package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/internships-backend/test-backend-bibihoya-1/api/docs"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/config"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/database"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/handler"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/middleware"
	"github.com/internships-backend/test-backend-bibihoya-1/internal/storage"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           Meeting Room Booking Service API
// @version         1.0
// @description     Сервис бронирования переговорок с автоматической генерацией слотов
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите токен в формате "Bearer {token}"
func main() {
	cfg := config.Load()

	db, err := database.NewPostgresDB(
		cfg.DBHost, cfg.DBPort, cfg.DBUser,
		cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer db.Close()

	if err = database.Migrate(db); err != nil {
		log.Fatalf("cannot migrate database: %v", err)
	}

	userStorage := storage.NewUserStorage(db)
	roomStorage := storage.NewRoomStorage(db)
	scheduleStorage := storage.NewScheduleStorage(db)
	slotStorage := storage.NewSlotStorage(db)
	bookingStorage := storage.NewBookingStorage(db)

	authHandler := handler.NewAuthHandler(cfg.JWTSecret, cfg.JWTTTL, userStorage)
	roomHandler := handler.NewRoomHandler(roomStorage)
	scheduleHandler := handler.NewScheduleHandler(scheduleStorage)
	slotHandler := handler.NewSlotHandler(slotStorage, roomStorage)
	bookingHandler := handler.NewBookingHandler(bookingStorage)

	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	)).Methods("GET")

	router.HandleFunc("/_info", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	router.HandleFunc("/dummyLogin", authHandler.DummyLogin).Methods("POST")
	router.HandleFunc("/register", authHandler.Register).Methods("POST")
	router.HandleFunc("/login", authHandler.Login).Methods("POST")

	router.Handle("/rooms/list", middleware.AuthMiddleware(cfg.JWTSecret)(http.HandlerFunc(roomHandler.RoomList))).Methods("GET")
	router.Handle("/rooms/create", middleware.AuthMiddleware(cfg.JWTSecret)(middleware.AdminCheck(http.HandlerFunc(roomHandler.RoomCreate)))).Methods("POST")

	router.Handle("/rooms/{roomId}/schedule/create", middleware.AuthMiddleware(cfg.JWTSecret)(middleware.AdminCheck(http.HandlerFunc(scheduleHandler.CreateSchedule)))).Methods("POST")

	router.Handle("/rooms/{roomId}/slots/list", middleware.AuthMiddleware(cfg.JWTSecret)(http.HandlerFunc(slotHandler.ListAvailableSlots))).Methods("GET")

	router.Handle("/bookings/create", middleware.AuthMiddleware(cfg.JWTSecret)(middleware.UserCheck(http.HandlerFunc(bookingHandler.CreateBooking)))).Methods("POST")
	router.Handle("/bookings/my", middleware.AuthMiddleware(cfg.JWTSecret)(middleware.UserCheck(http.HandlerFunc(bookingHandler.GetMyBooking)))).Methods("GET")
	router.Handle("/bookings/{bookingId}/cancel", middleware.AuthMiddleware(cfg.JWTSecret)(middleware.UserCheck(http.HandlerFunc(bookingHandler.CancelBooking)))).Methods("POST")
	router.Handle("/bookings/list", middleware.AuthMiddleware(cfg.JWTSecret)(middleware.AdminCheck(http.HandlerFunc(bookingHandler.GetAllBookings)))).Methods("GET")

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %s\n", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server gracefully stopped")
}
