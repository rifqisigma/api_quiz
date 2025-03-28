package dto

import "time"

type Submission struct {
	QuizID  uint               `json:"quiz_id"`
	UserID  uint               `json:"-"`
	Answers []SubmissionAnswer `json:"answers"`
}

type SubmissionAnswer struct {
	QuestionID uint `json:"question_id"`
	AnswerID   uint `json:"answer_id"`
}

type SubmissionUpdate struct {
	SubmissionID uint      `json:"-"`
	Score        float32   `json:"score"`
	UpdatedAt    time.Time `json:"-"`
}
type JustSubmissionResponse struct {
	ID        uint      `json:"id"`
	QuizID    uint      `json:"quiz_id"`
	UserID    uint      `json:"user_id"`
	Score     float32   `json:"score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SubmissionResponse struct {
	ID        uint                       `json:"id"`
	QuizID    uint                       `json:"quiz_id"`
	UserID    uint                       `json:"user_id"`
	Score     float32                    `json:"score"`
	CreatedAt time.Time                  `json:"created_at"`
	UpdatedAt time.Time                  `json:"updated_at"`
	Answers   []SubmissionAnswerResponse `json:"answer"`
}

type SubmissionAnswerResponse struct {
	QuestionID    uint `json:"question_id"`
	CorrectAnswer uint `json:"correct_id"`
	AnswerUser    uint `json:"answer_id"`
	IsCorrect     bool `json:"is_correct"`
}
