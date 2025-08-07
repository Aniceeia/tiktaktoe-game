package service

import (
	"database/sql"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req SignUpRequest) error
	Authenticate(login, password string) (string, error)
	GetUserByUUID(uuid string) (*User, error)
}

type User struct {
	UUID         string
	Username     string
	PasswordHash string
	Score        int
}

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type PostgresAuthService struct {
	db *sql.DB
}

func NewPostgresAuthService(db *sql.DB) *PostgresAuthService {
	return &PostgresAuthService{db: db}
}

func (s *PostgresAuthService) Register(req SignUpRequest) error {
	log.Printf("Attempting to register user: %s", req.Login)

	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, req.Login).Scan(&exists)
	if err != nil {
		log.Printf("Error checking user existence: %v", err)
		return err
	}

	if exists {
		log.Printf("User %s already exists", req.Login)
		return ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}

	_, err = s.db.Exec(`
        INSERT INTO users (uuid, username, password_hash, score)
        VALUES (gen_random_uuid(), $1, $2, 0)
    `, req.Login, string(hash))

	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}

	log.Printf("User %s registered successfully", req.Login)
	return nil
}

func (s *PostgresAuthService) Authenticate(login, password string) (string, error) {
	var uuid, hash string
	err := s.db.QueryRow(`
		SELECT uuid, password_hash 
		FROM users 
		WHERE username = $1
	`, login).Scan(&uuid, &hash)

	if errors.Is(err, sql.ErrNoRows) {
		return "", errors.New("invalid credentials")
	}
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	return uuid, nil
}

func (s *PostgresAuthService) GetUserByUUID(uuid string) (*User, error) {
	user := &User{}
	err := s.db.QueryRow(`
		SELECT uuid, username, score 
		FROM users 
		WHERE uuid = $1
	`, uuid).Scan(&user.UUID, &user.Username, &user.Score)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	return user, err
}
