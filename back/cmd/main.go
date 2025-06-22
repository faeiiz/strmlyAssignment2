package main

import (
	"back/handlers"
	"back/initializers"
	"back/middlewares"
	"back/repository"
	"back/services"
	"log"
	"net/http"
	"os"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectDB()
	initializers.ConnectCloudinary()
}

func main() {

	// This is dependcy injection to mimic clean and abstract cdoe

	userRepo := repository.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	videoRepo := repository.NewVideoRepository()
	videoService := services.NewVideoService(videoRepo)
	videoHandler := handlers.NewVideoHandler(videoService)

	fs := http.FileServer(http.Dir("front"))
	http.Handle("/", fs)

	// pub routes -> W/o JWT auth

	http.Handle("/signup",
		middlewares.RateLimit(
			middlewares.SecurityHeaders(
				http.HandlerFunc(userHandler.Signup),
			),
		),
	)

	http.Handle("/login",
		middlewares.RateLimit(
			middlewares.SecurityHeaders(
				http.HandlerFunc(userHandler.Login),
			),
		),
	)

	// protected routes -> With JWT auth

	http.Handle("/profile",
		middlewares.RateLimit(
			middlewares.SecurityHeaders(
				middlewares.JWTAuth(
					http.HandlerFunc(userHandler.Profile),
				),
			),
		),
	)

	http.Handle("/upload",
		middlewares.RateLimit(
			middlewares.SecurityHeaders(
				middlewares.JWTAuth(
					http.HandlerFunc(videoHandler.UploadVideo),
				),
			),
		),
	)

	http.Handle("/videos",
		middlewares.RateLimit(
			middlewares.SecurityHeaders(
				middlewares.JWTAuth(
					http.HandlerFunc(videoHandler.GetVideos),
				),
			),
		),
	)
	corsWrappedMux := middlewares.CORSMiddleware(http.DefaultServeMux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server starting on port", port)
	err := http.ListenAndServe(":"+port, corsWrappedMux)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
