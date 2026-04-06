package usecase

import (
	"errors"

	"github.com/azharf99/enterprise-lms/internal/domain"
)

type enrollmentUsecase struct {
	enrollRepo domain.EnrollmentRepository
}

func NewEnrollmentUsecase(cr domain.EnrollmentRepository) domain.EnrollmentUsecase {
	return &enrollmentUsecase{enrollRepo: cr}
}

func (u *enrollmentUsecase) EnrollStudent(courseID, userID uint) error {
	// Bisa tambahkan validasi: cek apakah course & user eksis terlebih dahulu
	isEnrolled, err := u.enrollRepo.CheckEnrollment(courseID, userID)
	if err != nil {
		return errors.New("Tidak dapat mengecek enrollment")
	}

	if isEnrolled {
		return nil
	} else {
		return u.enrollRepo.Enroll(courseID, userID)
	}

}

func (u *enrollmentUsecase) UnenrollStudent(courseID, userID uint) error {
	isEnrolled, err := u.enrollRepo.CheckEnrollment(courseID, userID)
	if err != nil {
		return errors.New("Tidak dapat mengecek enrollment")
	}

	if isEnrolled {
		return u.enrollRepo.Unenroll(courseID, userID)
	} else {
		return nil
	}
}

func (u *enrollmentUsecase) GetEnrolledStudents(courseID uint) ([]domain.Enrollment, error) {
	return u.enrollRepo.GetEnrolledUsers(courseID)
}
