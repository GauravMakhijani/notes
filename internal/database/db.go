package database

import (
	"log"

	"github.com/GauravMakhijani/notes/models"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Storer represents the database operations interface
type Storer interface {
	AutoMigrate() error
	CreateNewUser(user *models.User) error
	GetUserByUsername(username string) (*models.User, error)

	// Note related methods
	CreateNewNote(note *models.Note) (*models.Note, error)
	GetNoteByID(userId, id string) (*models.Note, error)
	ListNotes(userID string) ([]*models.Note, error)
	DeleteNoteByID(userId, id string) error
	UpdateNoteByID(userId, id string, note *models.Note) (*models.Note, error)
	ShareNoteWithUser(noteID string, fromUserID string, toUsersID []string) error
	SearchNotes(userID, query string) ([]*models.Note, error)
}

// store is the concrete implementation of the Storer interface
type store struct {
	db *gorm.DB
}

// NewStore creates a new instance of the database store
func NewStore() Storer {
	dsn := "host=localhost port=5432 user=postgres dbname=notes sslmode=disable password=postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}

	return &store{db: db}
}

// AutoMigrate performs automatic migration of database tables
func (s *store) AutoMigrate() error {
	return s.db.AutoMigrate(&models.User{}, &models.Note{}, &models.SharedNote{})
}

// CreateNewUser creates a new user in the database
func (s *store) CreateNewUser(user *models.User) error {
	return s.db.Create(user).Error
}

// GetUserByUsername fetches the user from the database by username
func (s *store) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateNewNote creates a new note in the database
func (s *store) CreateNewNote(note *models.Note) (*models.Note, error) {
	err := s.db.Create(note).Error
	if err != nil {
		return nil, err
	}
	return note, nil
}

// GetNoteByID fetches the note from the database by ID
func (s *store) GetNoteByID(userId, id string) (*models.Note, error) {
	var note models.Note
	err := s.db.Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, false).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

// ListNotes fetches all the notes from the database for the given user
func (s *store) ListNotes(userID string) ([]*models.Note, error) {
	var notes []*models.Note
	err := s.db.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&notes).Error
	if err != nil {
		return nil, err
	}
	// add shared notes also to the list
	var sharedNotes []*models.SharedNote
	err = s.db.Where("to_user_id = ?", userID).Find(&sharedNotes).Error
	if err != nil {
		return nil, err
	}
	for _, sharedNote := range sharedNotes {
		note, err := s.GetNoteByID(sharedNote.FromUserID, sharedNote.NoteID)
		if err != nil {
			logrus.Errorf("error getting note by id\nError: %s", err.Error())
			continue
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *store) DeleteNoteByID(userId, id string) error {

	err := s.db.Model(&models.Note{}).Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, false).Update("is_deleted", true).Error
	if err != nil {
		logrus.Errorf("error deleting note\nError: %s", err.Error())
		return err
	}
	return nil
}

func (s *store) UpdateNoteByID(userId, id string, note *models.Note) (*models.Note, error) {
	err := s.db.Where("id = ? AND user_id = ? AND is_deleted = ?", id, userId, false).Updates(note).Error
	if err != nil {
		return nil, err
	}
	note, err = s.GetNoteByID(userId, id)
	if err != nil {
		return nil, err
	}

	return note, nil
}

// ShareNoteWithUser shares the note with the given user
func (s *store) ShareNoteWithUser(noteID string, fromUserID string, toUsersName []string) error {
	var sharedNote []*models.SharedNote
	for _, toUserName := range toUsersName {
		toUser, err := s.GetUserByUsername(toUserName)
		if err != nil {
			logrus.Errorf("error getting user by username\nError: %s", err.Error())
			continue
		}
		sharedNote = append(sharedNote, &models.SharedNote{
			NoteID:     noteID,
			FromUserID: fromUserID,
			ToUserID:   toUser.ID,
		})
	}
	err := s.db.Create(sharedNote).Error
	if err != nil {
		return err
	}

	return nil
}

func (s *store) SearchNotes(userID, query string) ([]*models.Note, error) {
	var notes []*models.Note
	// get all the notes for the user from the database matching the query and not deleted query should be in title or content
	err := s.db.Where("user_id = ? AND is_deleted = ? AND (title LIKE ? OR content LIKE ?)", userID, false, "%"+query+"%", "%"+query+"%").Find(&notes).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return notes, nil
}
