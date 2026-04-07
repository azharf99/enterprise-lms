package usecase

import (
	"errors"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/datatypes"
)

type questionUsecase struct {
	questionRepo domain.QuizQuestionRepository
}

func NewQuizQuestionUsecase(qr domain.QuizQuestionRepository) domain.QuizQuestionUsecase {
	return &questionUsecase{questionRepo: qr}
}

func (u *questionUsecase) CreateQuestion(quizID uint, qType domain.QuestionType, text string, options, correctAnswer datatypes.JSON, points int, explanation string) (*domain.Question, error) {
	if text == "" {
		return nil, errors.New("teks soal tidak boleh kosong")
	}

	question := &domain.Question{
		QuizID:        quizID,
		Type:          qType,
		Text:          text,
		Options:       options,
		CorrectAnswer: correctAnswer,
		Points:        points,
		Explanation:   explanation,
	}

	if err := u.questionRepo.CreateQuizQuestion(question); err != nil {
		return nil, err
	}
	return question, nil
}

func (u *questionUsecase) GetQuestionsByQuizID(quizID uint, isRandomized bool) ([]domain.Question, error) {
	return u.questionRepo.GetQuizQuestionsByQuizID(quizID, isRandomized)
}

func (u *questionUsecase) GetQuestionByID(id uint) (domain.Question, error) {
	return u.questionRepo.GetQuizQuestionByID(id)
}

func (u *questionUsecase) UpdateQuestion(id uint, qType domain.QuestionType, text string, options, correctAnswer datatypes.JSON, points int, explanation string) (*domain.Question, error) {
	question, err := u.questionRepo.GetQuizQuestionByID(id)
	if err != nil {
		return nil, errors.New("pertanyaan tidak ditemukan")
	}

	question.Type = qType
	question.Text = text
	question.Options = options
	question.CorrectAnswer = correctAnswer
	question.Points = points
	question.Explanation = explanation

	if err := u.questionRepo.UpdateQuizQuestion(&question); err != nil {
		return nil, err
	}
	return &question, nil
}

func (u *questionUsecase) DeleteQuestion(id uint) error {
	return u.questionRepo.DeleteQuizQuestion(id)
}
