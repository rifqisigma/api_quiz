package usecase

import (
	"api_quiz/dto"
	"api_quiz/internal/repository"
	"api_quiz/utils/helper"
)

type QuizUseCase interface {

	//quiz
	GetAllQuiz() ([]dto.JustQuizResponse, error)
	GetQuizFromId(quizId uint) (*dto.QuizResponseWithQS, error)
	CreateQuiz(input *dto.Quiz) (*dto.JustQuizResponse, error)
	UpdateQuiz(input *dto.UpdatedQuiz, userId uint) (*dto.JustQuizResponse, error)
	DeleteQuiz(userId, quizId uint) error

	//question
	GetQuestionAnswerByQuizId(quizId uint) ([]dto.QuestionResponse, error)
	GetQuestionById(questionId, quizId uint) (*dto.QuestionResponse, error)
	CreateQuestionAndAnswer(inputQuestion *dto.Question, userId uint) (*dto.QuestionResponse, error)
	UpdateQuestion(input *dto.QuestionUpdate, userId uint) (*dto.JustQuestionResponse, error)
	DeleteQuestion(questionId, quizId, userId uint) error

	//answer
	GetAnswerByQuestionId(questionId uint) ([]dto.AnswerResponse, error)
	UpdateAnswer(userId, quizId uint, input dto.Answer) ([]dto.AnswerResponse, error)
	DeleteAnswer(answerId, questionId, quizId, userId uint) error
	AddAnswer(userId, quizId uint, input []dto.Answer) ([]dto.AnswerResponse, error)
}

type quizUseCase struct {
	quizRepo repository.QuizRepository
}

func NewQuizUseCase(quizRepo repository.QuizRepository) QuizUseCase {
	return &quizUseCase{quizRepo}
}

func (u *quizUseCase) GetAllQuiz() ([]dto.JustQuizResponse, error) {
	return u.quizRepo.GetAllQuiz()
}

func (u *quizUseCase) GetQuizFromId(quizId uint) (*dto.QuizResponseWithQS, error) {
	return u.quizRepo.GetQuizById(quizId)
}

func (u *quizUseCase) CreateQuiz(input *dto.Quiz) (*dto.JustQuizResponse, error) {
	result, err := u.quizRepo.CreateQuiz(input)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (u *quizUseCase) UpdateQuiz(input *dto.UpdatedQuiz, userId uint) (*dto.JustQuizResponse, error) {
	valid, err := u.quizRepo.IsCreator(userId, input.ID)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, helper.ErrUnauhorized
	}

	result, err := u.quizRepo.UpdateQuiz(input, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *quizUseCase) DeleteQuiz(userId, quizId uint) error {
	valid, err := u.quizRepo.IsCreator(userId, quizId)
	if err != nil {
		return err
	}
	if !valid {
		return helper.ErrUnauhorized
	}

	if err := u.quizRepo.DeleteQuiz(quizId); err != nil {
		return err
	}

	return nil
}

// question
func (u *quizUseCase) GetQuestionAnswerByQuizId(quizId uint) ([]dto.QuestionResponse, error) {
	return u.quizRepo.GetQuestionAnswerByQuizId(quizId)
}

func (u *quizUseCase) GetQuestionById(questionId, quizId uint) (*dto.QuestionResponse, error) {
	return u.quizRepo.GetQuestionById(questionId, quizId)
}

func (u *quizUseCase) CreateQuestionAndAnswer(inputQuestion *dto.Question, userId uint) (*dto.QuestionResponse, error) {

	valid, err := u.quizRepo.IsCreator(userId, inputQuestion.QuizID)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, helper.ErrUnauhorized
	}

	result, err := u.quizRepo.CreateQuestionAndAnswer(inputQuestion)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (u *quizUseCase) UpdateQuestion(input *dto.QuestionUpdate, userId uint) (*dto.JustQuestionResponse, error) {
	valid, err := u.quizRepo.IsCreator(userId, input.QuizID)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, helper.ErrUnauhorized
	}

	result, err := u.quizRepo.UpdateQuestion(input)
	if err != nil {
		return nil, err
	}
	return result, err
}

func (u *quizUseCase) DeleteQuestion(questionId, quizId, userId uint) error {
	valid, err := u.quizRepo.IsCreator(userId, quizId)
	if err != nil {
		return err
	}
	if !valid {
		return helper.ErrUnauhorized
	}

	if err := u.quizRepo.DeleteQuestion(quizId, questionId); err != nil {
		return err
	}

	return nil
}

// answer
func (u *quizUseCase) GetAnswerByQuestionId(questionId uint) ([]dto.AnswerResponse, error) {
	return u.quizRepo.GetAnswerByQuestionId(questionId)
}

func (u *quizUseCase) UpdateAnswer(userId, quizId uint, input dto.Answer) ([]dto.AnswerResponse, error) {
	valid, err := u.quizRepo.IsCreator(userId, quizId)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, helper.ErrUnauhorized
	}

	result, err := u.quizRepo.UpdateAnswer(input)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (u *quizUseCase) AddAnswer(userId, quizId uint, input []dto.Answer) ([]dto.AnswerResponse, error) {
	valid, err := u.quizRepo.IsCreator(userId, quizId)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, helper.ErrUnauhorized
	}

	result, err := u.quizRepo.AddAnswer(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (u *quizUseCase) DeleteAnswer(answerId, questionId, quizId, userId uint) error {
	valid, err := u.quizRepo.IsCreator(userId, quizId)
	if err != nil {
		return err
	}
	if !valid {
		return helper.ErrUnauhorized
	}

	if err := u.quizRepo.DeleteAnswer(answerId, questionId); err != nil {
		return err
	}

	return nil
}
