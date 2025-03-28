package route

import (
	"api_quiz/internal/handler"
	"api_quiz/utils/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(authHandler *handler.AuthHandler, quizHandler *handler.QuizHandler, submissionHandler *handler.SubmissionHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/verification", authHandler.Verification).Methods(http.MethodGet)

	userRoute := r.PathPrefix("/user").Subrouter()
	userRoute.Use(middleware.JWTAuthMiddleware)

	userRoute.HandleFunc("/delete", authHandler.DeleteUser).Methods(http.MethodDelete)

	//quiz
	quizRoute := r.PathPrefix("/quiz").Subrouter()
	quizRoute.Use(middleware.JWTAuthMiddleware)

	//quiz
	quizRoute.HandleFunc("/get", quizHandler.GetAllQuiz).Methods(http.MethodGet)
	quizRoute.HandleFunc("/get/{quizid}", quizHandler.GetQuizById).Methods(http.MethodGet)
	quizRoute.HandleFunc("/create", quizHandler.CreateQuiz).Methods(http.MethodPost)
	quizRoute.HandleFunc("/update/{quizid}", quizHandler.CreateQuiz).Methods(http.MethodPut)
	quizRoute.HandleFunc("/delete/{quizid}", quizHandler.DeleteQuiz).Methods(http.MethodDelete)
	//question
	quizRoute.HandleFunc("/{quizid}/get/question", quizHandler.GetQuestionAnswerByQuizId).Methods(http.MethodGet)
	quizRoute.HandleFunc("/{quizid}/get/question/{questionid}", quizHandler.GetQuestionById).Methods(http.MethodGet)
	quizRoute.HandleFunc("/{quizid}/question/create", quizHandler.CreateQuestionAndAnswer).Methods(http.MethodPost)
	quizRoute.HandleFunc("/{quizid}/question/{questionid}/delete", quizHandler.DeleteQuestion).Methods(http.MethodDelete)
	quizRoute.HandleFunc("/{quizid}/question/{questionid}/update", quizHandler.UpdateQuestion).Methods(http.MethodPut)
	//answer
	quizRoute.HandleFunc("/get/question{questionid}", quizHandler.GetAnswerByQuestionId).Methods(http.MethodGet)
	quizRoute.HandleFunc("/{quizid}/question/{questionid}/answer/{answerid}/update", quizHandler.UpdateAnswer).Methods(http.MethodPut)
	quizRoute.HandleFunc("/{quizid}/question/{questionid}/answer/add", quizHandler.AddAnswer).Methods(http.MethodPost)
	quizRoute.HandleFunc("/{quizid}/question/{questionid}/answer/{answerid}/delete", quizHandler.DeleteAnswer).Methods(http.MethodDelete)

	//submission
	submissionRoute := r.PathPrefix("/submission").Subrouter()
	submissionRoute.Use(middleware.JWTAuthMiddleware)

	submissionRoute.HandleFunc("/get", submissionHandler.GetAllSubmission).Methods(http.MethodGet)
	submissionRoute.HandleFunc("/get/{submissionid}", submissionHandler.GetSubmissionById).Methods(http.MethodGet)
	submissionRoute.HandleFunc("/create/{quizid}", submissionHandler.CreateSubmission).Methods(http.MethodPost)
	submissionRoute.HandleFunc("/update/{submissionid}", submissionHandler.UpdateSubmission).Methods(http.MethodPut)
	submissionRoute.HandleFunc("/delete/{submissionid}", submissionHandler.DeleteSubmission).Methods(http.MethodDelete)

	return r

}
