package domain

import (
	"time"

	"gorm.io/gorm"
)

type Course struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"type:varchar(200);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	// Relasi Many-to-Many dengan User (Tutor)
	Tutors []User `gorm:"many2many:course_tutors;" json:"tutors,omitempty"`
	// Relasi One-to-Many dengan Module
	Modules   []Module       `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"modules,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Module merepresentasikan bab di dalam suatu mata pelajaran
type Module struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CourseID  uint           `gorm:"not null" json:"course_id"`
	Title     string         `gorm:"type:varchar(200);not null" json:"title"`
	Sequence  int            `gorm:"not null;default:1" json:"sequence"`
	Lessons   []Lesson       `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"lessons,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type LessonType string

const (
	LessonTypeText  LessonType = "Text"
	LessonTypeVideo LessonType = "Video"
	LessonTypeAudio LessonType = "Audio"
	LessonTypePDF   LessonType = "PDF"
)

// Lesson merepresentasikan materi spesifik
type Lesson struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	ModuleID  uint           `gorm:"not null" json:"module_id"`
	Title     string         `gorm:"type:varchar(200);not null" json:"title"`
	Type      LessonType     `gorm:"type:varchar(50);not null;default:'Text'" json:"type"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Sequence  int            `gorm:"not null;default:1" json:"sequence"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// KONTRAK REPOSITORY DAN USECASE

// KONTRAK REPOSITORY
type CourseRepository interface {
	Create(course *Course) error
	GetAll() ([]Course, error)
	GetByID(id uint) (Course, error)
	Update(course *Course) error
	Delete(id uint) error
	AddTutors(courseID uint, tutorIDs []uint) error
}

// KONTRAK USECASE
type CourseUsecase interface {
	CreateCourse(title, description string, tutorIDs []uint) (*Course, error)
	GetAllCourses() ([]Course, error)
	GetCourseByID(id uint) (Course, error)
	UpdateCourse(id uint, title, description string, tutorIDs []uint) (*Course, error)
	DeleteCourse(id uint) error
}
