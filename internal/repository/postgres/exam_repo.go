package postgres

import (
	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/gorm"
)

// --- EXAM REPOSITORY ---
type examRepository struct{ db *gorm.DB }

func NewExamRepository(db *gorm.DB) domain.ExamRepository { return &examRepository{db: db} }

func (r *examRepository) Create(exam *domain.Exam) error { return r.db.Create(exam).Error }
func (r *examRepository) GetByCourseID(courseID uint) ([]domain.Exam, error) {
	var exams []domain.Exam
	err := r.db.Where("course_id = ?", courseID).Find(&exams).Error
	return exams, err
}
func (r *examRepository) GetByID(id uint) (domain.Exam, error) {
	var exam domain.Exam
	err := r.db.Preload("Questions").First(&exam, id).Error
	return exam, err
}
func (r *examRepository) Update(exam *domain.Exam) error { return r.db.Model(exam).Updates(exam).Error }
func (r *examRepository) Delete(id uint) error           { return r.db.Delete(&domain.Exam{}, id).Error }

// --- EXAM QUESTION REPOSITORY ---
type examQuestionRepository struct{ db *gorm.DB }

func NewExamQuestionRepository(db *gorm.DB) domain.ExamQuestionRepository {
	return &examQuestionRepository{db: db}
}

func (r *examQuestionRepository) Create(q *domain.ExamQuestion) error { return r.db.Create(q).Error }
func (r *examQuestionRepository) GetByExamID(examID uint) ([]domain.ExamQuestion, error) {
	var questions []domain.ExamQuestion
	err := r.db.Where("exam_id = ?", examID).Order("id asc").Find(&questions).Error
	return questions, err
}

// --- EXAM ATTEMPT REPOSITORY ---
type examAttemptRepository struct{ db *gorm.DB }

func NewExamAttemptRepository(db *gorm.DB) domain.ExamAttemptRepository {
	return &examAttemptRepository{db: db}
}

func (r *examAttemptRepository) Create(attempt *domain.ExamAttempt) error {
	return r.db.Create(attempt).Error
}

func (r *examAttemptRepository) GetByID(id uint) (domain.ExamAttempt, error) {
	var attempt domain.ExamAttempt
	err := r.db.First(&attempt, id).Error
	return attempt, err
}

func (r *examAttemptRepository) GetLatestAttempt(examID, userID uint) (domain.ExamAttempt, error) {
	var attempt domain.ExamAttempt
	err := r.db.Where("exam_id = ? AND user_id = ?", examID, userID).Order("started_at desc").First(&attempt).Error
	return attempt, err
}

func (r *examAttemptRepository) Update(attempt *domain.ExamAttempt) error {
	return r.db.Model(attempt).Updates(attempt).Error
}
