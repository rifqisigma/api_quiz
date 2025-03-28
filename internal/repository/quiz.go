package repository

import (
	"api_quiz/dto"
	"api_quiz/entity"
	"api_quiz/utils/helper"

	"gorm.io/gorm"
)

type QuizRepository interface {
	//quiz
	GetAllQuiz() ([]dto.JustQuizResponse, error)
	GetQuizById(quizId uint) (*dto.QuizResponseWithQS, error)
	CreateQuiz(input *dto.Quiz) (*dto.JustQuizResponse, error)
	IsCreator(userId, quizId uint) (bool, error)
	UpdateQuiz(input *dto.UpdatedQuiz, userId uint) (*dto.JustQuizResponse, error)
	DeleteQuiz(quizId uint) error

	//question
	GetQuestionAnswerByQuizId(quizId uint) ([]dto.QuestionResponse, error)
	GetQuestionById(questionId, quizId uint) (*dto.QuestionResponse, error)
	CreateQuestionAndAnswer(inputQuestion *dto.Question) (*dto.QuestionResponse, error)
	UpdateQuestion(input *dto.QuestionUpdate) (*dto.JustQuestionResponse, error)
	DeleteQuestion(quizId, questionId uint) error

	//answer
	GetAnswerByQuestionId(questionId uint) ([]dto.AnswerResponse, error)
	CheckTotalAnswer(questionId uint) (int64, error)
	CheckTotalAnswerIscorrect(questionId uint) (int64, error)
	UpdateAnswer(input dto.Answer) ([]dto.AnswerResponse, error)
	AddAnswer(input []dto.Answer) ([]dto.AnswerResponse, error)
	DeleteAnswer(answerId, questionId uint) error
}

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &quizRepository{db}
}

// quiz
func (r *quizRepository) GetAllQuiz() ([]dto.JustQuizResponse, error) {
	var quiz []entity.Quiz
	if err := r.db.Select("id,title, creator_id").Find(&quiz).Error; err != nil {
		return nil, err
	}

	var response []dto.JustQuizResponse
	for _, q := range quiz {
		response = append(response, dto.JustQuizResponse{
			ID:      q.ID,
			Creator: q.CreatorID,
			Title:   q.Title,
		})
	}

	return response, nil
}
func (r *quizRepository) GetQuizById(quizId uint) (*dto.QuizResponseWithQS, error) {
	var quiz entity.Quiz
	if err := r.db.Preload("Questions.Answers").Where("id = ?", quizId).First(&quiz).Error; err != nil {
		return nil, err
	}

	var questions []dto.QuestionResponse

	for _, q := range quiz.Questions {
		var answers []dto.AnswerResponse
		for _, ans := range q.Answers {
			answers = append(answers, dto.AnswerResponse{
				ID:        ans.ID,
				Text:      ans.Text,
				IsCorrect: ans.IsCorrect,
			})
		}

		questions = append(questions, dto.QuestionResponse{
			ID:     q.ID,
			QuizID: quizId,
			Text:   q.Text,
			Answer: answers,
		})
	}

	response := dto.QuizResponseWithQS{
		ID:       quiz.ID,
		Creator:  *quiz.CreatorID,
		Title:    quiz.Title,
		Question: questions,
	}

	return &response, nil
}

func (r *quizRepository) CreateQuiz(input *dto.Quiz) (*dto.JustQuizResponse, error) {

	quiz := entity.Quiz{
		Title:     input.Title,
		CreatorID: &input.Creator,
	}

	result := r.db.Create(&quiz)
	if result.Error != nil {
		return nil, result.Error
	}
	response := dto.JustQuizResponse{
		ID:      quiz.ID,
		Creator: quiz.CreatorID,
		Title:   quiz.Title,
	}

	return &response, nil
}

