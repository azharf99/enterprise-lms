package domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Exam merepresentasikan Ujian Besar (CBT) di tingkat Mata Pelajaran
type Exam struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	CourseID    uint   `gorm:"not null" json:"course_id"` // Hierarkinya langsung ke Course
	Title       string `gorm:"type:varchar(200);not null" json:"title"`
	ExamType    string `gorm:"type:varchar(50);not null" json:"exam_type"` // "PTS", "PAS", "TryOut"
	Description string `gorm:"type:text" json:"description"`

	// Konfigurasi CBT
	TimeLimit    int    `gorm:"not null;default:0" json:"time_limit"`
	PassingScore int    `gorm:"not null;default:70" json:"passing_score"`
	CBTToken     string `gorm:"type:varchar(10)" json:"cbt_token"`
	IsRandomized bool   `gorm:"default:true" json:"is_randomized"`

	// Jadwal Pelaksanaan CBT (Siswa hanya bisa akses di rentang waktu ini)
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`

	Questions []ExamQuestion `gorm:"foreignKey:ExamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"questions,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
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

// --- KONTRAK REPOSITORY ---
type ExamRepository interface {
	Create(exam *Exam) error
	GetByCourseID(courseID uint) ([]Exam, error)
	GetByID(id uint) (Exam, error)
	Update(exam *Exam) error
	Delete(id uint) error
}

type ExamQuestionRepository interface {
	Create(question *ExamQuestion) error
	GetByExamID(examID uint) ([]ExamQuestion, error)
}

type ExamAttemptRepository interface {
	Create(attempt *ExamAttempt) error
	GetByID(id uint) (ExamAttempt, error)
	GetLatestAttempt(examID, userID uint) (ExamAttempt, error)
	Update(attempt *ExamAttempt) error
}

// --- KONTRAK USECASE ---
type ExamUsecase interface {
	CreateExam(courseID uint, title, examType, description string, timeLimit, passingScore int, startTime, endTime *time.Time) (*Exam, error)
	GenerateCBTToken(examID uint) (string, error)
	GenerateQuestionsWithAI(examID uint, topic, qType string, count int) ([]ExamQuestion, error)

	// CBT Execution
	StartAttempt(examID, userID uint, inputToken string) (*ExamAttempt, []ExamQuestion, error)
	SubmitAttempt(attemptID uint, answers datatypes.JSON) (*ExamAttempt, error)
}
