package usecase

import (
	"errors"

	"github.com/azharf99/enterprise-lms/internal/domain"
)

type courseUsecase struct {
	courseRepo domain.CourseRepository
}

func NewCourseUsecase(cr domain.CourseRepository) domain.CourseUsecase {
	return &courseUsecase{courseRepo: cr}
}

// ... (CreateCourse, GetAllCourses tetap sama) ...
func (u *courseUsecase) CreateCourse(title, description string, tutorIDs []uint) (*domain.Course, error) {
	if title == "" {
		return nil, errors.New("judul mata pelajaran tidak boleh kosong")
	}

	course := &domain.Course{
		Title:       title,
		Description: description,
	}

	if err := u.courseRepo.Create(course); err != nil {
		return nil, err
	}

	if len(tutorIDs) > 0 {
		_ = u.courseRepo.AddTutors(course.ID, tutorIDs)
	}

	return course, nil
}

func (u *courseUsecase) GetAllCourses() ([]domain.Course, error) {
	return u.courseRepo.GetAll()
}

func (u *courseUsecase) GetCourseByID(id uint) (domain.Course, error) {
	return u.courseRepo.GetByID(id)
}

func (u *courseUsecase) UpdateCourse(id uint, title, description string, tutorIDs []uint) (*domain.Course, error) {
	// Pastikan course ada
	course, err := u.courseRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("mata pelajaran tidak ditemukan")
	}

	course.Title = title
	course.Description = description

	// Perbarui data dasar
	if err := u.courseRepo.Update(&course); err != nil {
		return nil, err
	}

	// Perbarui relasi tutor
	if len(tutorIDs) > 0 {
		if err := u.courseRepo.AddTutors(id, tutorIDs); err != nil {
			return nil, errors.New("gagal memperbarui daftar tutor")
		}
	} else {
		// Jika array kosong dikirim, kita hapus semua relasi tutornya
		_ = u.courseRepo.AddTutors(id, []uint{})
	}

	// Ambil data terbaru untuk dikembalikan
	updatedCourse, _ := u.courseRepo.GetByID(id)
	return &updatedCourse, nil
}

func (u *courseUsecase) DeleteCourse(id uint) error {
	// Pastikan course ada sebelum dihapus
	_, err := u.courseRepo.GetByID(id)
	if err != nil {
		return errors.New("mata pelajaran tidak ditemukan")
	}
	return u.courseRepo.Delete(id)
}