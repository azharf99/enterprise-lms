package postgres

import (
	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/gorm"
)

type enrollmentRepositoy struct{ db *gorm.DB }

func NewEnrollmentRepository(db *gorm.DB) domain.EnrollmentRepository {
	return &enrollmentRepositoy{db}
}

func (r *enrollmentRepositoy) Enroll(courseID, userID uint) error {
	return r.db.Create(&domain.Enrollment{CourseID: courseID, UserID: userID}).Error
}

func (r *enrollmentRepositoy) Unenroll(courseID, userID uint) error {
	return r.db.Where("course_id = ? AND user_id = ?", courseID, userID).Delete(&domain.Enrollment{}).Error
}

func (r *enrollmentRepositoy) CheckEnrollment(courseID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&domain.Enrollment{}).Where("course_id = ? AND user_id = ?", courseID, userID).Count(&count).Error
	return count > 0, err
}

func (r *enrollmentRepositoy) GetEnrolledUsers(courseID uint) ([]domain.Enrollment, error) {
	var enrollments []domain.Enrollment
	err := r.db.Preload("User").Where("course_id = ?", courseID).Find(&enrollments).Error
	return enrollments, err
}

// CheckQuizAccess mengecek akses hanya dengan 1 kali ke database
func (r *enrollmentRepositoy) CheckQuizAccess(quizID, userID uint) (bool, error) {
	var count int64
	// SQL Setara:
	// SELECT count(*) FROM enrollments e
	// JOIN modules m ON e.course_id = m.course_id
	// JOIN quizzes q ON m.id = q.module_id
	// WHERE q.id = ? AND e.user_id = ?
	err := r.db.Table("enrollments").
		Joins("JOIN modules ON enrollments.course_id = modules.course_id").
		Joins("JOIN quizzes ON modules.id = quizzes.module_id").
		Where("quizzes.id = ? AND enrollments.user_id = ?", quizID, userID).
		Count(&count).Error

	return count > 0, err
}
// CheckQuestionAccess mengecek akses hanya dengan 1 kali ke database
func (r *enrollmentRepositoy) CheckQuestionAccess(questionID, userID uint) (bool, error) {
	var count int64
	// SQL Setara:
	// SELECT count(*) FROM enrollments e
	// JOIN modules m ON e.course_id = m.course_id
	// JOIN quizzes q ON m.id = q.module_id
	// WHERE q.id = ? AND e.user_id = ?
	err := r.db.Table("enrollments").
		Joins("JOIN modules ON enrollments.course_id = modules.course_id").
		Joins("JOIN quizzes ON modules.id = quizzes.module_id").
		Joins("JOIN questions ON quizzes.id = questions.quiz_id").
		Where("questions.id = ? AND enrollments.user_id = ?", questionID, userID).
		Count(&count).Error

	return count > 0, err
}

// CheckQuizAccess mengecek akses hanya dengan 1 kali ke database
func (r *enrollmentRepositoy) CheckLessonAccess(lessonID, userID uint) (bool, error) {
	var count int64
	// SQL Setara:
	// SELECT count(*) FROM enrollments e
	// JOIN modules m ON e.course_id = m.course_id
	// JOIN lessons q ON m.id = q.module_id
	// WHERE q.id = ? AND e.user_id = ?
	err := r.db.Table("enrollments").
		Joins("JOIN modules ON enrollments.course_id = modules.course_id").
		Joins("JOIN lessons ON modules.id = lessons.module_id").
		Where("lessons.id = ? AND enrollments.user_id = ?", lessonID, userID).
		Count(&count).Error

	return count > 0, err
}

// CheckModuleAccess mengecek akses modul langsung ke tabel enrollments
func (r *enrollmentRepositoy) CheckModuleAccess(moduleID, userID uint) (bool, error) {
	var count int64
	// SQL Setara: SELECT count(*) FROM enrollments e JOIN modules m ON e.course_id = m.course_id WHERE m.id = ? AND e.user_id = ?
	err := r.db.Table("enrollments").
		Joins("JOIN modules ON enrollments.course_id = modules.course_id").
		Where("modules.id = ? AND enrollments.user_id = ?", moduleID, userID).
		Count(&count).Error

	return count > 0, err
}

// CheckExamAccess mengecek akses ujian langsung ke tabel enrollments
func (r *enrollmentRepositoy) CheckExamAccess(examID, userID uint) (bool, error) {
	var count int64
	// SQL Setara: SELECT count(*) FROM enrollments e JOIN exams ex ON e.course_id = ex.course_id WHERE ex.id = ? AND e.user_id = ?
	err := r.db.Table("enrollments").
		Joins("JOIN exams ON enrollments.course_id = exams.course_id").
		Where("exams.id = ? AND enrollments.user_id = ?", examID, userID).
		Count(&count).Error

	return count > 0, err
}
