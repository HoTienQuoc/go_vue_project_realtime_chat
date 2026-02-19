package controllers

import (
	"Server/database"
	"Server/models"
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

// Register
// @Summary Register a new user
// @Description Register a new user by providing email, password, first name, last name
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "user register details"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} map[string]interface{}
// @Router /user/signup [post]
func Register(c fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.CreateUser
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	CheckUser := UserSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&body)

	if CheckUser == nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "user with email " + body.Email + ". Already Exist!",
		})
	}
	// Hashing password
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	newUser := models.UserModel{
		Name:      body.FirstName + " " + body.LastName,
		Email:     body.Email,
		Password:  string(hashPassword),
		Followers: make([]string, 0),
		Following: make([]string, 0),
	}

	result, err := UserSchema.InsertOne(ctx, &newUser)

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(err)
	}

	// Get the new User
	var createdUser *models.UserModel
	query := bson.M{"_id": result.InsertedID}

	UserSchema.FindOne(ctx, query).Decode(&createdUser)

	// Create jwt token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    createdUser.ID.Hex(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	JWTSecret := os.Getenv("JWT_SECRET")

	token, _ := claims.SignedString([]byte(JWTSecret))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": createdUser,
		"token":  token,
	})
}

// Login
// @Summary Login a user
// @Description Login a new user by providing email, password, first name, last name
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.LoginUser true "user login details"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} map[string]interface{}
// @Router /user/signup [post]
func Login(c fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.LoginUser
	if err := c.Bind().Body(&body); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	var user models.UserModel
	CheckEmail := UserSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&user)

	// check if user with provided email exist or not
	if CheckEmail != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "user with email " + body.Email,
		})
	}

	// check if we have the same pass or not
	CheckPass := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if CheckPass != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"message": "user with email " + body.Email,
		})
	}

	// Create jwt token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    user.ID.Hex(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	JWTSecret := os.Getenv("JWT_SECRET")

	token, _ := claims.SignedString([]byte(JWTSecret))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": user,
		"token":  token,
	})
}
