package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"log"
	"net/http"
	"redditclone/pkg/auth"
	"redditclone/pkg/repository"
	"strings"
)

type commentRequest struct {
	Comment string `json:"comment"`
}

type deleteResponse struct {
	Message string `json:"message"`
}
type PostHandler struct {
	PostRepo *repository.InMemoryPostRepo
	logger   *zap.SugaredLogger
}

func NewPostHandler(logger *zap.SugaredLogger) *PostHandler {
	return &PostHandler{
		PostRepo: repository.NewInMemoryPostRepo(),
		logger:   logger,
	}
}

func (h *PostHandler) ListAllPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostRepo.ListAll()
	if err != nil {
		h.logger.Errorw("error while listing posts", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Infow("!!!Listing posts", "posts", posts)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		h.logger.Errorw("error while encoding posts", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var req repository.PostRequest
	h.logger.Infow("received register request", "r.body", r.Body)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("error while decoding post request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.logger.Infow("received post request", "post", req)

	fmt.Printf("\n\n\t%#v\n", r.Header)
	fmt.Printf("\t%+v\n\n", r.Header)

	inToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	session, err := auth.ParseToken(inToken)
	if err != nil {
		h.logger.Errorw("error while parsing token", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	post, err := h.PostRepo.Create(req, session)
	if err != nil {
		h.logger.Errorw("error while creating post", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Infow("post created", "post", post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		h.logger.Errorw("encoding new post", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *PostHandler) ListCategoryPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postCatID := vars["CATEGORY_NAME"]

	posts, err := h.PostRepo.GetByCategory(postCatID)
	if err != nil {
		h.logger.Errorw("getting posts by Category", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Infow("got posts by category", "posts", posts)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		h.logger.Errorw("encoding posts category", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *PostHandler) ListPostByID(w http.ResponseWriter, r *http.Request) {
	log.Println("listByPostID")
	vars := mux.Vars(r)
	log.Printf("mux vars: %#v", vars)
	postID := vars["POST_ID"]
	log.Printf("postid: %#v", postID)
	post, err := h.PostRepo.ListByID(postID)
	if err != nil {
		h.logger.Errorw("getting post by ID", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Infow("got post by ID", "post", post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		h.logger.Errorw("encoding to json post by ID", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *PostHandler) AddCommentPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("\n\tmux vars: %#v", vars)
	postID := vars["POST_ID"]
	log.Printf("\tpostid: %#v", postID)
	post, err := h.PostRepo.GetByID(postID)
	if err != nil {
		h.logger.Errorw("getting post by ID", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req commentRequest
	h.logger.Infow("received comment request", "r.body", r.Body)
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("decoding comment request", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	h.logger.Infow("received comment request", "comment", req)

	post, err = h.PostRepo.AddCommentToPost(req.Comment, post.ID)
	if err != nil {
		h.logger.Errorw("adding comment to post", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		h.logger.Errorw("encoding new post comment", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *PostHandler) DeleteCommentPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("mux vars: %#v", vars)
	postID, commentID := vars["POST_ID"], vars["COMMENT_ID"]
	log.Printf("postid: %#v, commID: %#v", postID, commentID)

	post, err := h.PostRepo.DeleteComment(commentID, postID)
	if err != nil {
		h.logger.Errorw("deleting comment", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		h.logger.Errorw("encoding post comment delete", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) UpVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("UPVOTE mux vars: %#v", vars)
	postID := vars["POST_ID"]
	post, err := h.PostRepo.GetByID(postID)
	if err != nil {
		h.logger.Errorw("getting post by ID", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	session, err := auth.ParseToken(inToken)
	if err != nil {
		h.logger.Errorw("error while parsing token", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.PostRepo.UpVote(session.ID, post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		h.logger.Errorw("encoding post upvote", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) DownVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("DOWNVOTE mux vars: %#v", vars)
	postID := vars["POST_ID"]
	post, err := h.PostRepo.GetByID(postID)
	if err != nil {
		h.logger.Errorw("getting post by ID", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	session, err := auth.ParseToken(inToken)
	if err != nil {
		h.logger.Errorw("error while parsing token", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.PostRepo.DownVote(session.ID, post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		h.logger.Errorw("encoding post downvote", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) UnVote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("UNVOTE mux vars: %#v", vars)
	postID := vars["POST_ID"]
	post, err := h.PostRepo.GetByID(postID)
	if err != nil {
		h.logger.Errorw("getting post by ID", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	session, err := auth.ParseToken(inToken)
	if err != nil {
		h.logger.Errorw("error while parsing token", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	h.PostRepo.UnVote(session.ID, post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		h.logger.Errorw("encoding post unvote", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) DeletePostByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("DeletePost mux vars: %#v", vars)
	postID := vars["POST_ID"]
	post, err := h.PostRepo.GetByID(postID)
	if err != nil {
		h.logger.Errorw("getting post by ID", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	inToken := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	session, err := auth.ParseToken(inToken)
	if err != nil {
		h.logger.Errorw("error while parsing token", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if post.Author.ID != session.ID {
		h.logger.Errorw("invalid post author", "error", errors.New("invalid post author"))
		http.Error(w, "post.Author.ID != session.ID", http.StatusUnauthorized)
		return
	}
	fmt.Printf("\n\tREADY TO DELETE POST, postid: %s", post.ID)
	h.PostRepo.DeletePost(post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(deleteResponse{Message: "success"})
	if err != nil {
		h.logger.Errorw("encoding post delete", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *PostHandler) GetPostsUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("UNVOTE mux vars: %#v", vars)
	userLogin := vars["USER_LOGIN"]

	posts, err := h.PostRepo.GetAllPostsUser(userLogin)
	if err != nil {
		h.logger.Errorw("getting all posts user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		h.logger.Errorw("encoding posts user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
