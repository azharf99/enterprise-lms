package domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Exam merepresentasikan Ujian Besar (CBT) di tingkat Mata Pelajaran
type Exam struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CourseID     uint           `gorm:"not null" json:"course_id"` // Hierarkinya langsung ke Course
	Title        string         `gorm:"type:varchar(200);not null" json:"title"`
	ExamType     string         `gorm:"type:varchar(50);not null" json:"exam_type"` // "PTS", "PAS", "TryOut"
	Description  string         `gorm:"type:text" json:"description"`
	TimeLimit    int            `gorm:"not null;default:0" json:"time_limit"`
	PassingScore int            `gorm:"not null;default:70" json:"passing_score"`
	CBTToken     string         `gorm:"type:varchar(10)" json:"cbt_token"`
	IsRandomized bool           `gorm:"default:true" json:"is_randomized"`
	StartTime    *time.Time     `json:"start_time"`
	EndTime      *time.Time     `json:"end_time"`
	Status       string         `gorm:"type:varchar(50);default:'draft'" json:"status"` // Draft. Published
	Questions    []ExamQuestion `gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"questions,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// ExamQuestion merepresentasikan butir soal khusus ujian
type ExamQuestion struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	ExamID        uint           `gorm:"not null" json:"exam_id"`
	Type          QuestionType   `gorm:"type:varchar(50);not null" json:"type"`
	Text          string         `gorm:"type:text;not null" json:"text"`
	Options       datatypes.JSON `gorm:"type:jsonb" json:"options"`
	CorrectAnswer datatypes.JSON `gorm:"type:jsonb" json:"correct_answer"`
	Points        int            `gorm:"not null;default:10" json:"points"`
	Explanation   string         `gorm:"type:text" json:"explanation"`
}

// ExamAttempt mencatat pengerjaan CBT oleh siswa
type ExamAttempt struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ExamID      uint           `gorm:"not null" json:"exam_id"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	Score       float64        `gorm:"type:decimal(5,2);default:0" json:"score"`
	Answers     datatypes.JSON `gorm:"type:jsonb" json:"answers"`
	StartedAt   time.Time      `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at"`
}

// Format Request untuk Generate AI
type AIGenerateRequest struct {
	Topic string `json:"topic" binding:"required"`
	QType string `json:"q_type" binding:"required"`
	Count int    `json:"count" binding:"required,min=1,max=25"`
}

// --- Struct untuk Request Body ---
type CreateExamRequest struct {
	Title        string     `json:"title" binding:"required"`
	ExamType     string     `json:"exam_type" binding:"required"` // Misalnya: "PTS", "PAS"
	Description  string     `json:"description"`
	TimeLimit    int        `json:"time_limit"`
	PassingScore int        `json:"passing_score"`
	StartTime    *time.Time `json:"start_time"` // Format JSON harus RFC3339, misal: "2026-10-01T08:00:00Z"
	EndTime      *time.Time `json:"end_time"`
	CBTToken     string     `json:"cbt_token"`
	Status       string     `json:"status"`
	IsRandomized *bool      `json:"is_randomized"`
}

// --- KONTRAK REPOSITORY ---
type ExamRepository interface {
	CreateExam(exam *Exam) error
	GetExamsByCourseID(courseID uint) ([]Exam, error)
	GetExamByID(id uint) (Exam, error)
	UpdateExam(exam *Exam) error
	DeleteExam(id uint) error
}

type ExamQuestionRepository interface {
	CreateExamQuestion(question *ExamQuestion) error
	GetExamQuestionsByExamID(examID uint, isRandomized bool) ([]ExamQuestion, error)
	GetExamQuestionByID(id uint) (ExamQuestion, error)
	UpdateExamQuestion(exam *ExamQuestion) error
	DeleteExamQuestion(id uint) error
}

type ExamAttemptRepository interface {
	CreateExamAttempt(attempt *ExamAttempt) error
	GetExamAttemptByID(id uint) (ExamAttempt, error)
	GetExamAttemptsByExamID(examID uint) ([]ExamAttempt, error)
	GetLatestExamAttempt(examID, userID uint) (ExamAttempt, error)
	UpdateExamAttempt(attempt *ExamAttempt) error
}

// --- KONTRAK USECASE ---
type ExamUsecase interface {
	CreateExam(courseID uint, req CreateExamRequest) (*Exam, error)
	GenerateCBTToken(examID uint) (string, error)
	GenerateExamQuestionsWithAI(examID uint, topic, qType string, count int) ([]ExamQuestion, error)
	GetExamsByCourseID(courseID uint) ([]Exam, error)
	GetExamByID(id uint) (Exam, error)
	UpdateExam(id uint, req *CreateExamRequest) (*Exam, error)
	DeleteExam(id uint) error

	// CBT Execution
	StartExamAttempt(examID, userID uint, inputToken string) (*ExamAttempt, []ExamQuestion, error)
	SubmitExamAttempt(examAttemptID uint, answers datatypes.JSON) (*ExamAttempt, error)
}

type ExamQuestionUsecase interface {
	CreateExamQuestion(examID uint, qType QuestionType, text string, options, correctAnswer datatypes.JSON, points int, explanation string) (*ExamQuestion, error)
	GetExamQuestionsByExamID(examID uint, israndomized bool) ([]ExamQuestion, error)
	GetExamQuestionByID(id uint) (ExamQuestion, error)
	UpdateExamQuestion(id uint, qType QuestionType, text string, options, correctAnswer datatypes.JSON, points int, explanation string) (*ExamQuestion, error)
	DeleteExamQuestion(id uint) error
}
