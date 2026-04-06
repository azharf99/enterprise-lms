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
