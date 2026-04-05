package utils

import (
	"encoding/json"
	"math"
)

// Pastikan Anda mengimpor "math" dan "encoding/json" di bagian atas file

// CalculatePartialCredit menghitung nilai untuk Pilihan Ganda Kompleks
func CalculatePartialCredit(userAnswerJSON, correctAnswerJSON []byte, maxPoints int) int {
	var userAnswers []string
	var correctAnswers []string

	// 1. Parsing JSON ke bentuk Array String
	// Asumsi format JSON: ["A", "C"]
	if err := json.Unmarshal(userAnswerJSON, &userAnswers); err != nil {
		return 0 // Jika format salah, nilainya 0
	}
	if err := json.Unmarshal(correctAnswerJSON, &correctAnswers); err != nil {
		return 0
	}

	totalCorrectOptions := len(correctAnswers)
	if totalCorrectOptions == 0 {
		return 0
	}

	// 2. Hitung bobot per opsi
	// Misal soal bobotnya 10, dan ada 2 jawaban benar (A dan C). Maka tiap opsi bernilai 5.
	pointPerItem := float64(maxPoints) / float64(totalCorrectOptions)
	earnedPoints := 0.0

	// Buat Map untuk pencarian jawaban benar dengan cepat (O(1))
	correctMap := make(map[string]bool)
	for _, ans := range correctAnswers {
		correctMap[ans] = true
	}

	// 3. Evaluasi jawaban siswa
	for _, ans := range userAnswers {
		if correctMap[ans] {
			// Jika siswa memilih opsi yang benar, tambah poin
			earnedPoints += pointPerItem
		} else {
			// Jika siswa memilih opsi yang salah, kurangi poin (Penalti)
			// Ini mencegah siswa menceklis semua opsi
			earnedPoints -= pointPerItem
		}
	}

	// 4. Skor tidak boleh minus
	if earnedPoints < 0 {
		earnedPoints = 0
	}

	// Membulatkan hasil (misal 7.5 menjadi 8)
	return int(math.Round(earnedPoints))
}
