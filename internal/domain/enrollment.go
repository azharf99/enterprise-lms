package domain

import "time"

// Enrollment merepresentasikan kepesertaan siswa dalam suatu kursus
type Enrollment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CourseID  uint      `gorm:"not null;uniqueIndex:idx_course_user" json:"course_id"`
	UserID    uint      `gorm:"not null;uniqueIndex:idx_course_user" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relasi opsional untuk mempermudah preload
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

type EnrollmentRepository interface {
	Enroll(courseID, userID uint) error
	Unenroll(courseID, userID uint) error
	CheckEnrollment(courseID, userID uint) (bool, error)
	GetEnrolledUsers(courseID uint) ([]Enrollment, error)
	CheckQuizAccess(quizID, userID uint) (bool, error)
	CheckModuleAccess(moduleID, userID uint) (bool, error)
	CheckExamAccess(examID, userID uint) (bool, error)
	CheckLessonAccess(lessonID, userID uint) (bool, error)
	CheckQuestionAccess(questionID, userID uint) (bool, error)
}

type EnrollmentUsecase interface {
	EnrollStudent(courseID, userID uint) error
	UnenrollStudent(courseID, userID uint) error
	GetEnrolledStudents(courseID uint) ([]Enrollment, error)
}
