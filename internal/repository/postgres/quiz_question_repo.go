package postgres

import (
	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

func NewQuizQuestionRepository(db *gorm.DB) domain.QuizQuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) CreateQuizQuestion(question *domain.Question) error {
	return r.db.Create(question).Error
}

func (r *questionRepository) GetQuizQuestionsByQuizID(quizID uint, isRandomized bool) ([]domain.Question, error) {
	var questions []domain.Question
	db := r.db.Where("quiz_id = ?", quizID)

	if isRandomized {
		db = db.Order("RANDOM()")
	} else {
		db = db.Order("id asc")
	}

	err := db.Find(&questions).Error
	return questions, err
}

func (r *questionRepository) GetQuizQuestionByID(id uint) (domain.Question, error) {
	var question domain.Question
	err := r.db.First(&question, id).Error
	return question, err
}

func (r *questionRepository) UpdateQuizQuestion(question *domain.Question) error {
	return r.db.Model(question).Updates(question).Error
}

func (r *questionRepository) DeleteQuizQuestion(id uint) error {
	return r.db.Delete(&domain.Question{}, id).Error
}
