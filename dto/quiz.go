package dto

type Quiz struct {
	Creator uint   `json:"-"`
	Title   string `json:"title"`
}

type UpdatedQuiz struct {
	ID    uint   `json:"-"`
	Title string `json:"title"`
}

type JustQuizResponse struct {
	ID      uint   `json:"id"`
	Creator *uint  `json:"creator"`
	Title   string `json:"title"`
}

type QuizResponseWithQS struct {
	ID       uint               `json:"id"`
	Creator  uint               `json:"creator"`
	Title    string             `json:"title"`
	Question []QuestionResponse `json:"question"`
}

// question
type Question struct {
	ID      uint     `json:"-"`
	QuizID  uint     `json:"-"`
	Text    string   `json:"text"`
	Answers []Answer `json:"answer"`
}

type QuestionUpdate struct {
	QuizID uint   `json:"-"`
	ID     uint   `json:"-"`
	Text   string `json:"text"`
}

type JustQuestionResponse struct {
	ID     uint   `json:"id"`
	QuizID uint   `json:"quiz_id"`
	Text   string `json:"text"`
}
type QuestionResponse struct {
	ID     uint             `json:"ID"`
	QuizID uint             `json:"quiz_id"`
	Text   string           `json:"text"`
	Answer []AnswerResponse `json:"answer"`
}

// answer
type Answer struct {
	ID         uint   `json:"-"`
	QuestionID uint   `json:"-"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}

type AnswerResponse struct {
	ID         uint   `json:"id"`
	QuestionID uint   `json:"question_id"`
	Text       string `json:"text"`
	IsCorrect  bool   `json:"is_correct"`
}