func (r *quizRepository) IsCreator(userId, quizId uint) (bool, error) {
	var quiz entity.Quiz
	if err := r.db.Where("id = ? AND creator_id = ?", quizId, userId).First(&quiz).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (r *quizRepository) UpdateQuiz(input *dto.UpdatedQuiz, userId uint) (*dto.JustQuizResponse, error) {
	updated := r.db.Model(&entity.Quiz{}).Where("id = ?", input.ID).Update("title", input.Title)
	if updated.Error != nil {
		return nil, updated.Error
	}
	if updated.RowsAffected == 0 {
		return nil, helper.ErrQuizNotFound
	}

	response := dto.JustQuizResponse{
		ID:      input.ID,
		Creator: &userId,
		Title:   input.Title,
	}

	return &response, nil
}

func (r *quizRepository) DeleteQuiz(quizId uint) error {

	deleted := r.db.Where("id = ?", quizId).Delete(&entity.Quiz{})
	if deleted.Error != nil {
		return deleted.Error
	}

	if deleted.RowsAffected == 0 {
		return helper.ErrUnauhorized
	}

	return nil
}

// question of quizz
func (r *quizRepository) GetQuestionAnswerByQuizId(quizId uint) ([]dto.QuestionResponse, error) {
	var question []entity.Question
	if err := r.db.Model(&entity.Question{}).Preload("Answers").Where("quiz_id = ?", quizId).Find(&question).Error; err != nil {
		return nil, err
	}

	questionResponse := make([]dto.QuestionResponse, len(question))
	for i, q := range question {
		answers := make([]dto.AnswerResponse, len(q.Answers))
		for i, ans := range q.Answers {
			answers[i] = dto.AnswerResponse{
				ID:         ans.ID,
				QuestionID: ans.QuestionID,
				Text:       ans.Text,
				IsCorrect:  ans.IsCorrect,
			}
		}
		questionResponse[i] = dto.QuestionResponse{
			ID:     q.ID,
			QuizID: quizId,
			Text:   q.Text,
			Answer: answers,
		}
	}

	return questionResponse, nil
}

func (r *quizRepository) GetQuestionById(questionId, quizId uint) (*dto.QuestionResponse, error) {
	var question entity.Question
	if err := r.db.Model(&entity.Question{}).Preload("Answers").Where("id = ? AND quiz_id = ? ", questionId, quizId).First(&question).Error; err != nil {
		return nil, err
	}

	responseAnswer := make([]dto.AnswerResponse, len(question.Answers))
	for i, ans := range question.Answers {
		responseAnswer[i] = dto.AnswerResponse{
			ID:         ans.ID,
			QuestionID: ans.QuestionID,
			Text:       ans.Text,
			IsCorrect:  ans.IsCorrect,
		}
	}

	response := dto.QuestionResponse{
		ID:     question.ID,
		QuizID: question.QuizID,
		Text:   question.Text,
		Answer: responseAnswer,
	}

	return &response, nil
}
func (r *quizRepository) CreateQuestionAndAnswer(inputQuestion *dto.Question) (*dto.QuestionResponse, error) {
	question := entity.Question{
		QuizID: inputQuestion.QuizID,
		Text:   inputQuestion.Text,
	}
	tx := r.db.Begin()

	if err := tx.Create(&question).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	answers := make([]entity.Answer, 0, len(inputQuestion.Answers))
	responseAnswers := make([]dto.AnswerResponse, 0, len(inputQuestion.Answers))

	for _, ans := range inputQuestion.Answers {
		answerEntity := entity.Answer{
			QuestionID: question.ID,
			Text:       ans.Text,
			IsCorrect:  ans.IsCorrect,
		}

		answers = append(answers, answerEntity)

		responseAnswers = append(responseAnswers, dto.AnswerResponse{
			ID:         answerEntity.ID,
			QuestionID: answerEntity.QuestionID,
			Text:       answerEntity.Text,
			IsCorrect:  answerEntity.IsCorrect,
		})
	}

	if err := tx.Create(&answers).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	response := dto.QuestionResponse{
		ID:     question.ID,
		QuizID: question.QuizID,
		Text:   question.Text,
		Answer: responseAnswers,
	}

	return &response, nil

}

func (r *quizRepository) DeleteQuestion(quizId, questionId uint) error {
	delete := r.db.Where("id = ? AND quiz_id = ?", questionId, quizId).Delete(&entity.Question{})
	if delete.Error != nil {
		return delete.Error
	}
	if delete.RowsAffected == 0 {
		return helper.ErrQuestionNotFound
	}

	return nil
}

func (r *quizRepository) UpdateQuestion(input *dto.QuestionUpdate) (*dto.JustQuestionResponse, error) {
	question := entity.Question{
		ID:     input.ID,
		QuizID: input.QuizID,
		Text:   input.Text,
	}

	updated := r.db.Model(&entity.Question{}).Where("id = ? AND quiz_id = ? ", input.ID, input.QuizID).Updates(&question)
	if updated.Error != nil {
		return nil, updated.Error
	}
	if updated.RowsAffected == 0 {
		return nil, helper.ErrQuestionNotFound
	}

	response := dto.JustQuestionResponse{
		ID:     question.ID,
		QuizID: question.QuizID,
		Text:   question.Text,
	}
	return &response, nil
}

// answer of quizz
func (r *quizRepository) GetAnswerByQuestionId(questionId uint) ([]dto.AnswerResponse, error) {
	var answer []entity.Answer
	if err := r.db.Model(&entity.Answer{}).Where("question_id = ?", questionId).Find(&answer).Error; err != nil {
		return nil, err
	}

	response := make([]dto.AnswerResponse, len(answer))
	for i, ans := range answer {
		response[i] = dto.AnswerResponse{
			ID:         ans.ID,
			QuestionID: ans.QuestionID,
			Text:       ans.Text,
			IsCorrect:  ans.IsCorrect,
		}
	}
	return response, nil
}

func (r *quizRepository) UpdateAnswer(input dto.Answer) ([]dto.AnswerResponse, error) {
	tx := r.db.Begin()

	if err := tx.Model(&entity.Answer{}).
		Where("id = ? AND question_id = ?", input.ID, input.QuestionID).
		Select("text", "is_correct").
		Updates(entity.Answer{
			Text:      input.Text,
			IsCorrect: input.IsCorrect,
		}).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	correctAnswers, err := r.CheckTotalAnswerIscorrect(input.QuestionID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if correctAnswers != 1 {
		tx.Rollback()
		return nil, helper.ErrCorrectAnswer
	}

	tx.Commit()

	var answers []entity.Answer
	if err := r.db.Where("question_id = ?", input.QuestionID).Find(&answers).Error; err != nil {
		return nil, err
	}

	response := make([]dto.AnswerResponse, 0, len(answers))
	for _, ans := range answers {
		response = append(response, dto.AnswerResponse{
			ID:         ans.ID,
			QuestionID: ans.QuestionID,
			Text:       ans.Text,
			IsCorrect:  ans.IsCorrect,
		})
	}

	return response, nil
}

func (r *quizRepository) AddAnswer(input []dto.Answer) ([]dto.AnswerResponse, error) {
	tx := r.db.Begin()
	questionID := input[0].QuestionID

	totalAnswers, err := r.CheckTotalAnswer(questionID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	totalCorrect, err := r.CheckTotalAnswerIscorrect(questionID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	newCorrectCount := 0
	for _, ans := range input {
		if ans.IsCorrect {
			newCorrectCount++
		}
	}

	if totalAnswers+int64(len(input)) > 5 {
		tx.Rollback()
		return nil, helper.ErrToomuchAnswer
	}

	if totalCorrect+int64(newCorrectCount) != 1 {
		tx.Rollback()
		return nil, helper.ErrCorrectAnswer
	}

	answers := make([]entity.Answer, len(input))
	for i, ans := range input {
		answers[i] = entity.Answer{
			QuestionID: ans.QuestionID,
			Text:       ans.Text,
			IsCorrect:  ans.IsCorrect,
		}
	}

	if err := tx.Create(&answers).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	var ResponseAnswers []entity.Answer
	if err := r.db.Where("question_id = ?", questionID).Find(&ResponseAnswers).Error; err != nil {
		return nil, err
	}

	response := make([]dto.AnswerResponse, len(ResponseAnswers))
	for i, ans := range ResponseAnswers {
		response[i] = dto.AnswerResponse{
			ID:         ans.ID,
			QuestionID: ans.QuestionID,
			Text:       ans.Text,
			IsCorrect:  ans.IsCorrect,
		}
	}

	return response, nil
}

func (r *quizRepository) DeleteAnswer(answerId, questionId uint) error {
	tx := r.db.Begin()

	if err := tx.Where("id = ? AND  question_id = ?", answerId, questionId).Delete(&entity.Answer{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	totalAnswers, err := r.CheckTotalAnswer(questionId)
	if err != nil {
		tx.Rollback()
		return err
	}

	if totalAnswers <= 2 {
		tx.Rollback()
		return helper.ErrAnswerNotEnough
	}

	return tx.Commit().Error
}

func (r *quizRepository) CheckTotalAnswer(questionId uint) (int64, error) {
	var total int64
	if err := r.db.Model(&entity.Answer{}).Where("question_id = ?", questionId).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *quizRepository) CheckTotalAnswerIscorrect(questionId uint) (int64, error) {
	var total int64
	if err := r.db.Model(&entity.Answer{}).Where("question_id = ? AND is_correct = ?", questionId, true).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}
