package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GenerateQuizJSON menghubungi Google Gemini AI menggunakan REST API standar
func GenerateQuizJSON(topic, qType string, count int) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY tidak ditemukan di environment")
	}

	// Endpoint resmi Gemini 1.5 Flash
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s", apiKey)

	prompt := fmt.Sprintf(`Buatkan %d soal dengan tipe %s untuk ujian tingkat menengah atas (SMA) dengan topik: "%s".

Aturan format:
1. Kembalikan murni array JSON, tanpa backtick, tanpa tulisan 'json' di awal.
2. Jika aku hanya meminta satu jenis soal, maka seluruh soal yang kamu buat harus satu jenis.
3. Aturan untuk jenis soal: 
   - MultipleChoice: Ini adalah Pilihan Ganda. Harus memiliki "options" (array string) dengan 5 opsi, misalnya ["A", "B", "C", "D", "E"] dan "correct_answer" (string yang merupakan salah satu opsi).
   - MultipleAnswer: Ini adalah Pilihan Ganda Kompleks. Harus memiliki "options" (array string) dengan 5 opsi, misalnya ["A", "B", "C", "D", "E"] dan "correct_answer" (array string yang merupakan subset dari opsi), misalnya ["A", "C"].
   - TrueFalse: Ini adalah Benar/Salah. "options" harus ["True", "False"] dan "correct_answer" harus salah satu dari itu. 
   - Matching: Ini adalah Menjodohkan. "options" harus berupa array objek dengan "left" dan "right", dan "correct_answer" harus berupa array pasangan yang benar. 
   - ShortAnswer: Ini adalah isian singkat. Tidak perlu "options", tetapi "correct_answer" harus berupa string atau array string yang benar. 
   - Essay: Ini adalah essai/uraian. Tidak perlu "options", tetapi "correct_answer" harus berupa string yang merupakan kunci jawaban ideal. 
4. "explanation" harus berisi penjelasan detail mengapa jawaban tersebut benar, dan ini akan ditampilkan kepada siswa setelah mereka menyelesaikan kuis. 
5. "points" harus berupa integer yang menunjukkan bobot nilai untuk setiap soal, dengan nilai default 10 jika tidak ditentukan. 
6. Pastikan semua teks dalam bahasa Indonesia yang baik dan benar, dan gunakan istilah yang sesuai dengan konteks pendidikan di Indonesia. 
7. Sehubungan respons kamu akan ditembak langsung ke API, maka JANGAN PERNAH sertakan teks penjelasan atau instruksi tambahan di luar format JSON. Hanya kembalikan array JSON yang valid. 
8. Struktur atau format untuk setiap objek WAJIB seperti ini:
{
  "type": "MultipleChoice",
  "text": "Teks pertanyaan secara lengkap",
  "options": ["A", "B", "C", "D", "E"],
  "correct_answer": "C",
  "points": 10,
  "explanation": "Penjelasan detail mengapa jawaban tersebut benar"
}`, count, qType, topic)

	// Menyusun struktur payload sesuai spesifikasi REST API Gemini
	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"responseMimeType": "application/json",
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("gagal menyusun payload JSON: %v", err)
	}

	// Membuat HTTP Request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("gagal membuat HTTP request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Mengeksekusi Request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal menghubungi server Gemini: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gagal membaca respons: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error dari Gemini API (Status %d): %s", resp.StatusCode, string(body))
	}

	// Parsing struktur respons JSON Gemini
	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", fmt.Errorf("gagal memparsing JSON balasan AI: %v", err)
	}

	// Mengekstrak teks balasan
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("respons dari AI kosong atau tidak sesuai format")
}
