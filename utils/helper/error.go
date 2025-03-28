package helper

import "errors"

var (

	//auth
	ErrUserNotFound = errors.New("uset not found")
	ErrServerError  = errors.New("server error")
	ErrInvalidEmail = errors.New("invalid email")
	ErrUnauhorized  = errors.New("you unauthorized for this action")

	//quiz
	ErrQuizNotFound     = errors.New("quiz not found")
	ErrQuestionNotFound = errors.New("question not found")
	ErrAnswerNotEnough  = errors.New("answer must 2 or more")
	ErrCorrectAnswer    = errors.New("correct answer just only 1 ")
	ErrToomuchAnswer    = errors.New("answer max is 5")

	//submission
	ErrSubmissionNotFound = errors.New("submission not found")
)
