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
	Questions []Question `gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE;"`
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
	ID        uint `gorm:"primaryKey"`
	QuizID    uint `gorm:"index"`
	UserID    uint `gorm:"not null;index"`
	Score     float32
	CreatedAt time.Time `gorm:"not null;autoCreateTime"`
	UpdatedAt time.Time `gorm:"not null;autoUpdateTime"`

	User User `gorm:"foreignKey:UserID"`
	Quiz Quiz `gorm:"foreignKey:QuizID"`
}
