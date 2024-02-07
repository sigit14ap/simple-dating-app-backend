package services

import (
	"fmt"
	"time"

	"github.com/sigit14ap/simple-dating-app-backend/internal/models"
	"github.com/sigit14ap/simple-dating-app-backend/internal/repositories"
)

type MatchServiceInterface interface {
	RecordMatch(userID, TargetUserID uint, match bool) error
	GetNextUser(userID uint) (*models.User, error)
	FindMatch(userID uint, targetUserID uint, date time.Time) (models.Match, error)
	UpdateMatch(userID uint, targetUserID uint, match bool) (bool, error)
}

type matchService struct {
	matchRepository repositories.MatchRepository
	userRepository  repositories.UserRepository
}

func NewMatchService(matchRepository repositories.MatchRepository, userRepository repositories.UserRepository) MatchServiceInterface {
	return &matchService{
		userRepository:  userRepository,
		matchRepository: matchRepository,
	}
}

func (s *matchService) RecordMatch(userID, TargetUserID uint, match bool) error {
	return s.matchRepository.RecordMatch(userID, TargetUserID, match)
}

func (s *matchService) GetNextUser(userID uint) (*models.User, error) {
	currentDate := time.Now()
	targetUsers, err := s.matchRepository.GetMatch(userID, currentDate)

	if err != nil {
		return nil, err
	}
	fmt.Print(targetUsers)
	var targetUserIds []uint
	for _, user := range targetUsers {
		targetUserIds = append(targetUserIds, user.TargetUserID)
	}

	targetUserIds = append(targetUserIds, userID)

	user, err := s.userRepository.FindNextUser(targetUserIds)

	if err != nil {
		return nil, err
	}

	err = s.RecordMatch(userID, user.ID, false)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *matchService) FindMatch(userID uint, targetUserID uint, date time.Time) (models.Match, error) {
	return s.matchRepository.FindMatch(userID, targetUserID, date)
}

func (s *matchService) UpdateMatch(userID uint, targetUserID uint, match bool) (bool, error) {
	return s.matchRepository.UpdateMatch(userID, targetUserID, match)
}
