package postgres

import (
	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/gorm"
)

type lessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) domain.LessonRepository {
	return &lessonRepository{db: db}
}

func (r *lessonRepository) Create(lesson *domain.Lesson) error {
	return r.db.Create(lesson).Error
}

func (r *lessonRepository) GetByModuleID(moduleID uint) ([]domain.Lesson, error) {
	var lessons []domain.Lesson
	err := r.db.Where("module_id = ?", moduleID).Order("sequence asc").Find(&lessons).Error
	return lessons, err
}

func (r *lessonRepository) GetByID(id uint) (domain.Lesson, error) {
	var lesson domain.Lesson
	err := r.db.First(&lesson, id).Error
	return lesson, err
}

func (r *lessonRepository) Update(lesson *domain.Lesson) error {
	return r.db.Model(lesson).Updates(lesson).Error
}

func (r *lessonRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Lesson{}, id).Error
}
