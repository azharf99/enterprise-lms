package postgres

import (
	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/gorm"
)

type quizAttemptRepository struct {
	db *gorm.DB
}

func NewQuizAttemptRepository(db *gorm.DB) domain.QuizAttemptRepository {
	return &quizAttemptRepository{db: db}
}

func (r *quizAttemptRepository) CreateQuizAttempt(attempt *domain.QuizAttempt) error {
	return r.db.Create(attempt).Error
}

func (r *quizAttemptRepository) GetQuizAttemptByID(id uint) (domain.QuizAttempt, error) {
	var attempt domain.QuizAttempt
	err := r.db.First(&attempt, id).Error
	return attempt, err
}

func (r *quizAttemptRepository) GetLatestQuizAttempt(quizID, userID uint, status string) (domain.QuizAttempt, error) {
	var attempt domain.QuizAttempt
	err := r.db.Where("quiz_id = ? AND user_id = ? AND status = ?", quizID, userID, status).
		Order("attempt_number desc").
		First(&attempt).Error
	return attempt, err
}

func (r *quizAttemptRepository) GetQuizAttemptsByUser(quizID, userID uint) ([]domain.QuizAttempt, error) {
	var attempts []domain.QuizAttempt
	err := r.db.Where("quiz_id = ? AND user_id = ?", quizID, userID).Order("attempt_number asc").Find(&attempts).Error
	return attempts, err
}

func (r *quizAttemptRepository) CheckCompletedQuizAttempt(quizID, userID uint, status string) int64 {
	var completedCount int64
	r.db.Model(&domain.QuizAttempt{}).
		Where("quiz_id = ? AND user_id = ? AND status = ?", quizID, userID, "completed").
		Count(&completedCount)

	return completedCount
}

func (r *quizAttemptRepository) UpdateQuizAttempt(attempt *domain.QuizAttempt) error {
	return r.db.Model(attempt).Updates(attempt).Error
}
