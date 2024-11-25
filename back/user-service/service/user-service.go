package service

import (
	"context"

	"github.com/milly013/trello-project/back/user-service/model"
	"github.com/milly013/trello-project/back/user-service/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Provera da li korisnik postoji prema korisničkom imenu ili emailu
func (s *UserService) CheckUserExists(ctx context.Context, username, email string) (bool, error) {
	return s.repo.CheckUserExists(ctx, username, email)
}

// Čuvanje verifikacionog koda za korisnika
func (s *UserService) SaveVerificationCode(ctx context.Context, user model.User, code string) error {
	return s.repo.SaveVerificationCode(ctx, user, code)
}

// Provera verifikacionog koda
func (s *UserService) VerifyCode(ctx context.Context, email, code string) (bool, error) {
	return s.repo.VerifyCode(ctx, email, code)
}

// Preuzimanje korisnika na osnovu emaila
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := s.repo.GetUserByEmail(ctx, email, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Kreiranje novog korisnika
func (s *UserService) CreateUser(ctx context.Context, user model.User) error {
	return s.repo.CreateUser(ctx, user)
}

// Brisanje korisnika po ID-u
func (s *UserService) DeleteUserByID(ctx context.Context, id string) error {
	return s.repo.DeleteUserByID(ctx, id)
}

// Preuzimanje korisnika po ID-u
func (s *UserService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := s.repo.GetUserByID(ctx, id, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Preuzimanje svih korisnika
func (s *UserService) GetAllUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.GetAllUsers(ctx)
}

// Verifikacija i aktivacija korisnika na osnovu emaila i koda
func (s *UserService) VerifyUserAndActivate(ctx context.Context, email, code string) (bool, error) {
	return s.repo.VerifyUserAndActivate(ctx, email, code)
}

// Promena lozinke korisnika
func (s *UserService) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	// Preuzmi korisnika prema ID-u
	var user model.User
	err := s.repo.GetUserByID(ctx, userID, &user)
	if err != nil {
		return err
	}

	// Proveri da li trenutna lozinka odgovara
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		return err // Trenutna lozinka nije validna
	}

	// Heširaj novu lozinku
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Ažuriraj korisnika sa novom lozinkom
	return s.repo.UpdatePassword(ctx, userID, string(hashedPassword))
}
