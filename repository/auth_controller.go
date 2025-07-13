package repository

import (
	"context"
	"os"
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

// Register User godoc
// @Summary Register a new user
// @Description Mendaftarkan pengguna baru
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body AuthRequest true "User credentials"
// @Success 201 {object} models.SuccessResponse "User created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 409 {object} models.ErrorResponse "Username already exists"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /api/register [post]
func Register(c *fiber.Ctx) error {

	var input struct {
		Username             string `json:"username"`
		Password             string `json:"password"`
		PasswordConfirmation string `json:"password_confirmation"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	
	// Validasi konfirmasi password
	if input.Password != input.PasswordConfirmation {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password confirmation does not match"})
	}

	// validasi input
	if input.Username == "" || input.Password == "" || input.PasswordConfirmation == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username and password are required"})
	}

	// Cek apakah username sudah ada
	count, err := getUserCollection().CountDocuments(context.TODO(), bson.M{"username": input.Username})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error checking username"})
	}
	if count > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Username already exists"})
	}

	hashedPassword, err := hashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	newUser := models.User{
		Username: input.Username,
		Password: hashedPassword,
		Role:     "user",
	}

	_, err = getUserCollection().InsertOne(context.TODO(), newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User created successfully"})
}

// Login User godoc
// @Summary Login a user
// @Description Login untuk mendapatkan token autentikasi
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

	// Cari user di database
	var user models.User
	err := getUserCollection().FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	// Cek password
	if !checkPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	}

	// Buat token JWT
	claims := jwt.MapClaims{
		"username": user.Username,
		"user_id":  user.ID.Hex(),
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	// Dapatkan secret key dari .env
	jwtSecret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not login"})
	}

	return c.JSON(fiber.Map{
		"token": t,
		"role":  user.Role,
	})
}