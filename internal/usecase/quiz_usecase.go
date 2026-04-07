package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/azharf99/enterprise-lms/pkg/utils"
	"gorm.io/datatypes"
)

type quizUsecase struct {
	quizRepo     domain.QuizRepository
	attemptRepo  domain.QuizAttemptRepository
	questionRepo domain.QuizQuestionRepository
}

func NewQuizUsecase(qr domain.QuizRepository, ar domain.QuizAttemptRepository, qtr domain.QuizQuestionRepository) domain.QuizUsecase {
	return &quizUsecase{
		quizRepo:     qr,
		attemptRepo:  ar,
		questionRepo: qtr,
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
	if err := u.quizRepo.CreateQuiz(quiz); err != nil {
		return nil, err
	}
	return quiz, nil
}

func (u *quizUsecase) GenerateQuizQuestionsWithAI(quizID uint, qType, topic string, count int) ([]domain.Question, error) {
	jsonResponse, err := utils.GenerateQuizJSON(topic, qType, count) // Memakai fungsi AI di pkg/utils/gemini.go
	if err != nil {
		return nil, errors.New("gagal generate soal: " + err.Error())
	}

	type AIQuestion struct {
		Type          string      `json:"type"`
		Text          string      `json:"text"`
		Options       interface{} `json:"options"`
		CorrectAnswer interface{} `json:"correct_answer"`
		Points        int         `json:"points"`
		Explanation   string      `json:"explanation"`
	}

	var aiQuestions []AIQuestion
	if err := json.Unmarshal([]byte(jsonResponse), &aiQuestions); err != nil {
		return nil, errors.New("format AI tidak valid")
	}

	var saved []domain.Question
	for _, aiQ := range aiQuestions {
		optJSON, _ := json.Marshal(aiQ.Options)
		ansJSON, _ := json.Marshal(aiQ.CorrectAnswer)

		newQ := &domain.Question{
			QuizID:        quizID,
			Type:          domain.QuestionType(aiQ.Type),
			Text:          aiQ.Text,
			Options:       datatypes.JSON(optJSON),
			CorrectAnswer: datatypes.JSON(ansJSON),
			Points:        aiQ.Points,
			Explanation:   aiQ.Explanation,
		}
		if u.questionRepo.CreateQuizQuestion(newQ) == nil {
			saved = append(saved, *newQ)
		}
	}
	return saved, nil
}

func (u *quizUsecase) GetQuizzesByModule(moduleID uint) ([]domain.Quiz, error) {
	return u.quizRepo.GetQuizzesByModuleID(moduleID)
}

func (u *quizUsecase) GetQuizByID(id uint) (domain.Quiz, error) {
	return u.quizRepo.GetQuizByID(id)
}

func (u *quizUsecase) UpdateQuiz(id uint, title, description string, timeLimit, passingScore int) (*domain.Quiz, error) {
	module, err := u.quizRepo.GetQuizByID(id)
	if err != nil {
		return nil, errors.New("quiz tidak ditemukan")
	}

	module.Title = title
	module.Description = description
	module.TimeLimit = timeLimit
	module.PassingScore = passingScore

	if err := u.quizRepo.UpdateQuiz(&module); err != nil {
		return nil, err
	}
	return &module, nil
}

func (u *quizUsecase) DeleteQuiz(id uint) error {
	return u.quizRepo.DeleteQuiz(id)
}

func (u *quizUsecase) StartAttempt(quizID, userID uint, status string) (*domain.AttemptResponse, error) {
	// Ambil data kuis untuk mengecek batas max_attempt (misal default 1)
	quiz, _ := u.quizRepo.GetQuizByID(quizID)
	
	// 1. CEK ATTEMPT YANG SEDANG BERJALAN (RESUME)
	activeAttempt, err := u.attemptRepo.GetLatestQuizAttempt(quizID, userID, status)
	if err == nil && activeAttempt.ID != 0 {
		// Ditemukan Attempt yang masih berjalan!
		// Langsung ambilkan daftar soal dan kembalikan (tanpa menambah kuota attempt)
		questions, _ := u.questionRepo.GetQuizQuestionsByQuizID(uint(quizID), quiz.IsRandomized)

		// Jangan lupa filter kunci jawaban agar tidak bocor ke siswa!
		filteredQuestions := filterAnswersOut(questions)

		return &domain.AttemptResponse{
			Attempt:   activeAttempt,
			Questions: filteredQuestions,
		}, nil
	}

	// 2. JIKA TIDAK ADA YANG AKTIF, CEK KUOTA MAKSIMAL ATTEMPT
	completedCount := u.attemptRepo.CheckCompletedQuizAttempt(quizID, userID, "completed")


	// Jika di tabel quiz tidak ada kolom max_attempts, asumsikan batasnya 1
	// Anda bisa menyesuaikan ini jika ada kolom MaxAttempts di database Anda
	if completedCount >= int64(quiz.MaxAttempts) {
		return nil, errors.New("batas maksimal percobaan kuis telah tercapai")
	}

	// 3. BUAT ATTEMPT BARU
	newAttempt := domain.QuizAttempt{
		QuizID: quizID,
		UserID: userID,
		Status: "in_progress",
		StartedAt: time.Now(),
	}

	if err := u.attemptRepo.CreateQuizAttempt(&newAttempt); err != nil {
		return nil, errors.New("gagal memulai kuis baru")
	}

	// Ambil soal untuk attempt baru
	questions, _ := u.questionRepo.GetQuizQuestionsByQuizID(uint(quizID), quiz.IsRandomized)
	filteredQuestions := filterAnswersOut(questions)

	return &domain.AttemptResponse{
		Attempt:   newAttempt,
		Questions: filteredQuestions,
	}, nil
}

func filterAnswersOut(questions []domain.Question) []domain.QuestionAttemptDTO {
	var filtered []domain.QuestionAttemptDTO
	for _, q := range questions {
		filtered = append(filtered, domain.QuestionAttemptDTO{
			ID:      q.ID,
			Type:    q.Type,
			Text:    q.Text,
			Options: q.Options,
			Points:  q.Points,
		})
	}
	return filtered
}

func (u *quizUsecase) SubmitAttempt(attemptID uint, userAnswers datatypes.JSON) (*domain.QuizAttempt, error) {
	// 1. Dapatkan data attempt
	// attempt, err := u.attemptRepo.GetLatestAttempt(0, 0) // Kita gunakan ID langsung jika ada method GetByID di repo attempt. Untuk kasus ini, mari asumsikan kita butuh GetByID.
	// Catatan: Anda harus menambahkan method GetByID(id uint) ke QuizAttemptRepository!
	// Karena di interface sebelumnya belum ada, mari kita asumsikan kita punya method itu.
	attempt, err := u.attemptRepo.GetQuizAttemptByID(attemptID)

	// Untuk contoh ini, saya buat query langsung agar kode ini tidak gagal. Tapi disarankan ditambahkan di Repo.
	// Jika attempt sudah selesai, tolak.
	if attempt.CompletedAt != nil {
		return nil, errors.New("kuis ini sudah disubmit sebelumnya")
	}

	// 2. Dapatkan kuis dan soal-soalnya untuk dicocokkan
	quiz, err := u.quizRepo.GetQuizByID(attempt.QuizID)
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
	attempt.AttemptNumber += 1
	attempt.Passed = finalScore >= float64(quiz.PassingScore)

	if err := u.attemptRepo.UpdateQuizAttempt(&attempt); err != nil {
		return nil, errors.New("gagal menyimpan hasil kuis")
	}

	return &attempt, nil
}
