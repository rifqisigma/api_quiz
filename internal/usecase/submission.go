package usecase

import (
	"api_quiz/dto"
	"api_quiz/internal/repository"
	"api_quiz/utils/helper"
)

type SubmissionUseCase interface {
	GetAllSubmission() ([]dto.JustSubmissionResponse, error)
	GetSubmissionById(submissionId uint) (*dto.SubmissionResponse, error)
	CreateSubmission(input *dto.Submission) (*dto.SubmissionResponse, error)
	UpdateSubmision(input *dto.SubmissionUpdate, userId uint) (*dto.JustSubmissionResponse, error)
	DeleteSubmision(submissionId, userId uint) error
}

type submissionUseCase struct {
	submissionRepo repository.SubmissionRepository
	quizRepo       repository.QuizRepository
}

func NewSubmissionUseCase(submissionRepo repository.SubmissionRepository, quizRepo repository.QuizRepository) SubmissionUseCase {
	return &submissionUseCase{submissionRepo, quizRepo}
}

func (u *submissionUseCase) GetAllSubmission() ([]dto.JustSubmissionResponse, error) {
	return u.submissionRepo.GetAllSubmission()
}

func (u *submissionUseCase) GetSubmissionById(submissionId uint) (*dto.SubmissionResponse, error) {
	return u.submissionRepo.GetSubmissionById(submissionId)
}

func (u *submissionUseCase) CreateSubmission(input *dto.Submission) (*dto.SubmissionResponse, error) {
	return u.submissionRepo.CreateSubmission(input)
}

func (u *submissionUseCase) UpdateSubmision(input *dto.SubmissionUpdate, userId uint) (*dto.JustSubmissionResponse, error) {
	quizId, err := u.submissionRepo.GetQuizIdFromSubmisionId(input.SubmissionID)
	if err != nil {
		return nil, err
	}
	if quizId == 0 {
		return nil, helper.ErrQuizNotFound
	}

	valid, err := u.quizRepo.IsCreator(userId, quizId)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, helper.ErrUnauhorized
	}

	return u.submissionRepo.UpdateSubmission(input)
}

func (u *submissionUseCase) DeleteSubmision(submissionId, userId uint) error {
	quizId, err := u.submissionRepo.GetQuizIdFromSubmisionId(submissionId)
	if err != nil {
		return err
	}
	if quizId == 0 {
		return helper.ErrQuizNotFound
	}

	valid, err := u.quizRepo.IsCreator(userId, quizId)
	if err != nil {
		return err
	}
	if !valid {
		return helper.ErrUnauhorized
	}

	return u.submissionRepo.DeleteSubmission(submissionId)
}
