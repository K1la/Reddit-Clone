package handlers

import (
	"encoding/json"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"redditclone/pkg/auth"
	"redditclone/pkg/repository"
)

type UserHandler struct {
	UserRepo *repository.InMemoryUserRepo
	Sessions *repository.InMemorySessionRepo
	logger   *zap.SugaredLogger
}

func NewUserHandler(logger *zap.SugaredLogger) *UserHandler {
	return &UserHandler{
		UserRepo: repository.NewInMemoryUserRepo(),
		Sessions: repository.NewInMemorySessionRepo(),
		logger:   logger,
	}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

func (h *UserHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/html/index.html")
	h.logger.Infoln("INDEX servedFile, redirected to /api/posts/")
}

func (h *UserHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	h.logger.Infow("received register request", "r.body", r.Body)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("error while decoding request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Infow("received register request", "req", req)

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Errorw("error while hashing password", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.UserRepo.Create(req.Username, string(hash))
	if err != nil {
		h.logger.Errorw("error while creating user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := h.Sessions.Create(user.Username)
	if err != nil {
		h.logger.Errorw("error while creating session", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Infow("session register", "session", session)

	token, err := auth.GenerateToken(session.ID, session.Username)
	if err != nil {
		h.logger.Errorw("error while generating token", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Infow("token register", "token", token)

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(authResponse{Token: token})
	if err != nil {
		h.logger.Errorw("error while encoding response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorw("error while decoding request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.logger.Infow("received login request", "req", req)

	user, err := h.UserRepo.GetByUsername(req.Username)
	if err != nil {
		h.logger.Errorw("error while getting user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		h.logger.Errorw("error while comparing password", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := h.Sessions.Create(user.Username)
	if err != nil {
		h.logger.Errorw("error while creating session", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	h.logger.Infow("session login", "session", session)

	token, err := auth.GenerateToken(session.ID, session.Username)
	if err != nil {
		h.logger.Errorw("error while generating token", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.logger.Infow("token login", "token", token)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(authResponse{Token: token})
	if err != nil {
		h.logger.Errorw("error while encoding response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
