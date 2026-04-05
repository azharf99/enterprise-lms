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

func (r *quizRepository) Create(quiz *domain.Quiz) error {
	return r.db.Create(quiz).Error
}

func (r *quizRepository) GetByModuleID(modulID uint) ([]domain.Quiz, error) {
	var quizzes []domain.Quiz
	err := r.db.Where("module_id = ?", modulID).Find(&quizzes).Error
	return quizzes, err
}

func (r *quizRepository) GetByID(id uint) (domain.Quiz, error) {
	var quiz domain.Quiz
	err := r.db.Preload("Questions").First(&quiz, id).Error
	return quiz, err
}

func (r *quizRepository) Update(quiz *domain.Quiz) error {
	return r.db.Model(quiz).Updates(quiz).Error
}

func (r *quizRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Quiz{}, id).Error
}
