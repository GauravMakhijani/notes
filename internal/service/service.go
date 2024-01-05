package service

import (
	"context"

	"github.com/GauravMakhijani/notes/internal/database"
	"github.com/GauravMakhijani/notes/internal/domain"
	"github.com/GauravMakhijani/notes/internal/jwt"
	"github.com/GauravMakhijani/notes/models"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	// User related methods
	CreateNewUser(ctx context.Context, signupReq domain.SignupRequest) error
	LoginUser(ctx context.Context, loginReq domain.LoginRequest) (domain.LoginResponse, error)

	// Note related methods
	CreateNote(ctx context.Context, noteReq domain.NoteRequest) (domain.NoteResponse, error)
	GetNoteByID(ctx context.Context, id string) (domain.NoteResponse, error)
	ListNotes(ctx context.Context) ([]domain.NoteResponse, error)
	DeleteNoteByID(ctx context.Context, id string) error
	UpdateNoteByID(ctx context.Context, id string, noteReq domain.NoteRequest) (domain.NoteResponse, error)

	// Share related methods
	ShareNoteWithUser(ctx context.Context, noteID string, shareReq domain.SharedNoteRequest) error
	SearchNotes(ctx context.Context, query string) ([]domain.NoteResponse, error)
}

type service struct {
	store database.Storer
}

func NewService(store database.Storer) Service {
	return &service{store: store}
}

func (s *service) CreateNewUser(ctx context.Context, signupReq domain.SignupRequest) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(signupReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     signupReq.Username,
		PasswordHash: string(hashedPassword),
	}

	return s.store.CreateNewUser(user)
}

func (s *service) LoginUser(ctx context.Context, loginReq domain.LoginRequest) (domain.LoginResponse, error) {
	user, err := s.store.GetUserByUsername(loginReq.Username)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		return domain.LoginResponse{}, err
	}

	// Generate JWT token
	accessToken, err := jwt.GenerateToken(user.ID, user.Username)

	return domain.LoginResponse{
		Username:    user.Username,
		AccessToken: accessToken,
	}, nil
}

// Note related methods
func (s *service) CreateNote(ctx context.Context, noteReq domain.NoteRequest) (domain.NoteResponse, error) {

	userId := ctx.Value("user_id").(string)

	note := &models.Note{
		Title:   noteReq.Title,
		Content: noteReq.Body,
		UserID:  userId,
	}

	note, err := s.store.CreateNewNote(note)
	if err != nil {
		return domain.NoteResponse{}, err
	}
	return domain.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Body:      note.Content,
		CreatedBy: note.UserID,
	}, nil

}

func (s *service) GetNoteByID(ctx context.Context, id string) (domain.NoteResponse, error) {

	userId := ctx.Value("user_id").(string)

	note, err := s.store.GetNoteByID(userId, id)
	if err != nil {
		return domain.NoteResponse{}, err
	}

	return domain.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Body:      note.Content,
		CreatedBy: note.UserID,
	}, nil
}

func (s *service) ListNotes(ctx context.Context) ([]domain.NoteResponse, error) {

	userId := ctx.Value("user_id").(string)

	notes, err := s.store.ListNotes(userId)
	if err != nil {
		return []domain.NoteResponse{}, err
	}

	noteResponses := make([]domain.NoteResponse, 0)
	for _, note := range notes {
		noteResponses = append(noteResponses, domain.NoteResponse{
			ID:        note.ID,
			Title:     note.Title,
			Body:      note.Content,
			CreatedBy: note.UserID,
		})
	}

	return noteResponses, nil
}

func (s *service) DeleteNoteByID(ctx context.Context, id string) error {

	userId := ctx.Value("user_id").(string)
	return s.store.DeleteNoteByID(userId, id)
}

func (s *service) UpdateNoteByID(ctx context.Context, id string, noteReq domain.NoteRequest) (domain.NoteResponse, error) {
	userID := ctx.Value("user_id").(string)

	note := &models.Note{
		Title:   noteReq.Title,
		Content: noteReq.Body,
	}

	note, err := s.store.UpdateNoteByID(userID, id, note)
	if err != nil {
		return domain.NoteResponse{}, err
	}

	return domain.NoteResponse{
		ID:        note.ID,
		Title:     note.Title,
		Body:      note.Content,
		CreatedBy: note.UserID,
	}, nil
}

// ShareNoteWithUser shares the note with the given user
func (s *service) ShareNoteWithUser(ctx context.Context, noteID string, shareReq domain.SharedNoteRequest) error {
	fromID := ctx.Value("user_id").(string)
	return s.store.ShareNoteWithUser(noteID, fromID, shareReq.ToUsersID)
}

func (s *service) SearchNotes(ctx context.Context, query string) ([]domain.NoteResponse, error) {
	userId := ctx.Value("user_id").(string)

	notes, err := s.store.SearchNotes(userId, query)
	if err != nil {
		return []domain.NoteResponse{}, err
	}

	noteResponses := make([]domain.NoteResponse, 0)
	for _, note := range notes {
		noteResponses = append(noteResponses, domain.NoteResponse{
			ID:        note.ID,
			Title:     note.Title,
			Body:      note.Content,
			CreatedBy: note.UserID,
		})
	}

	return noteResponses, nil
}
