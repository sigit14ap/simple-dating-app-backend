package models

import "time"

type Match struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id"`
	TargetUserID uint      `json:"profile_id"`
	Match        bool      `json:"liked"`
	CreatedAt    time.Time `json:"created_at"`
}

type SwipeRequest struct {
	TargetUserID uint `json:"target_user_id" validate:"required"`
	Match        bool `json:"match" validate:"required"`
}
