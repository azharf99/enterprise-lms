package postgres

import (
	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) domain.QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(question *domain.Question) error {
	return r.db.Create(question).Error
}

func (r *questionRepository) GetByQuizID(quizID uint, isRandomized bool) ([]domain.Question, error) {
	var questions []domain.Question
	db := r.db.Where("quiz_id = ?", quizID)

	if isRandomized {
		// PostgreSQL: Menggunakan RANDOM() untuk mengacak urutan baris
		db = db.Order("RANDOM()")
	} else {
		db = db.Order("id asc")
	}

	err := db.Find(&questions).Error
	return questions, err
}

func (r *questionRepository) GetByID(id uint) (domain.Question, error) {
	var question domain.Question
	err := r.db.First(&question, id).Error
	return question, err
}

func (r *questionRepository) Update(question *domain.Question) error {
	return r.db.Model(question).Updates(question).Error
}

func (r *questionRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Question{}, id).Error
}
