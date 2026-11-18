package handlers

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

type KehadiranGuruHandler struct{}

func NewKehadiranGuruHandler() *KehadiranGuruHandler {
	return &KehadiranGuruHandler{}
}

type DummyTeacherAttendance struct {
	ID       int    `json:"id"`
	NamaGuru string `json:"nama_guru"`
	Status   string `json:"status"` // "Hadir" atau "Tidak Hadir"
	Tanggal  string `json:"tanggal"`
}

// GetAll mengembalikan data dummy kehadiran guru hari ini (Read-Only, menunggu API Presensi Guru selesai)
func (h *KehadiranGuruHandler) GetAll(c fiber.Ctx) error {
	today := time.Now().Format("2006-01-02")

	dummyTeachers := []DummyTeacherAttendance{
		{ID: 1, NamaGuru: "Agus Septiano, S.Pd.", Status: "Hadir", Tanggal: today},
		{ID: 2, NamaGuru: "Albertus Kevin, S.Kom.", Status: "Hadir", Tanggal: today},
		{ID: 3, NamaGuru: "Alverttus Calvin Chong", Status: "Hadir", Tanggal: today},
		{ID: 4, NamaGuru: "Budi Santoso, M.Pd.", Status: "Hadir", Tanggal: today},
		{ID: 5, NamaGuru: "Siti Aminah, S.Kom.", Status: "Tidak Hadir", Tanggal: today},
		{ID: 6, NamaGuru: "Dedi Kurniawan, S.Pd.", Status: "Tidak Hadir", Tanggal: today},
		{ID: 7, NamaGuru: "Dimas Maulana Ishaq, S.Kom.", Status: "Hadir", Tanggal: today},
		{ID: 8, NamaGuru: "Mieke Rahmawaty, S.Kom.", Status: "Hadir", Tanggal: today},
		{ID: 9, NamaGuru: "Asit Rahmah, S.Pd.", Status: "Hadir", Tanggal: today},
		{ID: 10, NamaGuru: "Hendra Wijaya, S.Kom.", Status: "Hadir", Tanggal: today},
		{ID: 11, NamaGuru: "Winchet Asep, S.Pd.", Status: "Hadir", Tanggal: today},
		{ID: 12, NamaGuru: "Rina Marlina, S.Pd.", Status: "Tidak Hadir", Tanggal: today},
		{ID: 13, NamaGuru: "Fajar Ramadhan, M.Kom.", Status: "Hadir", Tanggal: today},
		{ID: 14, NamaGuru: "Eka Lestari, S.Pd.", Status: "Hadir", Tanggal: today},
		{ID: 15, NamaGuru: "Ahmad Fauzi, S.T.", Status: "Hadir", Tanggal: today},
		{ID: 16, NamaGuru: "Nurul Hidayati, S.Kom.", Status: "Tidak Hadir", Tanggal: today},
		{ID: 17, NamaGuru: "Rizky Pratama, S.Kom.", Status: "Hadir", Tanggal: today},
		{ID: 18, NamaGuru: "Dewi Anggraini, S.Pd.", Status: "Hadir", Tanggal: today},
		{ID: 19, NamaGuru: "Iwan Setiawan, S.Pd.", Status: "Hadir", Tanggal: today},
		{ID: 20, NamaGuru: "Gheraldy Moses (WS-RPL)", Status: "Hadir", Tanggal: today},
	}

	search := strings.ToLower(strings.TrimSpace(c.Query("search")))
	statusFilter := strings.ToLower(strings.TrimSpace(c.Query("status")))

	var filtered []DummyTeacherAttendance
	totalHadir := 0
	totalTidakHadir := 0

	for _, t := range dummyTeachers {
		if t.Status == "Hadir" {
			totalHadir++
		} else {
			totalTidakHadir++
		}

		if search != "" && !strings.Contains(strings.ToLower(t.NamaGuru), search) {
			continue
		}

		if statusFilter != "" && statusFilter != "semua" && strings.ToLower(t.Status) != statusFilter {
			continue
		}

		filtered = append(filtered, t)
	}

	if filtered == nil {
		filtered = []DummyTeacherAttendance{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"source":  "dummy (menunggu API Presensi Guru)",
		"tanggal": today,
		"summary": fiber.Map{
			"total_guru":  len(dummyTeachers),
			"hadir":       totalHadir,
			"tidak_hadir": totalTidakHadir,
		},
		"total_filtered": len(filtered),
		"data":           filtered,
	})
}
