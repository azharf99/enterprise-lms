package usecase

import (
	"errors"

	"github.com/azharf99/enterprise-lms/internal/domain"
)

type lessonUsecase struct {
	lessonRepo domain.LessonRepository
}

func NewLessonUsecase(lr domain.LessonRepository) domain.LessonUsecase {
	return &lessonUsecase{lessonRepo: lr}
}

func (u *lessonUsecase) CreateLesson(moduleID uint, title string, lessonType domain.LessonType, content string, sequence int) (*domain.Lesson, error) {
	if title == "" || content == "" {
		return nil, errors.New("judul dan konten materi tidak boleh kosong")
	}

	lesson := &domain.Lesson{
		ModuleID: moduleID,
		Title:    title,
		Type:     lessonType,
		Content:  content,
		Sequence: sequence,
	}

	if err := u.lessonRepo.Create(lesson); err != nil {
		return nil, err
	}
	return lesson, nil
}

func (u *lessonUsecase) GetLessonsByModule(moduleID uint) ([]domain.Lesson, error) {
	return u.lessonRepo.GetByModuleID(moduleID)
}

func (u *lessonUsecase) GetLessonByID(id uint) (domain.Lesson, error) {
	return u.lessonRepo.GetByID(id)
}

func (u *lessonUsecase) UpdateLesson(id uint, title string, lessonType domain.LessonType, content string, sequence int) (*domain.Lesson, error) {
	lesson, err := u.lessonRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("materi tidak ditemukan")
	}

	lesson.Title = title
	lesson.Type = lessonType
	lesson.Content = content
	lesson.Sequence = sequence

	if err := u.lessonRepo.Update(&lesson); err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (u *lessonUsecase) DeleteLesson(id uint) error {
	return u.lessonRepo.Delete(id)
}
