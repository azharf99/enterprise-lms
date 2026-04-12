package usecase

import (
	"bytes"
	cripto_rand "crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/azharf99/enterprise-lms/internal/domain"
	"github.com/azharf99/enterprise-lms/pkg/utils"
	"gorm.io/datatypes"
)

type examUsecase struct {
	examRepo         domain.ExamRepository
	examQuestionRepo domain.ExamQuestionRepository
	examAttemptRepo  domain.ExamAttemptRepository
}

func NewExamUsecase(er domain.ExamRepository, eqr domain.ExamQuestionRepository, ear domain.ExamAttemptRepository) domain.ExamUsecase {
	return &examUsecase{examRepo: er, examQuestionRepo: eqr, examAttemptRepo: ear}
}

// ... (CreateExam menggunakan logika CRUD standar) ...
func (u *examUsecase) CreateExam(courseID uint, req domain.CreateExamRequest) (*domain.Exam, error) {
	isRandom := true // Default jika kosong
	if req.IsRandomized != nil {
		isRandom = *req.IsRandomized
	}
	exam := &domain.Exam{
		CourseID:     courseID,
		Title:        req.Title,
		ExamType:     req.ExamType,
		Description:  req.Description,
		TimeLimit:    req.TimeLimit,
		PassingScore: req.PassingScore,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		CBTToken:     req.CBTToken,
		IsRandomized: isRandom,
		Status:       req.Status,
	}
	if err := u.examRepo.CreateExam(exam); err != nil {
		return nil, err
	}
	return exam, nil
}

func (u *examUsecase) GenerateCBTToken(examID uint) (string, error) {
	exam, err := u.examRepo.GetExamByID(examID)
	if err != nil {
		return "", errors.New("ujian tidak ditemukan")
	}

	bytes := make([]byte, 3) // Menghasilkan 6 karakter hex (misal: "a1b2c3")
	cripto_rand.Read(bytes)
	token := hex.EncodeToString(bytes)

	exam.CBTToken = token
	if err := u.examRepo.UpdateExam(&exam); err != nil {
		return "", err
	}
	return token, nil
}

func (u *examUsecase) GenerateExamQuestionsWithAI(examID uint, qType, topic string, count int) ([]domain.ExamQuestion, error) {
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

	var saved []domain.ExamQuestion
	for _, aiQ := range aiQuestions {
		optJSON, _ := json.Marshal(aiQ.Options)
		ansJSON, _ := json.Marshal(aiQ.CorrectAnswer)

		newQ := &domain.ExamQuestion{
			ExamID:        examID,
			Type:          domain.QuestionType(aiQ.Type),
			Text:          aiQ.Text,
			Options:       datatypes.JSON(optJSON),
			CorrectAnswer: datatypes.JSON(ansJSON),
			Points:        aiQ.Points,
			Explanation:   aiQ.Explanation,
		}
		if u.examQuestionRepo.CreateExamQuestion(newQ) == nil {
			saved = append(saved, *newQ)
		}
	}
	return saved, nil
}

func (u *examUsecase) StartExamAttempt(examID, userID uint, inputToken string) (*domain.ExamAttempt, []domain.ExamQuestion, error) {
	exam, err := u.examRepo.GetExamByID(examID)
	if err != nil {
		return nil, nil, errors.New("ujian tidak ditemukan")
	}

	// 1. Validasi Jadwal Ujian
	now := time.Now()
	if exam.StartTime != nil && now.Before(*exam.StartTime) {
		return nil, nil, errors.New("ujian belum dimulai")
	}
	if exam.EndTime != nil && now.After(*exam.EndTime) {
		return nil, nil, errors.New("waktu ujian telah berakhir")
	}

	// 2. Validasi Token CBT
	if exam.CBTToken != "" && inputToken != exam.CBTToken {
		return nil, nil, errors.New("token ujian tidak valid")
	}

	// 3. Cek apakah sudah pernah mengerjakan (CBT biasanya 1 kali percobaan)
	_, err = u.examAttemptRepo.GetLatestExamAttempt(examID, userID)
	if err == nil {
		return nil, nil, errors.New("anda sudah mengerjakan ujian ini")
	}

	attempt := &domain.ExamAttempt{
		ExamID:    examID,
		UserID:    userID,
		StartedAt: now,
	}
	if err := u.examAttemptRepo.CreateExamAttempt(attempt); err != nil {
		return nil, nil, err
	}

	questions := exam.Questions
	if exam.IsRandomized {
		rand.New(rand.NewSource(time.Now().UnixNano()))
		rand.Shuffle(len(questions), func(i, j int) { questions[i], questions[j] = questions[j], questions[i] })
	}

	return attempt, questions, nil
}

