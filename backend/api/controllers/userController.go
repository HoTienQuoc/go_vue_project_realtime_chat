package controllers

import (
	"Server/database"
	"Server/models"
	"context"
	"slices"
	"sort"
	"time"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// GetUserById
// @Summary Get User By Id
// @Description Get User By Id a new user by providing email, password, first name, last name
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} map[string]interface{}
// @Router /user/getUser/{id} [Get]
func GetUserById(c fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var user models.UserModel
	objId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	// strId := c.Params("id")
	// Todo get and return user posts

	// Get user data
	userResult := UserSchema.FindOne(ctx, bson.M{"_id": objId})
	if userResult.Err() != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"success": false,
			"message": "User Not Found",
		})
	}

	userResult.Decode(&user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":  user,
		"posts": "posts",
	})
}

// UpdateUser
// @Summary Update User
// @Description Update User Details
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /user/Update/{id} [Patch]
func UpdateUser(c fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//
	extUid := c.Locals("userId").(string)

	if extUid != c.Params("id") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "You Are Not Authroized to Update This Profile",
		})
	}

	userId, _ := primitive.ObjectIDFromHex(c.Params("id"))

	var user models.UpdateUser
	if err := c.Bind().Body(&user); err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	update := bson.M{"name": user.Name, "imageUrl": user.ImageUrl, "bio": user.Bio}

	result, err := UserSchema.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$set": update})

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error":   "Cannot update the user data",
			"details": err.Error(),
		})
	}

	//
	var updateUser models.UserModel
	if result.MatchedCount == 1 {
		err := UserSchema.FindOne(ctx, bson.M{"_id": userId}).Decode(&updateUser)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Cannot update the user data",
				"details": err.Error(),
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": updateUser,
	})
}

// Following Users
// @Summary Follow/UnFollow User
// @Description follow or unfollow an user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User Id"
// @Success 201 {object} models.UserModel
// @Failure 400 {object} map[string]interface{}
// @Security BearerAuth
// @Router /user/Update/{id} [Patch]
func FollowingUser(c fiber.Ctx) error {
	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var FirstUser models.UserModel
	var SecondUser models.UserModel

	FirstUserId, _ := primitive.ObjectIDFromHex(c.Params("id"))
	SecondUserID, _ := primitive.ObjectIDFromHex(c.Locals("userId").(string))

	err := UserSchema.FindOne(ctx, bson.M{"_id": FirstUserId}).Decode(&FirstUser)
	err = UserSchema.FindOne(ctx, bson.M{"_id": SecondUserID}).Decode(&SecondUser)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"details": err.Error(),
		})
	}

	fuid := c.Params("id")
	suid := c.Locals("userId").(string)

	if slices.Contains(FirstUser.Followers, suid) {
		i := sort.SearchStrings(FirstUser.Followers, suid)
		FirstUser.Followers = slices.Delete(FirstUser.Followers, i, i+1)
		// remove form the following list on second user
		index := sort.SearchStrings(SecondUser.Followers, fuid)
		FirstUser.Followers = slices.Delete(FirstUser.Followers, i, i+1)
	}
	// //
	// extUid := c.Locals("userId").(string)

	// if extUid != c.Params("id") {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"success": false,
	// 		"message": "You Are Not Authroized to Update This Profile",
	// 	})
	// }

	// userId, _ := primitive.ObjectIDFromHex(c.Params("id"))

	// var user models.UpdateUser
	// if err := c.Bind().Body(&user); err != nil {
	// 	return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
	// 		"error":   "Invalid request body",
	// 		"details": err.Error(),
	// 	})
	// }

	// update := bson.M{"name": user.Name, "imageUrl": user.ImageUrl, "bio": user.Bio}

	// result, err := UserSchema.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$set": update})

	// if err != nil {
	// 	return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
	// 		"error":   "Cannot update the user data",
	// 		"details": err.Error(),
	// 	})
	// }

	// //
	// var updateUser models.UserModel
	// if result.MatchedCount == 1 {
	// 	err := UserSchema.FindOne(ctx, bson.M{"_id": userId}).Decode(&updateUser)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
	// 			"error":   "Cannot update the user data",
	// 			"details": err.Error(),
	// 		})
	// 	}
	// }

	// return c.Status(fiber.StatusOK).JSON(fiber.Map{
	// 	"data": updateUser,
	// })
}
