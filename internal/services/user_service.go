package services

import (
	"errors"
	"time"

	"github.com/sigit14ap/simple-dating-app-backend/internal/models"
	"github.com/sigit14ap/simple-dating-app-backend/internal/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	RegisterUser(username, password string) error
	GetUserByUsername(username string) (*models.User, error)
	AuthenticateUser(username, password string) (*models.User, error)
	GetUserById(userID uint) (*models.User, error)
	CheckEligibleSwipe(userID uint) (bool, error)
	UpdateUnlimitedSwipe(userID uint) error
}

type UserService struct {
	userRepository  repositories.UserRepository
	matchRepository repositories.MatchRepository
}

func NewUserService(userRepository repositories.UserRepository, matchRepository repositories.MatchRepository) UserServiceInterface {
	return &UserService{
		userRepository:  userRepository,
		matchRepository: matchRepository,
	}
}

func (s *UserService) RegisterUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}
	return s.userRepository.CreateUser(user)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.userRepository.GetUserByUsername(username)
}

func (s *UserService) GetUserById(userID uint) (*models.User, error) {
	return s.userRepository.GetUserById(userID)
}

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) CheckEligibleSwipe(userID uint) (bool, error) {
	currentDate := time.Now()
	currentUser, err := s.userRepository.GetUserById(userID)

	if err != nil {
		return false, errors.New("account not found")
	}

	totalSwipe, err := s.matchRepository.CountMatch(userID, currentDate)

	if err != nil {
		return false, err
	}

	unlimitedSwipe := currentUser.UnlimitedSwipeExpiredAt
	if !unlimitedSwipe.Valid && totalSwipe == 10 {
		return false, errors.New("reached swipe limit")
	}

	if unlimitedSwipe.Valid && unlimitedSwipe.Time.Before(currentDate) {
		return false, errors.New("unlimited swipe feature is expired")
	}

	return true, nil
}

func (s *UserService) UpdateUnlimitedSwipe(userID uint) error {
	return s.userRepository.UpdateUnlimitedSwipe(userID)
}
