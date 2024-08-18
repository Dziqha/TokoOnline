package controllers

import (
	"Clone-TokoOnline/configs"
	"Clone-TokoOnline/pkg/models"
	"Clone-TokoOnline/pkg/response-codes"
	"Clone-TokoOnline/pkg/responses"
	"Clone-TokoOnline/pkg/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type OrderController struct{}

func NewOrderController() *OrderController {
	return &OrderController{}
}

func (controller *OrderController) CreateOrder(c *fiber.Ctx, carts []models.Carts, isCheckout bool) error {
	var orderResponses []responses.Order

	// Jika tidak melakukan checkout, ambil data dari request body
	if !isCheckout {
		var req models.Orders
		if err := c.BodyParser(&req); err != nil {
			return utils.ErrorBadRequest(c, err)
		}

		if err := utils.ValidateStruct(req); err != nil {
			return utils.ErrorValidation(c, err)
		}

		var product models.Product
		err := configs.Database().Transaction(func(tx *gorm.DB) error {
			result := tx.First(&product, "id = ?", req.ProductId)
			if result.Error != nil {
				return result.Error
			}
			return nil
		})
		if err != nil {
			return utils.ErrorNotFoundProduct(c, err)
		}

		req.TotalPrice = req.Quantity * product.Price

		err = configs.Database().Transaction(func(tx *gorm.DB) error {
			result := tx.Create(&req)
			if result.Error != nil {
				return result.Error
			}
			return nil
		})
		if err != nil {
			return utils.ErrorInternalServerError(c, err)
		}
		if req.Status == "" {
			req.Status = "pending"
		}
		orderResponses = append(orderResponses, responses.Order{
			Id:         req.ID,
			UserId:     req.UserId,
			ProductId:  req.ProductId,
			Quantity:   req.Quantity,
			TotalPrice: req.TotalPrice,
			Status:     req.Status,
		})
	} else {
		// Jika melakukan checkout, proses cart items
		if len(carts) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(responsecodes.ResponseCode{
				Status:  fiber.StatusBadRequest,
				Message: "No carts provided",
			})
		}

		userId, ok := c.Locals("userId").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(responsecodes.ResponseCode{
				Status:  fiber.StatusUnauthorized,
				Message: "User ID not found in context",
			})
		}

		for _, cart := range carts {
			var product models.Product
			err := configs.Database().Transaction(func(tx *gorm.DB) error {
				result := tx.First(&product, "id = ?", cart.ProductId)
				if result.Error != nil {
					return result.Error
				}
				return nil
			})
			if err != nil {
				return utils.ErrorNotFoundProduct(c, err)
			}

			order := models.Orders{
				UserId:     userId,
				ProductId:  cart.ProductId,
				Quantity:   cart.Quantity,
				TotalPrice: cart.Quantity * product.Price,
				Status:     "pending",
			}

			err = configs.Database().Transaction(func(tx *gorm.DB) error {
				result := tx.Create(&order)
				if result.Error != nil {
					return result.Error
				}
				return nil
			})
			if err != nil {
				return utils.ErrorInternalServerError(c, err)
			}

			orderResponses = append(orderResponses, responses.Order{
				Id:         order.ID,
				UserId:     order.UserId,
				ProductId:  order.ProductId,
				Quantity:   order.Quantity,
				TotalPrice: order.TotalPrice,
				Status:     order.Status,
			})
		}
	}
	return c.JSON(responsecodes.ResponseCode{
		Status:  fiber.StatusOK,
		Message: "Orders created successfully",
		Data:    orderResponses,
	})
}

func (controller *OrderController) NewOrder(c *fiber.Ctx) error {
	var carts []models.Carts
	var req models.Orders
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorBadRequest(c, err)
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorValidation(c, err)
	}

	// Tidak menggunakan cart dalam kasus ini
	err := controller.CreateOrder(c, carts, false)
    if err != nil {
        return utils.ErrorInternalServerError(c, err)
    }

	err = configs.OrderQueue(req.ID)
	if err != nil {
		return utils.ErrorInternalServerError(c, err)
	}

    cacheKey := os.Getenv("CACHE_KEY_ORDERS")
    configs.RedisClient.Del(c.Context(), cacheKey)
	return nil
}


func (controller *OrderController) ViewOrderAll(c *fiber.Ctx) error {
	startTime := time.Now()

	cacheKey := os.Getenv("CACHE_KEY_ORDERS_ALL")
	val, err := configs.RedisClient.Get(c.Context(), cacheKey).Result()
	if err == nil && val != "" {
		var data []responses.Order
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

	
	var res []models.Orders
	var resData []responses.Order
	err = configs.Database().Transaction(func(tx *gorm.DB) error {
		result := tx.Find(&res)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return utils.ErrorInternalServerError(c,err)
	}

	for _, responseOrder := range res {
		resData = append(resData, responses.Order{
			Id: responseOrder.ID,
			UserId: responseOrder.UserId,
			ProductId: responseOrder.ProductId,
			Quantity: responseOrder.Quantity,
			TotalPrice: responseOrder.TotalPrice,
			Status: responseOrder.Status,
		})
	}

	cachedata, err := json.Marshal(resData)
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
		Message: "Success",
		Data: resData,
	})
}

