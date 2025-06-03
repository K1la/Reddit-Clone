package repository

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"redditclone/pkg/models"
	"sync"
	"time"
)

type (
	InMemoryPostRepo struct {
		posts map[string]*models.Post
		mu    sync.RWMutex
	}
	PostRequest struct {
		Category string `json:"category"`
		Text     string `json:"text,omitempty"`
		URL      string `json:"url,omitempty"`
		Title    string `json:"title"`
		Type     string `json:"type"`
	}
)

func NewInMemoryPostRepo() *InMemoryPostRepo {
	return &InMemoryPostRepo{
		posts: make(map[string]*models.Post),
	}
}

func (h *InMemoryPostRepo) ListAll() ([]*models.Post, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	res := make([]*models.Post, 0, len(h.posts))
	log.Printf("listAll: %v", len(h.posts))
	log.Printf("listAll: %v", h.posts)
	for _, p := range h.posts {
		res = append(res, p)
	}
	return res, nil
}

func (h *InMemoryPostRepo) Create(postReq PostRequest, session *models.Session) (*models.Post, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	post := &models.Post{
		Score:      1,
		Views:      0,
		Type:       postReq.Type,
		Title:      postReq.Title,
		Category:   postReq.Category,
		Created:    time.Now(),
		UpVotePerc: 100,
		ID:         uuid.NewString(),
		Votes:      make([]*models.Vote, 0),
		Comments:   make([]*models.Comment, 0),
		Author:     session,
	}
	switch postReq.Type {
	case "text":
		post.Text = postReq.Text
	case "link":
		post.URL = postReq.URL
	}
	h.posts[post.ID] = post
	h.UpVote(session.ID, post)
	return post, nil
}

func (h *InMemoryPostRepo) ListByID(id string) (*models.Post, error) {
	post, err := h.GetByID(id)
	if err != nil {
		return nil, err
	}
	post.Views++
	return post, nil
}

func (h *InMemoryPostRepo) GetByID(id string) (*models.Post, error) {
	post, ok := h.posts[id]
	if !ok {
		return nil, errors.New("post not found")
	}
	return post, nil
}

func (h *InMemoryPostRepo) GetByCategory(category string) ([]*models.Post, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var res []*models.Post
	for _, p := range h.posts {
		if p.Category == category {
			res = append(res, p)
		}
	}
	return res, nil
}

func (h *InMemoryPostRepo) AddCommentToPost(body string, postID string) (*models.Post, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	post, ok := h.posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}

	comm := h.addComment(body, post.ID)
	post.Comments = append(post.Comments, comm)
	return post, nil
}

func (h *InMemoryPostRepo) DeleteComment(commentID, postID string) (*models.Post, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	post, ok := h.posts[postID]
	if !ok {
		return nil, errors.New("post not found")
	}
	err := h.deleteComment(post.ID, commentID)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (h *InMemoryPostRepo) checkVote(postID string) bool {
	return len(h.posts[postID].Votes) > 0
}

func (h *InMemoryPostRepo) deleteVote(sessionID, postID string) {

	position := -1
	for i, vote := range h.posts[postID].Votes {
		if vote.User == sessionID {
			position = i
			break
		}
	}
	if position == -1 {
		return
	}
	h.posts[postID].Votes = append(h.posts[postID].Votes[:position], h.posts[postID].Votes[position+1:]...)
}

func upVote(userID string) *models.Vote {
	vote := &models.Vote{
		User: userID,
		Vote: 1,
	}
	return vote
}

func downVote(userID string) *models.Vote {
	vote := &models.Vote{
		User: userID,
		Vote: -1,
	}
	return vote
}

func (h *InMemoryPostRepo) calcUpVotePercent(post *models.Post) {
	if len(post.Votes) == 0 {
		h.posts[post.ID].UpVotePerc = 0
		return
	}
	up := 0
	for _, v := range post.Votes {
		if v.Vote == 1 {
			up++
		}
	}
	if up == 0 {
		h.posts[post.ID].UpVotePerc = 0
		return
	}
	percentage := up * 100 / len(post.Votes)
	h.posts[post.ID].UpVotePerc = percentage
}

func (h *InMemoryPostRepo) UpVote(sessionID string, post *models.Post) {
	if h.checkVote(post.ID) {
		h.deleteVote(sessionID, post.ID)
	}

	vote := upVote(sessionID)
	post.Votes = append(post.Votes, vote)
	h.calcUpVotePercent(post)
}

func (h *InMemoryPostRepo) DownVote(sessionID string, post *models.Post) {
	if h.checkVote(post.ID) {
		h.deleteVote(sessionID, post.ID)
	}
	vote := downVote(sessionID)
	post.Votes = append(post.Votes, vote)
	h.calcUpVotePercent(post)
}

func (h *InMemoryPostRepo) UnVote(sessionID string, post *models.Post) {
	if h.checkVote(post.ID) {
		h.deleteVote(sessionID, post.ID)
	}
	h.calcUpVotePercent(post)
}

func (h *InMemoryPostRepo) DeletePost(post *models.Post) {
	delete(h.posts, h.posts[post.ID].ID)
}

func (h *InMemoryPostRepo) GetAllPostsUser(userLogin string) ([]*models.Post, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	var res []*models.Post
	for _, p := range h.posts {
		if p.Author.Username == userLogin {
			res = append(res, p)
		}
	}
	if len(res) == 0 {
		return nil, errors.New("posts not found")
	}
	return res, nil
}

func (h *InMemoryPostRepo) addComment(body, postID string) *models.Comment {
	comm := &models.Comment{
		Created: time.Now(),
		Author:  h.posts[postID].Author,
		Body:    body,
		ID:      uuid.NewString(),
	}
	return comm
}

func (h *InMemoryPostRepo) getCommentPosition(postID, commentID string) (int, error) {
	comments := h.posts[postID].Comments
	for i, comment := range comments {
		if comment.ID == commentID {
			return i, nil
		}
	}
	return -1, fmt.Errorf("comment not found")
}

func (h *InMemoryPostRepo) deleteComment(postID, commentID string) error {
	position, err := h.getCommentPosition(postID, commentID)
	if err != nil {
		return err
	}
	h.posts[postID].Comments = append(h.posts[postID].Comments[:position], h.posts[postID].Comments[position+1:]...)
	return nil
}
