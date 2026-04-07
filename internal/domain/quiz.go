package domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Quiz merepresentasikan satu set ujian/kuis dalam sebuah modul
type Quiz struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ModuleID     uint           `gorm:"not null" json:"module_id"`
	Title        string         `gorm:"type:varchar(200);not null" json:"title"`
	Description  string         `gorm:"type:text" json:"description"`
	TimeLimit    int            `gorm:"not null;default:0" json:"time_limit"`     // Waktu pengerjaan dalam menit (0 = tanpa batas)
	PassingScore int            `gorm:"not null;default:70" json:"passing_score"` // Nilai KKM
	IsRandomized bool           `gorm:"default:false" json:"is_randomized"`
	MaxAttempts  int            `gorm:"default:1" json:"max_attempts"` // 0 = Unlimited
	Questions    []Question     `gorm:"foreignKey:QuizID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"questions,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// QuestionType membatasi jenis pertanyaan yang didukung sistem
type QuestionType string

const (
	TypeMultipleChoice QuestionType = "MultipleChoice" // 5 Opsi (A-E)
	TypeMultipleAnswer QuestionType = "MultipleAnswer" // Pilihan Ganda Kompleks (Partial Credit)
	TypeTrueFalse      QuestionType = "TrueFalse"      // Benar / Salah
	TypeMatching       QuestionType = "Matching"       // Menjodohkan
	TypeShortAnswer    QuestionType = "ShortAnswer"    // Isian Singkat
	TypeEssay          QuestionType = "Essay"          // Uraian
)

// Question merepresentasikan butir soal
type Question struct {
	ID     uint         `gorm:"primaryKey" json:"id"`
	QuizID uint         `gorm:"not null" json:"quiz_id"`
	Type   QuestionType `gorm:"type:varchar(50);not null" json:"type"`
	Text   string       `gorm:"type:text;not null" json:"text"` // Teks soal

	// Gunakan datatypes.JSON agar GORM otomatis mengubah struct Golang menjadi JSONB di PostgreSQL
	Options       datatypes.JSON `gorm:"type:jsonb" json:"options"`
	CorrectAnswer datatypes.JSON `gorm:"type:jsonb" json:"correct_answer"`

	Points      int    `gorm:"not null;default:10" json:"points"` // Bobot nilai per soal
	Explanation string `gorm:"type:text" json:"explanation"`      // Pembahasan setelah kuis selesai

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// QuizAttempt mencatat setiap kali siswa mengerjakan kuis
type QuizAttempt struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	QuizID        uint    `gorm:"not null" json:"quiz_id"`
	UserID        uint    `gorm:"not null" json:"user_id"`
	AttemptNumber int     `gorm:"not null" json:"attempt_number"`
	Score         float64 `gorm:"type:decimal(5,2);default:0" json:"score"`

	// Menyimpan jawaban siswa dalam bentuk JSON untuk audit/review
	// Format: [{"question_id": 1, "answer": "A"}, ...]
	Answers datatypes.JSON `gorm:"type:jsonb" json:"answers"`

	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at"` // Nilai null jika belum selesai
}

type QuizRepository interface {
	CreateQuiz(quiz *Quiz) error
	GetQuizzesByModuleID(moduleID uint) ([]Quiz, error)
	GetQuizByID(id uint) (Quiz, error)
	UpdateQuiz(quiz *Quiz) error
	DeleteQuiz(id uint) error
}

type QuizQuestionRepository interface {
	CreateQuizQuestion(question *Question) error
	GetQuizQuestionsByQuizID(quizID uint, isRandomized bool) ([]Question, error)
	GetQuizQuestionByID(id uint) (Question, error)
	UpdateQuizQuestion(question *Question) error
	DeleteQuizQuestion(id uint) error
}

type QuizAttemptRepository interface {
	CreateQuizAttempt(attempt *QuizAttempt) error
	GetQuizAttemptByID(id uint) (QuizAttempt, error)
	GetLatestQuizAttempt(quizID, userID uint) (QuizAttempt, error)
	GetQuizAttemptsByUser(quizID, userID uint) ([]QuizAttempt, error)
	UpdateQuizAttempt(attempt *QuizAttempt) error
}

type QuizUsecase interface {
	CreateQuiz(moduleID uint, title, description string, timeLimit, passingScore int, isRandomized bool, maxAttempts int) (*Quiz, error)
	GetQuizzesByModule(moduleID uint) ([]Quiz, error)
	GetQuizByID(id uint) (Quiz, error)
	UpdateQuiz(id uint, title, description string, timeLimit, passingScore int) (*Quiz, error)
	DeleteQuiz(id uint) error
	StartAttempt(quizID, userID uint) (*QuizAttempt, []Question, error)
	SubmitAttempt(attemptID uint, answers datatypes.JSON) (*QuizAttempt, error)
}

type QuizQuestionUsecase interface {
	CreateQuestion(quizID uint, qType QuestionType, text string, options, correctAnswer datatypes.JSON, points int, explanation string) (*Question, error)
	GetQuestionsByQuizID(quizID uint, israndomized bool) ([]Question, error)
	GetQuestionByID(id uint) (Question, error)
	UpdateQuestion(id uint, qType QuestionType, text string, options, correctAnswer datatypes.JSON, points int, explanation string) (*Question, error)
	DeleteQuestion(id uint) error
}