func (u *examUsecase) GetExamsByCourseID(courseID uint) ([]domain.Exam, error) {
	return u.examRepo.GetExamsByCourseID(courseID)
}

func (u *examUsecase) GetExamByID(id uint) (domain.Exam, error) {
	return u.examRepo.GetExamByID(id)
}

func (u *examUsecase) UpdateExam(id uint, req *domain.CreateExamRequest) (*domain.Exam, error) {
	exam, err := u.examRepo.GetExamByID(id)
	if err != nil {
		return nil, errors.New("ujian tidak ditemukan")
	}

	isRandom := true
	if req.IsRandomized != nil {
		isRandom = *req.IsRandomized
	}

	exam.Title = req.Title
	exam.ExamType = req.ExamType
	exam.Description = req.Description
	exam.TimeLimit = req.TimeLimit
	exam.PassingScore = req.PassingScore
	exam.StartTime = req.StartTime
	exam.EndTime = req.EndTime
	exam.CBTToken = req.CBTToken
	exam.IsRandomized = isRandom
	exam.Status = req.Status
	if err := u.examRepo.UpdateExam(&exam); err != nil {
		return nil, err
	}
	return &exam, nil
}

func (u *examUsecase) DeleteExam(id uint) error {
	return u.examRepo.DeleteExam(id)
}

func (u *examUsecase) SubmitExamAttempt(examAttemptID uint, userAnswers datatypes.JSON) (*domain.ExamAttempt, error) {
	// 1. Ambil data Attempt (Pastikan ada GetByID di repo Anda)
	// Catatan: Anda perlu menambahkan metode ini di examAttemptRepository Anda
	attempt, err := u.examAttemptRepo.GetExamAttemptByID(examAttemptID)
	if err != nil {
		return nil, errors.New("data pengerjaan ujian tidak ditemukan")
	}

	// Validasi jika sudah di-submit
	if attempt.CompletedAt != nil {
		return nil, errors.New("ujian ini sudah diselesaikan sebelumnya")
	}

	// 2. Ambil Soal-soal Ujian
	exam, err := u.examRepo.GetExamByID(attempt.ExamID)
	if err != nil {
		return nil, errors.New("data ujian tidak valid")
	}

	// 3. Parsing Jawaban Siswa
	var parsedAnswers map[string]interface{}
	if err := json.Unmarshal(userAnswers, &parsedAnswers); err != nil {
		return nil, errors.New("format jawaban tidak valid")
	}

	// 4. Proses Perhitungan Nilai (Grading Engine)
	totalMaxPoints := 0
	totalEarnedPoints := 0

	for _, question := range exam.Questions {
		totalMaxPoints += question.Points
		questionIDStr := fmt.Sprintf("%d", question.ID)

		userAns, exists := parsedAnswers[questionIDStr]
		if !exists {
			continue // Kosong/Tidak dijawab = 0 poin
		}

		userAnsBytes, _ := json.Marshal(userAns)

		switch question.Type {
		case domain.TypeMultipleAnswer:
			// Hitung dengan Partial Credit
			totalEarnedPoints += utils.CalculatePartialCredit(userAnsBytes, question.CorrectAnswer, question.Points)
		case domain.TypeEssay:
			// Esai diabaikan dari perhitungan otomatis (akan dinilai manual oleh Tutor nanti)
		default:
			// Hitung dengan Strict Match (Pilihan Ganda Biasa, Benar Salah, Isian)
			compactUser := new(bytes.Buffer)
			compactCorrect := new(bytes.Buffer)
			json.Compact(compactUser, userAnsBytes)
			json.Compact(compactCorrect, question.CorrectAnswer)
			if compactUser.String() == compactCorrect.String() {
				totalEarnedPoints += question.Points
			}
		}
	}

	// 5. Kalkulasi Skor Akhir (Skala 100)
	var finalScore float64 = 0
	if totalMaxPoints > 0 {
		finalScore = float64(totalEarnedPoints) / float64(totalMaxPoints) * 100
	}

	// 6. Simpan ke Database
	now := time.Now()
	attempt.CompletedAt = &now
	attempt.Score = finalScore
	attempt.Answers = userAnswers // Simpan JSON jawaban siswa untuk audit

	if err := u.examAttemptRepo.UpdateExamAttempt(&attempt); err != nil {
		return nil, errors.New("gagal menyimpan hasil ujian")
	}

	return &attempt, nil
}
