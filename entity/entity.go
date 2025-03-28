package entity

import (
	"time"
)

type User struct {
	ID         uint   `gorm:"primaryKey"`
	Email      string `gorm:"not null;uniqueIndex;size:255"`
	Password   string `gorm:"not null"`
	Username   string `gorm:"not null;uniqueIndex;size:50"`
	IsVerified bool   `gorm:"default:false"`
	CreatedAt  time.Time
}

type Quiz struct {
	ID        uint       `gorm:"primaryKey"`
	Title     string     `gorm:"not null"`
	CreatorID *uint      `gorm:"null:index"`
	Questions []Question `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time  `gorm:"not null;autoCreateTime"`
	User      User       `gorm:"foreignKey:CreatorID;constraint:OnDelete:SET NULL;"`
}

type Question struct {
	ID      uint     `gorm:"primaryKey"`
	QuizID  uint     `gorm:"not null;index"`
	Quiz    Quiz     `gorm:"foreignKey:QuizID"`
	Text    string   `gorm:"not null"`
	Answers []Answer `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE;"`
}

type Answer struct {
	ID         uint     `gorm:"primaryKey"`
	QuestionID uint     `gorm:"not null;index"`
	Question   Question `gorm:"foreignKey:QuestionID"`
	Text       string   `gorm:"not null"`
	IsCorrect  bool     `gorm:"not null"`
}

type Submission struct {
	ID        uint  `gorm:"primaryKey"`
	QuizID    uint  `gorm:"index"`
	UserID    *uint `gorm:"null;index"`
	Score     float32
	CreatedAt time.Time              `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time              `gorm:"not null;autoUpdateTime"`
	User      User                   `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL;"`
	Quiz      Quiz                   `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE;"`
	Answers   []SubmissionUserAnswer `gorm:"foreignKey:SubmissionID;constraint:OnDelete:CASCADE;"`
}

type SubmissionUserAnswer struct {
	ID           uint `gorm:"primaryKey"`
	SubmissionID uint `gorm:"not null;index"`
	QuestionID   uint `gorm:"not null"`
	UserAnswerID uint `gorm:"not null"`
	CorrectID    uint `gorm:"not null"`
	IsCorrect    bool `gorm:"not null"`
}
