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
	exam, err := u.examRepo.GetExamByID(examID)
	if err != nil {
		return nil, errors.New("ujian tidak ditemukan")
	}

	attempts, err := u.examAttemptRepo.GetExamAttemptsByExamID(examID)
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
	exam, err := u.examRepo.GetExamByID(examID)
	if err != nil {
		return nil, err
	}

	attempts, err := u.examAttemptRepo.GetExamAttemptsByExamID(examID)
	if err != nil {
		return nil, err
	}

	totalAttempts := len(attempts)
	var analysisResult []domain.ItemAnalysisDTO

	// Jika belum ada yang mengerjakan, kembalikan array kosong
	if totalAttempts == 0 {
		return analysisResult, nil
	}

	// Optimization: Pre-parse and pre-compact all attempts' answers once
	// We map each attempt to a map of [questionID]compactedAnswerString
	type processedAttempt struct {
		answers map[string]string
	}
	processedAttempts := make([]processedAttempt, totalAttempts)
	for i, att := range attempts {
		var rawAnswers map[string]json.RawMessage
		_ = json.Unmarshal(att.Answers, &rawAnswers)

		processed := processedAttempt{
			answers: make(map[string]string),
		}
		for qID, rawAns := range rawAnswers {
			compacted := new(bytes.Buffer)
			if err := json.Compact(compacted, rawAns); err == nil {
				processed.answers[qID] = compacted.String()
			}
		}
		processedAttempts[i] = processed
	}

	// Evaluasi setiap soal
	for _, question := range exam.Questions {
		correctCount := 0
		wrongCount := 0
		unansweredCount := 0
		qIDStr := fmt.Sprintf("%d", question.ID)

		// Optimization: Pre-compact correct answer for this question
		compactCorrect := new(bytes.Buffer)
		_ = json.Compact(compactCorrect, question.CorrectAnswer)
		correctStr := compactCorrect.String()

		// Cek jawaban seluruh siswa untuk soal ini
		for _, att := range processedAttempts {
			userAnsCompacted, exists := att.answers[qIDStr]
			if !exists {
				unansweredCount++
				continue
			}

			if userAnsCompacted == correctStr {
				correctCount++
			} else {
				wrongCount++
			}
		}

		// Hitung Indeks Kesukaran (P) = Jumlah Benar / Total Siswa
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
