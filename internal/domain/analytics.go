package domain

// ExamAnalyticsDTO berisi ringkasan statistik dari suatu ujian
type ExamAnalyticsDTO struct {
	ExamID        uint    `json:"exam_id"`
	ExamTitle     string  `json:"exam_title"`
	TotalStudents int     `json:"total_students"`
	AverageScore  float64 `json:"average_score"`
	HighestScore  float64 `json:"highest_score"`
	LowestScore   float64 `json:"lowest_score"`
	PassRate      float64 `json:"pass_rate"` // Persentase kelulusan (0-100)
}

// ItemAnalysisDTO berisi evaluasi kualitas per butir soal
type ItemAnalysisDTO struct {
	QuestionID   uint   `json:"question_id"`
	QuestionText string `json:"question_text"`
	Type         string `json:"type"`
	CorrectCount int    `json:"correct_count"`
	WrongCount   int    `json:"wrong_count"`
	Unanswered   int    `json:"unanswered"`
	Difficulty   string `json:"difficulty"` // "Mudah", "Sedang", "Sulit"
}

// AnalyticsUsecase adalah kontrak untuk logika bisnis analitik
type AnalyticsUsecase interface {
	GetExamAnalytics(examID uint) (*ExamAnalyticsDTO, error)
	GetItemAnalysis(examID uint) ([]ItemAnalysisDTO, error)
}
