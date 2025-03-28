package repository

import (
	"api_quiz/dto"
	"api_quiz/entity"
	"api_quiz/utils/helper"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SubmissionRepository interface {
	GetAllSubmission() ([]dto.JustSubmissionResponse, error)
	GetSubmissionById(submissionId uint) (*dto.SubmissionResponse, error)
	CreateSubmission(input *dto.Submission) (*dto.SubmissionResponse, error)
	GetQuizIdFromSubmisionId(submisionId uint) (uint, error)
	UpdateSubmission(input *dto.SubmissionUpdate) (*dto.JustSubmissionResponse, error)
	DeleteSubmission(submissionId uint) error
}

type submissionRepository struct {
	db *gorm.DB
}

func NewSubmissionRepository(db *gorm.DB) SubmissionRepository {
	return &submissionRepository{db}
}

func (r *submissionRepository) GetAllSubmission() ([]dto.JustSubmissionResponse, error) {
	var submission []entity.Submission
	if err := r.db.Model(&entity.Submission{}).Find(&submission).Error; err != nil {
		return nil, err
	}

	response := make([]dto.JustSubmissionResponse, len(submission))
	for i, s := range submission {
		response[i] = dto.JustSubmissionResponse{
			ID:        s.ID,
			QuizID:    s.QuizID,
			UserID:    *s.UserID,
			Score:     s.Score,
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		}
	}

	return response, nil
}

func (r *submissionRepository) GetSubmissionById(submissionId uint) (*dto.SubmissionResponse, error) {
	var submission entity.Submission

	if err := r.db.Preload("Answers").Where("id = ?", submissionId).First(&submission).Error; err != nil {
		return nil, err
	}

	answers := make([]dto.SubmissionAnswerResponse, len(submission.Answers))
	for i, ans := range submission.Answers {
		answers[i] = dto.SubmissionAnswerResponse{
			QuestionID:    ans.QuestionID,
			AnswerUser:    ans.UserAnswerID,
			CorrectAnswer: ans.CorrectID,
			IsCorrect:     ans.IsCorrect,
		}
	}

	response := dto.SubmissionResponse{
		ID:        submissionId,
		QuizID:    submission.QuizID,
		UserID:    *submission.UserID,
		Score:     submission.Score,
		CreatedAt: submission.CreatedAt,
		UpdatedAt: submission.UpdatedAt,
		Answers:   answers,
	}

	return &response, nil

}

func (r *submissionRepository) CreateSubmission(input *dto.Submission) (*dto.SubmissionResponse, error) {
	tx := r.db.Begin()

	var questions []entity.Question
	if err := tx.Preload("Answers").Where("quiz_id = ?", input.QuizID).Find(&questions).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	missingQuestions := []uint{}
	userAnswers := make(map[uint]uint)
	for _, ans := range input.Answers {
		userAnswers[ans.QuestionID] = ans.AnswerID
	}

	for _, question := range questions {
		if _, ok := userAnswers[question.ID]; !ok {
			missingQuestions = append(missingQuestions, question.ID)
		}
	}

	if len(missingQuestions) > 0 {
		tx.Rollback()
		return nil, fmt.Errorf("pertanyaan belum dijawab: %v", missingQuestions)
	}

	var correctCount int
	var submissionAnswers []entity.SubmissionUserAnswer

	for _, question := range questions {
		var correctAnswerID uint
		for _, ans := range question.Answers {
			if ans.IsCorrect {
				correctAnswerID = ans.ID
				break
			}
		}

		userAnswerID := userAnswers[question.ID]
		isCorrect := userAnswerID == correctAnswerID
		if isCorrect {
			correctCount++
		}

		submissionAnswers = append(submissionAnswers, entity.SubmissionUserAnswer{
			SubmissionID: 0,
			QuestionID:   question.ID,
			UserAnswerID: userAnswerID,
			CorrectID:    correctAnswerID,
			IsCorrect:    isCorrect,
		})
	}

	submission := entity.Submission{
		QuizID: input.QuizID,
		UserID: &input.UserID,
		Score:  float32(correctCount) / float32(len(questions)) * 100,
	}

	if err := tx.Create(&submission).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	for i := range submissionAnswers {
		submissionAnswers[i].SubmissionID = submission.ID
	}

	if err := tx.Create(&submissionAnswers).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	responseAnswer := make([]dto.SubmissionAnswerResponse, len(submissionAnswers))
	for i, answer := range submissionAnswers {
		responseAnswer[i] = dto.SubmissionAnswerResponse{
			QuestionID:    answer.QuestionID,
			AnswerUser:    answer.UserAnswerID,
			CorrectAnswer: answer.CorrectID,
			IsCorrect:     answer.IsCorrect,
		}
	}

	response := dto.SubmissionResponse{
		ID:        submission.ID,
		QuizID:    submission.QuizID,
		UserID:    *submission.UserID,
		Score:     submission.Score,
		CreatedAt: submission.CreatedAt,
		UpdatedAt: submission.UpdatedAt,
		Answers:   responseAnswer,
	}

	return &response, nil
}

func (r *submissionRepository) GetQuizIdFromSubmisionId(submisionId uint) (uint, error) {
	var quizId uint
	if err := r.db.Model(&entity.Submission{}).Select("quiz_id").Where("id  = ?", submisionId).First(&quizId).Error; err != nil {
		return 0, err
	}

	return quizId, nil
}

func (r *submissionRepository) UpdateSubmission(input *dto.SubmissionUpdate) (*dto.JustSubmissionResponse, error) {
	updated := r.db.Model(&entity.Submission{}).Where("id = ?", input.SubmissionID).Updates(map[string]interface{}{"score": input.Score, "updated_at": time.Now()})
	if updated.Error != nil {
		return nil, updated.Error
	}
	if updated.RowsAffected == 0 {
		return nil, helper.ErrSubmissionNotFound
	}

	var parsingResponse entity.Submission
	if err := r.db.Where("id = ?", input.SubmissionID).First(&parsingResponse).Error; err != nil {
		return nil, err
	}

	response := dto.JustSubmissionResponse{
		ID:        parsingResponse.ID,
		Score:     parsingResponse.Score,
		UpdatedAt: parsingResponse.UpdatedAt,
	}

	return &response, nil
}

func (r *submissionRepository) DeleteSubmission(submissionId uint) error {
	result := r.db.Where("id = ?", submissionId).Delete(&entity.Submission{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return helper.ErrSubmissionNotFound
	}

	return nil
}
