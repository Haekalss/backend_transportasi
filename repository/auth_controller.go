package repository

import (
	"context"
	"os"
	"regexp"
	"time"
	"transport-app/config"
	"transport-app/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserCollection() *mongo.Collection {
	return config.GetCollection("users")
}

// Fungsi untuk validasi format email
func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}

// Register User godoc
// @Summary Register a new user
// @Description Mendaftarkan pengguna baru dengan username dan email
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "User credentials"
// @Success 201 {object} models.SuccessResponse "User created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 409 {object} models.ErrorResponse "Username or email already exists"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/register [post]
func Register(c *fiber.Ctx) error {

	var input struct {
		Username             string `json:"username"`
		Email                string `json:"email"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Validasi field wajib diisi
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username, email, dan password wajib diisi"})
	}

	// Validasi format email
	if !isEmailValid(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format email tidak valid"})
	}

	// Validasi panjang password
	if len(input.Password) < 8 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password minimal harus 8 karakter"})
	}

	// Validasi konfirmasi password
	if input.Password != input.PasswordConfirmation {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Konfirmasi password tidak cocok"})
	}

	// Cek apakah username atau email sudah ada
	// Menggunakan $or untuk memeriksa keduanya dalam satu query
	filter := bson.M{
		"$or": []bson.M{
			{"username": input.Username},
			{"email": input.Email},
		},
	}
	count, err := getUserCollection().CountDocuments(context.TODO(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saat memeriksa data pengguna"})
	}
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username atau email sudah terdaftar"})
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengenkripsi password"})
	}

	newUser := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	_, err = getUserCollection().InsertOne(context.TODO(), newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat pengguna baru"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Pengguna berhasil dibuat"})
}

// Login User godoc
// @Summary Login a user
// @Description Login menggunakan username dan password untuk mendapatkan token
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body AuthRequest true "User credentials"
// @Success 200 {object} LoginResponse "Login successful, token returned"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Invalid username or password"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/login [post]
func Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Cari user di database berdasarkan username
	var user models.User
	err := getUserCollection().FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		// Jika tidak ditemukan, berikan pesan error yang generik
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username atau password salah"})
	}

	// Cek password
	if !checkPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Username atau password salah"})
	}

	// Buat token JWT
	claims := jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
		"user_id":  user.ID.Hex(),
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat token autentikasi"})
	}

	return c.JSON(fiber.Map{
		"token": t,
		"role":  user.Role,
	})
}
