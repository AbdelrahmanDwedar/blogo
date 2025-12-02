package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	deliveryHttp "AbdelrahmanDwedar/blogo/internal/delivery/http"
	"AbdelrahmanDwedar/blogo/internal/infrastructure/cache"
	"AbdelrahmanDwedar/blogo/internal/infrastructure/database"
	"AbdelrahmanDwedar/blogo/internal/usecase"
	"AbdelrahmanDwedar/blogo/pkg/auth"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	db, err := database.NewPostgresDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize tables
	if err := db.InitTables(); err != nil {
		log.Fatal("Failed to initialize tables:", err)
	}

	// Initialize cache (optional)
	redisCache := cache.NewRedisCache()
	if redisCache != nil {
		defer redisCache.Close()
	}

	// Initialize repositories
	userRepo := database.NewUserRepository(db)
	blogRepo := database.NewBlogRepository(db)

	// Initialize use cases
	userUC := usecase.NewUserUseCase(userRepo, redisCache)
	blogUC := usecase.NewBlogUseCase(blogRepo, redisCache)

	// Initialize HTTP handlers
	handler := deliveryHttp.NewHandler(userUC, blogUC)

	// Setup router
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/ping", handler.Ping).Methods("GET")

	// User routes
	r.HandleFunc("/api/u/new", handler.UserHandler.CreateUser).Methods("POST")
	r.HandleFunc("/api/u/{id:[0-9]+}", handler.UserHandler.GetUser).Methods("GET")
	r.HandleFunc("/api/u/{id:[0-9]+}", auth.AuthMiddleware(handler.UserHandler.FollowUser)).Methods("POST")
	r.HandleFunc("/api/u/{id:[0-9]+}/manage", auth.AuthMiddleware(handler.UserHandler.UpdateUser)).Methods("POST")
	r.HandleFunc("/api/u/{id:[0-9]+}/following", handler.UserHandler.GetUserFollowing).Methods("GET")
	r.HandleFunc("/api/u/{id:[0-9]+}/follows", handler.UserHandler.GetUserFollowers).Methods("GET")

	// Blog routes
	r.HandleFunc("/api/b", handler.BlogHandler.GetBlogs).Methods("GET")
	r.HandleFunc("/api/b/new", auth.AuthMiddleware(handler.BlogHandler.CreateBlog)).Methods("POST")
	r.HandleFunc("/api/b/{id:[0-9]+}", handler.BlogHandler.GetBlog).Methods("GET")
	r.HandleFunc("/api/b/{id:[0-9]+}", auth.AuthMiddleware(handler.BlogHandler.LikeBlog)).Methods("POST")
	r.HandleFunc("/api/b/{id:[0-9]+}/edit", auth.AuthMiddleware(handler.BlogHandler.UpdateBlog)).Methods("POST")
	r.HandleFunc("/api/b/{id:[0-9]+}/delete", auth.AuthMiddleware(handler.BlogHandler.DeleteBlog)).Methods("POST")
	r.HandleFunc("/api/b/{id:[0-9]+}/likes", handler.BlogHandler.GetBlogLikes).Methods("GET")

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 Server starting on port %s\n", port)
	log.Printf("📖 API documentation: http://localhost:%s/ping\n", port)
	log.Fatal(srv.ListenAndServe())
}