func (controller *OrderController) DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Orders
	if err := configs.Database().First(&order, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorNotFoundOrder(c,err)	
		}
	}

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&order).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return utils.ErrorInternalServerError(c,err)
	}

	err = configs.DeleteQueue(id)
	if err != nil {
		return utils.ErrorInternalServerError(c,err)
	}

	cacheKey := os.Getenv("CACHE_KEY_ORDERS_PREFIX") + id
    configs.RedisClient.Del(c.Context(), cacheKey)

    configs.RedisClient.Del(c.Context(), os.Getenv("CACHE_KEY_ORDERS_ALL"))

	return c.JSON(responsecodes.ResponseCode{
		Status: fiber.StatusOK,
		Message: "Success",
	})
}

func (controller *OrderController) CancelOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Orders
	if err := configs.Database().Where("id = ?", id).First(&order).Error; err != nil {
		return utils.ErrorNotFoundOrder(c,err)
	}

	if order.Status == "canceled" {
		return utils.ErrorAlreadyCanceled(c)
	}

	if err := configs.Database().Model(&order).Update("status", "canceled").Error; err != nil {
		return utils.ErrorInternalServerError(c,err)
	}

	if err := configs.Database().Where("id = ?", id).Delete(&order).Error; err != nil {
		return utils.ErrorInternalServerError(c,err)
	}

	err := configs.CancelQueue(id)
	if err != nil {
		return utils.ErrorInternalServerError(c,err)
	}
	
	return c.JSON(responsecodes.ResponseCode{
		Status: fiber.StatusOK,
		Message: "Order successfully canceled",
	})
}

func (controller *OrderController) CekOngkir(c *fiber.Ctx) error {
    var request struct {
        Origin      string `json:"origin"`
        Destination string `json:"destination"`
        Weight      int    `json:"weight"`
    }

    if err := c.BodyParser(&request); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  fiber.StatusBadRequest,
            "message": "Invalid request body",
        })
    }

    getCityID := func(cityName string) (int, error) {
        url := "https://api.rajaongkir.com/starter/city?key=637db82739e1241882f2c3ba7d3a1ea6" // Ganti dengan kunci API Anda
        resp, err := http.Get(url)
        if err != nil {
            return 0, err
        }
        defer resp.Body.Close()

        var result struct {
            RajaOngkir struct {
                Results []struct {
                    CityID   string `json:"city_id"`
                    CityName string `json:"city_name"`
                } `json:"results"`
            } `json:"rajaongkir"`
        }

        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
            return 0, err
        }

        for _, city := range result.RajaOngkir.Results {
            if strings.EqualFold(city.CityName, cityName) {
                id, err := strconv.Atoi(city.CityID)
                if err != nil {
                    return 0, err
                }
                return id, nil
            }
        }

        return 0, fmt.Errorf("city not found")
    }

    originID, err := getCityID(request.Origin)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  fiber.StatusBadRequest,
            "message": "Invalid origin city",
        })
    }

    destinationID, err := getCityID(request.Destination)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  fiber.StatusBadRequest,
            "message": "Invalid destination city",
        })
    }

    couriers := []string{"jne", "pos", "tiki"}
    var allCosts []fiber.Map

    for _, courier := range couriers {
        payload := map[string]string{
            "origin":      strconv.Itoa(originID),
            "destination": strconv.Itoa(destinationID),
            "weight":      strconv.Itoa(request.Weight),
            "courier":     courier,
        }

        jsonPayload, err := json.Marshal(payload)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  fiber.StatusInternalServerError,
                "message": "Failed to marshal request payload",
            })
        }

        url := "https://api.rajaongkir.com/starter/cost"
        req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  fiber.StatusInternalServerError,
                "message": "Failed to create HTTP request",
            })
        }

        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("key", "637db82739e1241882f2c3ba7d3a1ea6") // Ganti dengan kunci API Anda

        client := &http.Client{}
        resp, err := client.Do(req)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  fiber.StatusInternalServerError,
                "message": "Failed to make API request",
            })
        }
        defer resp.Body.Close()

        var apiResponse struct {
            RajaOngkir struct {
                Results []struct {
                    Code  string `json:"code"`
                    Costs []struct {
                        Service     string `json:"service"`
                        Description string `json:"description"`
                        Cost        []struct {
                            Value int    `json:"value"`
                            Etd   string `json:"etd"`
                        } `json:"cost"`
                    } `json:"costs"`
                } `json:"results"`
            } `json:"rajaongkir"`
        }

        if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  fiber.StatusInternalServerError,
                "message": "Failed to decode API response",
            })
        }

        for _, result := range apiResponse.RajaOngkir.Results {
            for _, cost := range result.Costs {
                for _, detail := range cost.Cost {
                    allCosts = append(allCosts, fiber.Map{
                        "courier":     result.Code,
                        "service":     cost.Service,
                        "description": cost.Description,
                        "value":       detail.Value,
                        "etd":         detail.Etd,
                    })
                }
            }
        }
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":  fiber.StatusOK,
        "message": "Success",
        "data": fiber.Map{
            "origin_city":      request.Origin,
            "destination_city": request.Destination,
            "costs":            allCosts,
        },
    })
}
