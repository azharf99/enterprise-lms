package usecase

import (
	"errors"
	"strings"
	"sync"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/azharf99/enterprise-lms/pkg/utils"
	"gorm.io/gorm"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

// NewUserUsecase adalah constructor untuk usecase
func NewUserUsecase(ur domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: ur}
}

func (u *userUsecase) ImportFromCSV(records [][]string) (int, error) {
	if len(records) < 2 {
		return 0, errors.New("data CSV kosong atau hanya berisi header")
	}

	var wg sync.WaitGroup
	userChan := make(chan domain.User, len(records)-1)
	errorChan := make(chan string, len(records)-1)

	for i, record := range records {
		if i == 0 || len(record) != 4 {
			continue // Lewati header dan baris tidak valid
		}

		wg.Add(1)
		go func(row []string) {
			defer wg.Done()

			name := strings.TrimSpace(row[0])
			email := strings.TrimSpace(row[1])
			rawPassword := strings.TrimSpace(row[2])
			role := domain.Role(strings.TrimSpace(row[3]))

			hashedPassword, err := utils.HashPassword(rawPassword)
			if err != nil {
				errorChan <- "error hashing"
				return
			}

			userChan <- domain.User{
				Name:     name,
				Email:    email,
				Password: hashedPassword,
				Role:     role,
			}
		}(record)
	}

	wg.Wait()
	close(userChan)
	close(errorChan)

	var users []domain.User
	for user := range userChan {
		users = append(users, user)
	}

	if len(users) == 0 {
		return 0, errors.New("tidak ada data valid untuk disimpan")
	}

	// Panggil repository untuk menyimpan data
	err := u.userRepo.BulkInsertUsers(users)
	if err != nil {
		return 0, err
	}

	return len(users), nil
}

func (u *userUsecase) CreateUser(name, email, password string, role domain.Role) (*domain.User, error) {
	if name == "" {
		return nil, errors.New("nama tidak boleh kosong")
	}

	if email == "" {
		return nil, errors.New("email tidak boleh kosong")
	}

	if password == "" {
		return nil, errors.New("password tidak boleh kosong")
	}

	user := &domain.User{
		Name:  name,
		Email: email,
		Role:  role,
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	if err := u.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) GetAllUsers() ([]domain.User, error) {
	return u.userRepo.GetAllUsers()
}

func (u *userUsecase) Login(email, password string) (*utils.TokenPair, error) {
	// 1. Cari user di database
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("email atau password salah")
		}
		return nil, err
	}

	// 2. Bandingkan password
	match := utils.CheckPasswordHash(password, user.Password)
	if !match {
		return nil, errors.New("email atau password salah")
	}

	// 3. Buat JWT
	return utils.GenerateTokenPair(user.ID, string(user.Role))
}

func (u *userUsecase) RefreshAccessToken(refreshToken string) (*utils.TokenPair, error) {
	// 1. Validasi token dan ambil ID User
	userID, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// 2. Cek apakah user masih ada di database dan belum di-banned/dihapus
	user, err := u.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan atau tidak aktif")
	}

	// 3. Buat pasangan token yang baru
	return utils.GenerateTokenPair(user.ID, string(user.Role))
}

func (u *userUsecase) UpdateUser(id uint, name, email string, role domain.Role) (*domain.User, error) {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return nil, errors.New("Pengguna tidak ditemukan")
	}

	user.Name = name
	user.Email = email
	user.Role = role

	// Perbarui data dasar
	if err := u.userRepo.UpdateUser(&user); err != nil {
		return nil, err
	}

	// Ambil data terbaru untuk dikembalikan
	updatedUser, _ := u.userRepo.GetUserByID(id)
	return &updatedUser, nil
}

func (u *userUsecase) DeleteUser(id uint) error {
	// Pastikan user ada sebelum dihapus
	_, err := u.userRepo.GetUserByID(id)
	if err != nil {
		return errors.New("Pengguna tidak ditemukan")
	}
	return u.userRepo.DeleteUser(id)
}
