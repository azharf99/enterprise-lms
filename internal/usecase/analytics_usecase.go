package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/azharf99/enterprise-lms/internal/domain"
)

type analyticsUsecase struct {
	examRepo        domain.ExamRepository
	examAttemptRepo domain.ExamAttemptRepository
}

func NewAnalyticsUsecase(er domain.ExamRepository, ear domain.ExamAttemptRepository) domain.AnalyticsUsecase {
	return &analyticsUsecase{
		examRepo:        er,
		examAttemptRepo: ear,
	}
}

func (u *analyticsUsecase) GetExamAnalytics(examID uint) (*domain.ExamAnalyticsDTO, error) {
	exam, err := u.examRepo.GetByID(examID)
	if err != nil {
		return nil, errors.New("ujian tidak ditemukan")
	}

	attempts, err := u.examAttemptRepo.GetByExamID(examID)
	if err != nil {
		return nil, errors.New("gagal mengambil data pengerjaan ujian")
	}

	totalStudents := len(attempts)
	if totalStudents == 0 {
		return &domain.ExamAnalyticsDTO{
			ExamID:    examID,
			ExamTitle: exam.Title,
		}, nil
	}

	var totalScore, highest, lowest float64
	passedCount := 0
	lowest = attempts[0].Score

	for _, att := range attempts {
		totalScore += att.Score
		if att.Score > highest {
			highest = att.Score
		}
		if att.Score < lowest {
			lowest = att.Score
		}
		if att.Score >= float64(exam.PassingScore) {
			passedCount++
		}
	}

	return &domain.ExamAnalyticsDTO{
		ExamID:        examID,
		ExamTitle:     exam.Title,
		TotalStudents: totalStudents,
		AverageScore:  totalScore / float64(totalStudents),
		HighestScore:  highest,
		LowestScore:   lowest,
		PassRate:      (float64(passedCount) / float64(totalStudents)) * 100,
	}, nil
}

func (u *analyticsUsecase) GetItemAnalysis(examID uint) ([]domain.ItemAnalysisDTO, error) {
	exam, err := u.examRepo.GetByID(examID)
	if err != nil {
		return nil, err
	}

	attempts, err := u.examAttemptRepo.GetByExamID(examID)
	if err != nil {
		return nil, err
	}

	totalAttempts := len(attempts)
	var analysisResult []domain.ItemAnalysisDTO

	// Jika belum ada yang mengerjakan, kembalikan array kosong
	if totalAttempts == 0 {
		return analysisResult, nil
	}

	// Evaluasi setiap soal
	for _, question := range exam.Questions {
		correctCount := 0
		wrongCount := 0
		unansweredCount := 0
		qIDStr := fmt.Sprintf("%d", question.ID)

		// Cek jawaban seluruh siswa untuk soal ini
		for _, att := range attempts {
			var parsedAnswers map[string]interface{}
			_ = json.Unmarshal(att.Answers, &parsedAnswers)

			userAns, exists := parsedAnswers[qIDStr]
			if !exists {
				unansweredCount++
				continue
			}

			// Menggunakan komparasi strict untuk analisis tingkat kesukaran
			userAnsBytes, _ := json.Marshal(userAns)
			compactUser := new(bytes.Buffer)
			compactCorrect := new(bytes.Buffer)
			json.Compact(compactUser, userAnsBytes)
			json.Compact(compactCorrect, question.CorrectAnswer)

			if compactUser.String() == compactCorrect.String() {
				correctCount++
			} else {
				wrongCount++
			}
		}

		// Hitung Indeks Kesukaran (P) = Jumlah Benar / Total Siswa
		// Kriteria standar evaluasi pendidikan:
		// P > 0.70        : Mudah
		// 0.30 <= P <= 0.70 : Sedang
		// P < 0.30        : Sulit
		difficultyIndex := float64(correctCount) / float64(totalAttempts)
		difficultyLabel := "Sedang"
		if difficultyIndex > 0.70 {
			difficultyLabel = "Mudah"
		} else if difficultyIndex < 0.30 {
			difficultyLabel = "Sulit"
		}

		analysisResult = append(analysisResult, domain.ItemAnalysisDTO{
			QuestionID:   question.ID,
			QuestionText: question.Text,
			Type:         string(question.Type),
			CorrectCount: correctCount,
			WrongCount:   wrongCount,
			Unanswered:   unansweredCount,
			Difficulty:   difficultyLabel,
		})
	}

	return analysisResult, nil
}
