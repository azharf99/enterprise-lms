package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/azharf99/enterprise-lms/pkg/utils"
	"gorm.io/datatypes"
)

type quizUsecase struct {
	quizRepo    domain.QuizRepository
	attemptRepo domain.QuizAttemptRepository
}

func NewQuizUsecase(qr domain.QuizRepository, ar domain.QuizAttemptRepository) domain.QuizUsecase {
	return &quizUsecase{
		quizRepo:    qr,
		attemptRepo: ar,
	}
}

func (u *quizUsecase) CreateQuiz(moduleID uint, title, description string, timeLimit, passingScore int, isRandomized bool, maxAttempts int) (*domain.Quiz, error) {
	quiz := &domain.Quiz{
		ModuleID:     moduleID,
		Title:        title,
		Description:  description,
		TimeLimit:    timeLimit,
		PassingScore: passingScore,
		IsRandomized: isRandomized,
		MaxAttempts:  maxAttempts,
	}
	if err := u.quizRepo.Create(quiz); err != nil {
		return nil, err
	}
	return quiz, nil
}

func (u *quizUsecase) GetQuizzesByModule(moduleID uint) ([]domain.Quiz, error) {
	return u.quizRepo.GetByModuleID(moduleID)
}

func (u *quizUsecase) GetQuizByID(id uint) (domain.Quiz, error) {
	return u.quizRepo.GetByID(id)
}

func (u *quizUsecase) UpdateQuiz(id uint, title, description string, timeLimit, passingScore int) (*domain.Quiz, error) {
	module, err := u.quizRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("quiz tidak ditemukan")
	}

	module.Title = title
	module.Description = description
	module.TimeLimit = timeLimit
	module.PassingScore = passingScore

	if err := u.quizRepo.Update(&module); err != nil {
		return nil, err
	}
	return &module, nil
}

func (u *quizUsecase) DeleteQuiz(id uint) error {
	return u.quizRepo.Delete(id)
}

func (u *quizUsecase) StartAttempt(quizID, userID uint) (*domain.QuizAttempt, []domain.Question, error) {
	quiz, err := u.quizRepo.GetByID(quizID)
	if err != nil {
		return nil, nil, errors.New("kuis tidak ditemukan")
	}

	attempts, _ := u.attemptRepo.GetAttemptsByUser(quizID, userID)
	if quiz.MaxAttempts > 0 && len(attempts) >= quiz.MaxAttempts {
		return nil, nil, errors.New("batas maksimal percobaan kuis telah tercapai")
	}

	attempt := &domain.QuizAttempt{
		QuizID:        quizID,
		UserID:        userID,
		AttemptNumber: len(attempts) + 1,
		StartedAt:     time.Now(),
	}

	if err := u.attemptRepo.Create(attempt); err != nil {
		return nil, nil, err
	}

	questions := quiz.Questions
	if quiz.IsRandomized {
		// Mengacak soal di level aplikasi (Golang) menggunakan Fisher-Yates shuffle
		rand.New(rand.NewSource(time.Now().UnixNano()))
		rand.Shuffle(len(questions), func(i, j int) {
			questions[i], questions[j] = questions[j], questions[i]
		})
	}

	return attempt, questions, nil
}

func (u *quizUsecase) SubmitAttempt(attemptID uint, userAnswers datatypes.JSON) (*domain.QuizAttempt, error) {
	// 1. Dapatkan data attempt
	// attempt, err := u.attemptRepo.GetLatestAttempt(0, 0) // Kita gunakan ID langsung jika ada method GetByID di repo attempt. Untuk kasus ini, mari asumsikan kita butuh GetByID.
	// Catatan: Anda harus menambahkan method GetByID(id uint) ke QuizAttemptRepository!
	// Karena di interface sebelumnya belum ada, mari kita asumsikan kita punya method itu.
	attempt, err := u.attemptRepo.GetByID(attemptID)

	// Untuk contoh ini, saya buat query langsung agar kode ini tidak gagal. Tapi disarankan ditambahkan di Repo.
	// Jika attempt sudah selesai, tolak.
	if attempt.CompletedAt != nil {
		return nil, errors.New("kuis ini sudah disubmit sebelumnya")
	}

	// 2. Dapatkan kuis dan soal-soalnya untuk dicocokkan
	quiz, err := u.quizRepo.GetByID(attempt.QuizID)
	if err != nil {
		return nil, err
	}

	// 3. Parsing jawaban siswa. Asumsi format JSON: {"1": "A", "2": ["A", "C"], "3": "DNA"}
	// dimana key adalah ID Question dalam bentuk string
	var parsedAnswers map[string]interface{}
	if err := json.Unmarshal(userAnswers, &parsedAnswers); err != nil {
		return nil, errors.New("format jawaban tidak valid")
	}

	// 4. Proses Kalkulasi Skor
	totalMaxPoints := 0
	totalEarnedPoints := 0

	for _, question := range quiz.Questions {
		totalMaxPoints += question.Points
		questionIDStr := fmt.Sprintf("%d", question.ID)

		userAnswer, exists := parsedAnswers[questionIDStr]
		if !exists {
			continue // Tidak dijawab = 0 poin
		}

		userAnswerBytes, _ := json.Marshal(userAnswer)

		switch question.Type {
		case domain.TypeMultipleAnswer:
			// Gunakan logika Partial Credit untuk Pilihan Ganda Kompleks
			earned := utils.CalculatePartialCredit(userAnswerBytes, question.CorrectAnswer, question.Points)
			totalEarnedPoints += earned
		case domain.TypeEssay:
			// Essay dinilai manual oleh Tutor
			// Poin ditambahkan nanti pada saat proses grading manual
		default:
			// Tipe soal lainnya (Pilihan Ganda biasa, Benar/Salah, Isian Singkat)
			// Menggunakan pencocokan string eksak (Strict Equality)
			compactUserAnswer := new(bytes.Buffer)
			compactCorrectAnswer := new(bytes.Buffer)
			json.Compact(compactUserAnswer, userAnswerBytes)
			json.Compact(compactCorrectAnswer, question.CorrectAnswer)

			if compactUserAnswer.String() == compactCorrectAnswer.String() {
				totalEarnedPoints += question.Points
			}
		}
	}

	// 5. Hitung persentase akhir (Skala 0 - 100)
	var finalScore float64 = 0
	if totalMaxPoints > 0 {
		finalScore = float64(totalEarnedPoints) / float64(totalMaxPoints) * 100
	}

	// 6. Simpan hasil ke database
	now := time.Now()
	attempt.CompletedAt = &now
	attempt.Score = finalScore
	attempt.Answers = userAnswers

	if err := u.attemptRepo.Update(&attempt); err != nil {
		return nil, errors.New("gagal menyimpan hasil kuis")
	}

	return &attempt, nil
}
