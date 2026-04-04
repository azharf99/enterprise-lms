package postgres

import (
	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/gorm"
)

type courseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) domain.CourseRepository {
	return &courseRepository{db: db}
}

func (r *courseRepository) Create(course *domain.Course) error {
	return r.db.Create(course).Error
}

func (r *courseRepository) GetAll() ([]domain.Course, error) {
	var courses []domain.Course
	err := r.db.Preload("Tutors").Find(&courses).Error
	return courses, err
}

func (r *courseRepository) GetByID(id uint) (domain.Course, error) {
	var course domain.Course
	err := r.db.Preload("Tutors").
		Preload("Modules").
		Preload("Modules.Lessons").
		First(&course, id).Error
	return course, err
}

// Update memperbarui data dasar mata pelajaran
func (r *courseRepository) Update(course *domain.Course) error {
	// GORM Updates akan memperbarui field yang tidak kosong
	return r.db.Model(course).Updates(course).Error
}

// Delete menghapus mata pelajaran (Soft Delete karena ada gorm.DeletedAt)
func (r *courseRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Course{}, id).Error
}

// AddTutors menimpa/memperbarui daftar tutor untuk course tertentu
func (r *courseRepository) AddTutors(courseID uint, tutorIDs []uint) error {
	var course domain.Course
	if err := r.db.First(&course, courseID).Error; err != nil {
		return err
	}

	var tutors []domain.User
	if err := r.db.Where("id IN ?", tutorIDs).Find(&tutors).Error; err != nil {
		return err
	}

	// Gunakan Replace agar jika ada update, daftar tutor yang lama diganti dengan yang baru
	return r.db.Model(&course).Association("Tutors").Replace(&tutors)
}
