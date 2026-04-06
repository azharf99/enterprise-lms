package domain

import (
	"time"

	"github.com/azharf99/enterprise-lms/pkg/utils"
	"gorm.io/gorm"
)

type Role string

const (
	RoleSiswa  Role = "Siswa"
	RoleTutor  Role = "Tutor"
	RoleAdmin  Role = "Admin"
	RoleEditor Role = "Editor"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Role      Role           `gorm:"type:varchar(20);not null;default:'Siswa'" json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserRepository adalah kontrak untuk berinteraksi dengan database
type UserRepository interface {
	BulkInsertUsers(users []User) error
	GetUserByEmail(email string) (User, error)
	GetAllUsers() ([]User, error)
	GetUserByID(id uint) (User, error)
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id uint) error
}

// UserUsecase adalah kontrak untuk logika bisnis
type UserUsecase interface {
	ImportFromCSV(records [][]string) (int, error)
	GetAllUsers() ([]User, error)
	CreateUser(name, email, password string, role Role) (*User, error)
	UpdateUser(id uint, name, email, password string, role Role) (*User, error)
	DeleteUser(id uint) error
	Login(email, password string) (*utils.TokenPair, error)
	RefreshAccessToken(refreshToken string) (*utils.TokenPair, error) // Fungsi baru
}

// LoginRequest adalah format data yang diharapkan saat user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"omitempty"`
	Role     Role   `json:"role" binding:"required,oneof=Siswa Tutor Admin Editor"`
}

type UserCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Role     Role   `json:"role" binding:"required,oneof=Siswa Tutor Admin Editor"`
}

// RefreshRequest adalah format data saat frontend meminta access token baru
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
