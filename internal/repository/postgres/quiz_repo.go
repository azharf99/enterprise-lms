package postgres

import (
	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/gorm"
)

type quizRepository struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) domain.QuizRepository {
	return &quizRepository{db: db}
}

func (r *quizRepository) CreateQuiz(quiz *domain.Quiz) error {
	return r.db.Create(quiz).Error
}

func (r *quizRepository) GetQuizzesByModuleID(moduleID uint) ([]domain.Quiz, error) {
	var quizzes []domain.Quiz
	err := r.db.Where("module_id = ?", moduleID).Find(&quizzes).Error
	return quizzes, err
}

func (r *quizRepository) GetQuizByID(id uint) (domain.Quiz, error) {
	var quiz domain.Quiz
	err := r.db.Preload("Questions").First(&quiz, id).Error
	return quiz, err
}

func (r *quizRepository) UpdateQuiz(quiz *domain.Quiz) error {
	return r.db.Model(quiz).Updates(quiz).Error
}

func (r *quizRepository) DeleteQuiz(id uint) error {
	return r.db.Delete(&domain.Quiz{}, id).Error
}
