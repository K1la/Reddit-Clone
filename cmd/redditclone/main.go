package main

import (
	"log"
	"net/http"
	"redditclone/pkg/auth"
	"redditclone/pkg/handlers"
	"redditclone/pkg/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to initialize zap logger", zap.Error(err))
		return
	}
	logger := zapLogger.Sugar()
	defer func() {
		err = zapLogger.Sync()
		if err != nil {
			logger.Fatal(err)
			return
		}
	}()

	_ = godotenv.Load(".env")
	auth.Init()

	r := mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/html/index.html")
	}).Methods("GET")

	authHandler := handlers.NewUserHandler(logger)
	postsHandler := handlers.NewPostHandler(logger)

	r.HandleFunc("/api/register", authHandler.RegisterPage).Methods("POST")
	r.HandleFunc("/api/login", authHandler.LoginPage).Methods("POST")

	r.HandleFunc("/api/posts/", postsHandler.ListAllPosts).Methods("GET")
	r.HandleFunc("/api/posts", postsHandler.CreatePost).Methods("POST")
	r.HandleFunc("/api/posts/{CATEGORY_NAME}", postsHandler.ListCategoryPosts).Methods("GET")

	r.HandleFunc("/api/post/{POST_ID}", postsHandler.ListPostByID).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}", postsHandler.AddCommentPost).Methods("POST")
	r.HandleFunc("/api/post/{POST_ID}/{COMMENT_ID}", postsHandler.DeleteCommentPost).Methods("DELETE")

	r.HandleFunc("/api/post/{POST_ID}/upvote", postsHandler.UpVote).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/downvote", postsHandler.DownVote).Methods("GET")
	r.HandleFunc("/api/post/{POST_ID}/unvote", postsHandler.UnVote).Methods("GET")

	r.HandleFunc("/api/post/{POST_ID}", postsHandler.DeletePostByID).Methods("DELETE")

	r.HandleFunc("/api/user/{USER_LOGIN}", postsHandler.GetPostsUser).Methods("GET")

	// MiddleWares
	muxMW := middleware.AccessLog(logger, r)
	muxMW = middleware.Panic(logger, muxMW)

	addr := ":8032"
	logger.Infof("Starting server on %s", addr)
	logger.Fatal(http.ListenAndServe(addr, muxMW))

}
