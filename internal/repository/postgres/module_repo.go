package postgres

import (
	"github.com/azharf99/enterprise-lms/internal/domain"
	"gorm.io/gorm"
)

type moduleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) domain.ModuleRepository {
	return &moduleRepository{db: db}
}

func (r *moduleRepository) Create(module *domain.Module) error {
	return r.db.Create(module).Error
}

func (r *moduleRepository) GetByCourseID(courseID uint) ([]domain.Module, error) {
	var modules []domain.Module
	// Mengambil semua modul milik sebuah Course dan mengurutkannya berdasarkan sequence
	err := r.db.Where("course_id = ?", courseID).Order("sequence asc").Find(&modules).Error
	return modules, err
}

func (r *moduleRepository) GetByID(id uint) (domain.Module, error) {
	var module domain.Module
	// Preload Lessons agar saat melihat detail modul, daftar materinya juga ikut terbawa
	err := r.db.Preload("Lessons").First(&module, id).Error
	return module, err
}

func (r *moduleRepository) Update(module *domain.Module) error {
	return r.db.Model(module).Updates(module).Error
}

func (r *moduleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Module{}, id).Error
}
