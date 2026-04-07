package postgres

import (
	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/gorm"
)

// --- EXAM REPOSITORY ---
type examRepository struct{ db *gorm.DB }

func NewExamRepository(db *gorm.DB) domain.ExamRepository { return &examRepository{db: db} }

func (r *examRepository) CreateExam(exam *domain.Exam) error { return r.db.Create(exam).Error }
func (r *examRepository) GetExamsByCourseID(courseID uint) ([]domain.Exam, error) {
	var exams []domain.Exam
	err := r.db.Where("course_id = ?", courseID).Find(&exams).Error
	return exams, err
}
func (r *examRepository) GetExamByID(id uint) (domain.Exam, error) {
	var exam domain.Exam
	err := r.db.Preload("Questions").First(&exam, id).Error
	return exam, err
}
func (r *examRepository) UpdateExam(exam *domain.Exam) error {
	return r.db.Model(exam).Updates(exam).Error
}
func (r *examRepository) DeleteExam(id uint) error { return r.db.Delete(&domain.Exam{}, id).Error }

// --- EXAM QUESTION REPOSITORY ---
type examQuestionRepository struct{ db *gorm.DB }

func NewExamQuestionRepository(db *gorm.DB) domain.ExamQuestionRepository {
	return &examQuestionRepository{db: db}
}

func (r *examQuestionRepository) CreateExamQuestion(q *domain.ExamQuestion) error {
	return r.db.Create(q).Error
}
func (r *examQuestionRepository) GetExamQuestionsByExamID(examID uint) ([]domain.ExamQuestion, error) {
	var questions []domain.ExamQuestion
	err := r.db.Where("exam_id = ?", examID).Order("id asc").Find(&questions).Error
	return questions, err
}

func (r *examQuestionRepository) GetExamQuestionByID(id uint) (domain.ExamQuestion, error) {
	var question domain.ExamQuestion
	err := r.db.First(&question, id).Error
	return question, err
}

func (r *examQuestionRepository) UpdateExamQuestion(question *domain.ExamQuestion) error {
	return r.db.Model(question).Updates(question).Error
}

func (r *examQuestionRepository) DeleteExamQuestion(id uint) error {
	return r.db.Delete(&domain.ExamQuestion{}, id).Error
}

// --- EXAM ATTEMPT REPOSITORY ---
type examAttemptRepository struct{ db *gorm.DB }

func NewExamAttemptRepository(db *gorm.DB) domain.ExamAttemptRepository {
	return &examAttemptRepository{db: db}
}

func (r *examAttemptRepository) CreateExamAttempt(attempt *domain.ExamAttempt) error {
	return r.db.Create(attempt).Error
}

func (r *examAttemptRepository) GetExamAttemptByID(id uint) (domain.ExamAttempt, error) {
	var attempt domain.ExamAttempt
	err := r.db.First(&attempt, id).Error
	return attempt, err
}

func (r *examAttemptRepository) GetLatestExamAttempt(examID, userID uint) (domain.ExamAttempt, error) {
	var attempt domain.ExamAttempt
	err := r.db.Where("exam_id = ? AND user_id = ?", examID, userID).Order("started_at desc").First(&attempt).Error
	return attempt, err
}

func (r *examAttemptRepository) UpdateExamAttempt(attempt *domain.ExamAttempt) error {
	return r.db.Model(attempt).Updates(attempt).Error
}

// Tambahkan fungsi ini di bagian bawah file exam_repo.go
func (r *examAttemptRepository) GetExamAttemptsByExamID(examID uint) ([]domain.ExamAttempt, error) {
	var attempts []domain.ExamAttempt
	// Kita hanya mengambil attempt yang sudah disubmit (CompletedAt tidak null)
	err := r.db.Where("exam_id = ? AND completed_at IS NOT NULL", examID).Find(&attempts).Error
	return attempts, err
}
