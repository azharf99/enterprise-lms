package usecase

import (
	"errors"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/datatypes"
)

type examQuestionUsecase struct {
	examQuestionRepo domain.ExamQuestionRepository
}

func NewExamQuestionUsecase(eqr domain.ExamQuestionRepository) domain.ExamQuestionUsecase {
	return &examQuestionUsecase{examQuestionRepo: eqr}
}

func (u *examQuestionUsecase) CreateExamQuestion(examID uint, qType domain.QuestionType, text string, options, correctAnswer datatypes.JSON, points int, explanation string) (*domain.ExamQuestion, error) {
	exam := &domain.ExamQuestion{
		ExamID:        examID,
		Type:          qType,
		Text:          text,
		Options:       options,
		CorrectAnswer: correctAnswer,
		Points:        points,
		Explanation:   explanation,
	}
	if err := u.examQuestionRepo.CreateExamQuestion(exam); err != nil {
		return nil, err
	}
	return exam, nil
}

func (u *examQuestionUsecase) GetExamQuestionsByExamID(examID uint, isRandomized bool) ([]domain.ExamQuestion, error) {
	return u.examQuestionRepo.GetExamQuestionsByExamID(examID, isRandomized)
}

func (u *examQuestionUsecase) GetExamQuestionByID(id uint) (domain.ExamQuestion, error) {
	return u.examQuestionRepo.GetExamQuestionByID(id)
}

func (u *examQuestionUsecase) UpdateExamQuestion(id uint, qType domain.QuestionType, text string, options, correctAnswer datatypes.JSON, points int, explanation string) (*domain.ExamQuestion, error) {
	question, err := u.examQuestionRepo.GetExamQuestionByID(id)
	if err != nil {
		return nil, errors.New("pertanyaan tidak ditemukan")
	}

	question.Type = qType
	question.Text = text
	question.Options = options
	question.CorrectAnswer = correctAnswer
	question.Points = points
	question.Explanation = explanation

	if err := u.examQuestionRepo.UpdateExamQuestion(&question); err != nil {
		return nil, err
	}
	return &question, nil
}

func (u *examQuestionUsecase) DeleteExamQuestion(id uint) error {
	return u.examQuestionRepo.DeleteExamQuestion(id)
}
