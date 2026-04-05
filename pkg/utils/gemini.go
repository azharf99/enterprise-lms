package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GenerateQuizJSON menghubungi Google Gemini AI untuk membuat soal format JSON murni
func GenerateQuizJSON(topic, qType string, count int) (string, error) {
	ctx := context.Background()

	// Ambil API Key dari Environment Variable (.env)
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY tidak ditemukan di environment")
	}

	// Inisialisasi client AI
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", fmt.Errorf("gagal membuat AI client: %v", err)
	}
	defer client.Close()

	// Pilih model (Flash sangat direkomendasikan untuk teks karena cepat dan murah)
	model := client.GenerativeModel("gemini-1.5-flash")

	// Fitur krusial: Memaksa model merespons HANYA dalam format JSON
	model.ResponseMIMEType = "application/json"

	// Merangkai instruksi (Prompt Engineering)
	prompt := fmt.Sprintf(`Buatkan %d soal %s untuk ujian tingkat menengah/atas dengan topik: "%s".

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
7. Struktur atau format untuk setiap objek WAJIB seperti ini:
{
  "type": "MultipleChoice",
  "text": "Teks pertanyaan secara lengkap",
  "options": ["A", "B", "C", "D", "E"],
  "correct_answer": "C",
  "points": 10,
  "explanation": "Penjelasan detail mengapa jawaban tersebut benar"
}`, count, qType, topic)

	// Mengirim permintaan ke server Gemini
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("gagal melakukan request ke AI: %v", err)
	}

	// Mengekstrak teks dari respons AI
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		part := resp.Candidates[0].Content.Parts[0]
		if textPart, ok := part.(genai.Text); ok {
			return string(textPart), nil
		}
	}

	return "", fmt.Errorf("respons dari AI kosong atau tidak sesuai format")
}
