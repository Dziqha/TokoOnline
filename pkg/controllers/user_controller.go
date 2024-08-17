package controllers

import (
	"Clone-TokoOnline/configs"
	"Clone-TokoOnline/pkg/models"
	"Clone-TokoOnline/pkg/response-codes"
	"Clone-TokoOnline/pkg/responses"
	"Clone-TokoOnline/pkg/utils"
	"encoding/base64"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController struct{}

func NewUserController() *UserController {
	return &UserController{}
}

func (controller *UserController)Register(c *fiber.Ctx) error {
	var req models.Users
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorBadRequest(c,err)
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorValidation(c,err)
	}

	// Validasi panjang password dan karakter yang diizinkan
	if len(req.Password) < 5 || len(req.Password) > 15 {
		return c.Status(fiber.StatusBadRequest).JSON(responsecodes.ResponseCode{
			Status:  fiber.StatusBadRequest,
			Message: "Password must be between 5 to 15 characters long.",
		})
	}

	// Validasi harus ada huruf
	hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(req.Password)
	// Validasi harus ada angka
	hasDigit := regexp.MustCompile(`\d`).MatchString(req.Password)
	// Validasi harus ada karakter khusus
	hasSpecialChar := regexp.MustCompile(`[@$!%*?&]`).MatchString(req.Password)

	if !hasLetter || !hasDigit || !hasSpecialChar {
		return c.Status(fiber.StatusBadRequest).JSON(responsecodes.ResponseCode{
			Status:  fiber.StatusBadRequest,
			Message: "Password must contain at least one letter, one number, and one special character (@, $, !, %, *, ?, &).",
		})
	}

	encode := base64.StdEncoding.EncodeToString([]byte(req.Password))

	user := models.Users{
		ID: req.ID,
		Username: req.Username,
		Password: encode,
	}
	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		tx.Create(&user)
		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
			Status: fiber.StatusInternalServerError,
			Message: "Failed to create user",
		})
	}

	res := responses.UserResponseRegister{
		Id: user.ID,
		Username: user.Username,
	}


	return c.JSON(responsecodes.ResponseCode{
		Status: fiber.StatusOK,
		Message: "Success",
		Data: res,
	})
}

func (controller *UserController) Login(c *fiber.Ctx) error {
	var TokenSecret = os.Getenv("TOKEN_SECRET")
	var req models.Users
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorBadRequest(c,err)
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorValidation(c,err)
	}

	encode := base64.StdEncoding.EncodeToString([]byte(req.Password))
		
	password := string(encode)
	
		var user models.Users
		err := configs.Database().Transaction(func(tx *gorm.DB) error {
			result := tx.Where("username = ? AND password = ?", req.Username, password).First(&user)
			if result.Error != nil {
				return result.Error
			}

			
			return nil
		})
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusUnauthorized).JSON(responsecodes.ResponseCode{
					Status:  fiber.StatusUnauthorized,
					Message: "Invalid username or password",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to login user",
			})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte(TokenSecret))
	if err != nil {
		return utils.ErrorInternalServerError(c, err)
	}

	cookie := &fiber.Cookie{
		Name: "myapp_user",
		Value: tokenString,
		Expires: time.Now().Add(time.Hour * 24),// Digunakan untuk mengatasi risiko CSRF 
		SameSite: fiber.CookieSameSiteStrictMode,
		HTTPOnly: true, // Cookie tidak bisa diakses dari JavaScript
		Secure: true, // Hanya dikirim melalui HTTPS
	}

	c.Cookie(cookie)
	res := responses.UserResponseLogin{
		Id: user.ID,
		Username: req.Username,
		Token: tokenString,
	}

	return c.JSON(responsecodes.ResponseCode{
		Status: fiber.StatusOK,
		Message: "Success",
		Data: res,
	})
}