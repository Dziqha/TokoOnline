package controllers

import (
	"Clone-TokoOnline/configs"
	"Clone-TokoOnline/pkg/models"
	"Clone-TokoOnline/pkg/response-codes"
	"Clone-TokoOnline/pkg/responses"
	"Clone-TokoOnline/pkg/utils"
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CartsController struct {
    orderController *OrderController
}

func NewCartsController(orderController *OrderController) *CartsController {
    return &CartsController{orderController: orderController}
}

func (contoller *CartsController) InsertItemToCarts(c *fiber.Ctx) error {
	var req models.Carts
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorBadRequest(c, err)
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorValidation(c, err)
	}

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		var exitingCartItem models.Carts
		result := tx.Where("product_id = ? AND user_id = ?", req.ProductId, req.UserId).First(&exitingCartItem)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		//jika cart item ada
		if result.RowsAffected > 0 {
			err := tx.Model(&exitingCartItem).Update("quantity", exitingCartItem.Quantity + req.Quantity).Error
			if err != nil {
				return err
			}
		}else {
			//jika cart item tidak ada
			err := tx.Create(&req).Error
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return utils.ErrorInternalServerError(c, err)
	}

	return c.JSON(responsecodes.ResponseCode{
		Status: fiber.StatusOK,
		Message: "Item successfully added to cart",
	})
}

func (controller *CartsController) ViewCarts(c *fiber.Ctx) error {
	// Ambil ID pengguna dari token JWT (dapatkan dari middleware atau konteks)
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(responsecodes.ResponseCode{
			Status:  fiber.StatusUnauthorized,
			Message: "User ID not found in context",
		})
	}

	startTime := time.Now()

	cacheKey := os.Getenv("CACHE_KEY_ORDERS_ALL")
	val, err := configs.RedisClient.Get(c.Context(), cacheKey).Result()
	if err == nil && val != "" {
		var data []responses.Carts
		if err := json.Unmarshal([]byte(val), &data); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
				Status: fiber.StatusInternalServerError,
				Message: "Failed to parse cache data",
                Data:    nil,
			})
		}
		c.Set("X-Source", "cache")
		log.Printf("FindAll from cache took %v", time.Since(startTime))
		return c.JSON(responsecodes.ResponseCode{
		Status: fiber.StatusOK,
		Message: "Product retrieved from cache",
        Data:    data,
	})
	}
	var carts []models.Carts

	err = configs.Database().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("user_id = ?", userId).Find(&carts)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return utils.ErrorInternalServerError(c, err)
	}

	var res []responses.Carts
	for _,resdata  := range carts {
		res = append(res, responses.Carts{
			Id: resdata.ID,
			UserId: resdata.UserId,
			ProductId: resdata.ProductId,
			Quantity: resdata.Quantity,
		})
	}
	cachedata, err := json.Marshal(res)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
            Status:  fiber.StatusInternalServerError,
            Message: "Failed to marshal data for caching",
            Data:    nil,
        })
	}

	err = configs.RedisClient.Set(c.Context(), cacheKey, cachedata, 5 * time.Minute).Err()
	if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
            Status:  fiber.StatusInternalServerError,
            Message: "Failed to cache product",
            Data:    nil,
        })
    }
	
	c.Set("X-Source", "database")
    log.Printf("FindAll from database took %v", time.Since(startTime))
	return c.JSON(responsecodes.ResponseCode{
		Status: fiber.StatusOK,
		Message: "Success view carts",
		Data: res,
	})
}

func (controller *CartsController) CheckOutCarts(c *fiber.Ctx) error {
	// Ambil ID pengguna dari token JWT (dapatkan dari middleware atau konteks)
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(responsecodes.ResponseCode{
			Status:  fiber.StatusUnauthorized,
			Message: "User ID not found in context",
		})
	}

	var carts []models.Carts

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("user_id = ?", userId).Find(&carts)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return utils.ErrorInternalServerError(c, err)
	}

	if len(carts) == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(responsecodes.ResponseCode{
            Status:  fiber.StatusBadRequest,
            Message: "Cart is empty",
        })
    }
	if err := controller.orderController.CreateOrder(c,carts, true); err != nil {
        return err
    }

	err = configs.Database().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("user_id = ?", userId ).Delete(&models.Carts{})
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	if err != nil {
		return utils.ErrorInternalServerError(c, err)
	}

	err = configs.CheckOutQueue(userId)
	if err != nil {
		return utils.ErrorInternalServerError(c, err)
	}

	return c.JSON(responsecodes.ResponseCode{
		Status: fiber.StatusOK,
		Message: "Success checkout carts",
	})
}