package repositories

import (
	"time"

	"github.com/sigit14ap/simple-dating-app-backend/internal/models"
	"gorm.io/gorm"
)

type MatchRepository interface {
	RecordMatch(userID, TargetUserID uint, match bool) error
	GetMatch(userID uint, date time.Time) ([]models.Match, error)
	FindMatch(userID uint, targetUserID uint, date time.Time) (models.Match, error)
	CountMatch(userID uint, date time.Time) (int64, error)
	UpdateMatch(userID uint, targetUserID uint, match bool) (bool, error)
}

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db}
}

func (r *matchRepository) RecordMatch(userID, TargetUserID uint, match bool) error {
	matchData := models.Match{
		UserID:       userID,
		TargetUserID: TargetUserID,
		Match:        match,
		CreatedAt:    time.Now(),
	}
	return r.db.Create(&matchData).Error
}

func (r *matchRepository) GetMatch(userID uint, date time.Time) ([]models.Match, error) {
	var targetUsers []models.Match

	err := r.db.Where("DATE(created_at) >= ?", date.Format("2006-01-02")).Where(&models.Match{UserID: userID}).Find(&targetUsers).Error

	if err != nil {
		return nil, err
	}
	return targetUsers, nil
}

func (r *matchRepository) FindMatch(userID uint, targetUserID uint, date time.Time) (models.Match, error) {
	var match models.Match

	err := r.db.Model(&models.Match{}).Where("DATE(created_at) >= ?", date.Format("2006-01-02")).Where(&models.Match{UserID: userID, TargetUserID: targetUserID}).First(&match).Error

	if err != nil {
		return match, err
	}

	return match, nil
}

func (r *matchRepository) CountMatch(userID uint, date time.Time) (int64, error) {
	var count int64

	err := r.db.Model(&models.Match{}).Where("DATE(created_at) >= ?", date.Format("2006-01-02")).Where(&models.Match{UserID: userID}).Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *matchRepository) UpdateMatch(userID uint, targetUserID uint, match bool) (bool, error) {
	err := r.db.Model(&models.Match{}).Where(&models.Match{UserID: userID, TargetUserID: targetUserID}).Update("match", match).Error

	if err != nil {
		return false, err
	}

	return true, nil
}
