package usecase

import (
	"errors"

	"github.com/azharf99/enterprise-lms/internal/domain"
)

type moduleUsecase struct {
	moduleRepo domain.ModuleRepository
}

func NewModuleUsecase(mr domain.ModuleRepository) domain.ModuleUsecase {
	return &moduleUsecase{moduleRepo: mr}
}

func (u *moduleUsecase) CreateModule(courseID uint, title string, sequence int) (*domain.Module, error) {
	if title == "" {
		return nil, errors.New("judul modul tidak boleh kosong")
	}

	module := &domain.Module{
		CourseID: courseID,
		Title:    title,
		Sequence: sequence,
	}

	if err := u.moduleRepo.Create(module); err != nil {
		return nil, err
	}
	return module, nil
}

func (u *moduleUsecase) GetModulesByCourse(courseID uint) ([]domain.Module, error) {
	return u.moduleRepo.GetByCourseID(courseID)
}

func (u *moduleUsecase) GetModuleByID(id uint) (domain.Module, error) {
	return u.moduleRepo.GetByID(id)
}

func (u *moduleUsecase) UpdateModule(id uint, title string, sequence int) (*domain.Module, error) {
	module, err := u.moduleRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("modul tidak ditemukan")
	}

	module.Title = title
	module.Sequence = sequence

	if err := u.moduleRepo.Update(&module); err != nil {
		return nil, err
	}
	return &module, nil
}

func (u *moduleUsecase) DeleteModule(id uint) error {
	return u.moduleRepo.Delete(id)
}
