package handlers

import (
	"strings"

	"e-piket-be/config"
	"e-piket-be/internal/models"

	"github.com/gofiber/fiber/v3"
)

type JamKosongHandler struct{}

func NewJamKosongHandler() *JamKosongHandler {
	return &JamKosongHandler{}
}

type JamKosongInput struct {
	Kelas          string `json:"kelas"`
	JamKe          string `json:"jam_ke"`
	GuruTidakHadir string `json:"guru_tidak_hadir"`
	GuruPengganti  string `json:"guru_pengganti"`
}

// GetAll mengambil semua data guru pengganti / jam kosong
func (h *JamKosongHandler) GetAll(c fiber.Ctx) error {
	var list []models.JamKosong
	if err := config.DB.Order("id DESC").Find(&list).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data jam kosong: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    list,
	})
}

// GetByID mengambil detail 1 data berdasarkan ID
func (h *JamKosongHandler) GetByID(c fiber.Ctx) error {
	id := c.Params("id")
	var item models.JamKosong

	if err := config.DB.First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    item,
	})
}

// Create menambahkan data guru pengganti / jam kosong baru
func (h *JamKosongHandler) Create(c fiber.Ctx) error {
	var input JamKosongInput
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format JSON tidak valid: " + err.Error(),
		})
	}

	input.Kelas = strings.TrimSpace(input.Kelas)
	input.JamKe = strings.TrimSpace(input.JamKe)
	input.GuruTidakHadir = strings.TrimSpace(input.GuruTidakHadir)
	input.GuruPengganti = strings.TrimSpace(input.GuruPengganti)

	if input.Kelas == "" || input.JamKe == "" || input.GuruTidakHadir == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Kelas, jam pelajaran ke, dan guru tidak hadir wajib diisi",
		})
	}

	if input.GuruPengganti == "" {
		input.GuruPengganti = "Belum Ada"
	}

	item := models.JamKosong{
		Kelas:          input.Kelas,
		JamKe:          input.JamKe,
		GuruTidakHadir: input.GuruTidakHadir,
		GuruPengganti:  input.GuruPengganti,
	}

	if err := config.DB.Create(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menyimpan data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Data guru pengganti berhasil ditambahkan",
		"data":    item,
	})
}

// Update mengubah data guru pengganti / jam kosong berdasarkan ID
func (h *JamKosongHandler) Update(c fiber.Ctx) error {
	id := c.Params("id")
	var item models.JamKosong

	if err := config.DB.First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak ditemukan",
		})
	}

	var input JamKosongInput
	if err := c.Bind().JSON(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format JSON tidak valid: " + err.Error(),
		})
	}

	if strings.TrimSpace(input.Kelas) != "" {
		item.Kelas = strings.TrimSpace(input.Kelas)
	}
	if strings.TrimSpace(input.JamKe) != "" {
		item.JamKe = strings.TrimSpace(input.JamKe)
	}
	if strings.TrimSpace(input.GuruTidakHadir) != "" {
		item.GuruTidakHadir = strings.TrimSpace(input.GuruTidakHadir)
	}
	if strings.TrimSpace(input.GuruPengganti) != "" {
		item.GuruPengganti = strings.TrimSpace(input.GuruPengganti)
	}

	if err := config.DB.Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal memperbarui data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data berhasil diperbarui",
		"data":    item,
	})
}

// Delete menghapus data guru pengganti / jam kosong berdasarkan ID
func (h *JamKosongHandler) Delete(c fiber.Ctx) error {
	id := c.Params("id")
	var item models.JamKosong

	if err := config.DB.First(&item, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Data tidak ditemukan",
		})
	}

	if err := config.DB.Delete(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghapus data: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Data berhasil dihapus",
	})
}
