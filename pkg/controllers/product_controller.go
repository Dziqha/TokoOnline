package controllers

import (
	"Clone-TokoOnline/configs"
	"Clone-TokoOnline/pkg/models"
	"Clone-TokoOnline/pkg/response-codes"
	"Clone-TokoOnline/pkg/responses"
	"Clone-TokoOnline/pkg/utils"
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductController struct{}

func NewProductController() *ProductController {
	return &ProductController{}
}

func (controller *ProductController) InsertProduct(c *fiber.Ctx) error {
	var req models.Product

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorBadRequest(c,err)
	}

	if err := utils.ValidateStruct(req); err != nil {
		return utils.ErrorValidation(c,err)
	}

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		result := tx.Create(&req)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return utils.ErrorInternalServerError(c,err)
	}

	err = configs.AddProductToIndex(req)
	if err != nil {
		log.Fatalf("Error adding product to index: %s", err)
	}

	cacheKey := os.Getenv("CACHE_KEY_INSERT_PRODUCT")
    if err := configs.RedisClient.Del(c.Context(), cacheKey).Err(); err != nil {
        return utils.ErrorInternalServerError(c, err)
    }

	res := responses.Product{
		Id: req.ID,
		Name: req.Name,
		Description: req.Description,
		Price: req.Price,
		Stock: req.Stock,
	}

	return c.JSON(responsecodes.ResponseCode{
		Status: fiber.StatusOK,
		Message: "Success",
		Data: res,
	})
}

func (controller *ProductController) FindByIdProduct(c *fiber.Ctx) error {
    startTime := time.Now()

    id := c.Params("id")
    cacheKey := os.Getenv("CACHE_KEY_PRODUCT_PREFIX") + id

    // Cek cache Redis terlebih dahulu
    val, err := configs.RedisClient.Get(c.Context(), cacheKey).Result()
    if err == nil && val != "" {
        var data responses.Product
        if err := json.Unmarshal([]byte(val), &data); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
                Status:  fiber.StatusInternalServerError,
                Message: "Failed to parse cache data",
                Data:    nil,
            })
        }
        c.Set("X-Source", "cache")
        log.Printf("FindById from cache took %v", time.Since(startTime))
        return c.JSON(responsecodes.ResponseCode{
            Status:  fiber.StatusOK,
            Message: "Product retrieved from cache",
            Data:    data,
        })
    }
	
    var res models.Product
    err = configs.Database().Transaction(func(tx *gorm.DB) error {
        result := tx.Where("id = ?", id).First(&res)
        if result.Error != nil {
            return result.Error
        }
        return nil
    })

    if err != nil {
        return utils.ErrorNotFoundProduct(c, err)
    }

    result := responses.Product{
        Id:          res.ID,
        Name:        res.Name,
        Description: res.Description,
        Price:       res.Price,
        Stock:       res.Stock,
    }

    // Simpan data ke Redis cache
    cacheData, err := json.Marshal(result)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
            Status:  fiber.StatusInternalServerError,
            Message: "Failed to marshal data for caching",
            Data:    nil,
        })
    }
    err = configs.RedisClient.Set(c.Context(), cacheKey, cacheData, 5*time.Minute).Err()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
            Status:  fiber.StatusInternalServerError,
            Message: "Failed to cache product",
            Data:    nil,
        })
    }

    c.Set("X-Source", "database")
    log.Printf("FindById from database took %v", time.Since(startTime))
    return c.JSON(responsecodes.ResponseCode{
        Status:  fiber.StatusOK,
        Message: "Product found",
        Data:    result,
    })
}

func (controller *ProductController) FindAllProduct(c *fiber.Ctx) error {
	startTime := time.Now()

    cacheKey := os.Getenv("CACHE_KEY_PRODUCT_ALL")

    // Cek cache Redis terlebih dahulu
    val, err := configs.RedisClient.Get(c.Context(), cacheKey).Result()
    if err == nil && val != "" {
        var data []responses.Product
        if err := json.Unmarshal([]byte(val), &data); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
                Status:  fiber.StatusInternalServerError,
                Message: "Failed to parse cache data",
                Data:    nil,
            })
        }
        c.Set("X-Source", "cache")
        log.Printf("FindAll from cache took %v", time.Since(startTime))
        return c.JSON(responsecodes.ResponseCode{
            Status:  fiber.StatusOK,
            Message: "Product retrieved from cache",
            Data:    data,
        })
    }
	var resModels []models.Product
	var responseData []responses.Product

	err = configs.Database().Transaction(func(tx *gorm.DB) error {
		result := tx.Find(&resModels)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return utils.ErrorInternalServerError(c,err)
	}

	for _, resModel := range resModels {
		responseData = append(responseData, responses.Product{
			Id:          resModel.ID,
			Name:        resModel.Name,
			Description: resModel.Description,
			Price:       resModel.Price,
			Stock:       resModel.Stock,
		})
	}

	cacheData, err := json.Marshal(responseData)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
            Status:  fiber.StatusInternalServerError,
            Message: "Failed to marshal data for caching",
            Data:    nil,
        })
    }
    err = configs.RedisClient.Set(c.Context(), cacheKey, cacheData, 5*time.Minute).Err()
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
        Status:  fiber.StatusOK,
        Message: "Product found",
        Data:    responseData,
    })
}


func (controller *ProductController) UpdatedProduct(c *fiber.Ctx) error {
	startTime := time.Now()

    id := c.Params("id")
    cacheKey := os.Getenv("CACHE_KEY_PRODUCT_PREFIX") + id

    // Cek cache Redis terlebih dahulu
    val, err := configs.RedisClient.Get(c.Context(), cacheKey).Result()
    if err == nil && val != "" {
        var data responses.Product
        if err := json.Unmarshal([]byte(val), &data); err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
                Status:  fiber.StatusInternalServerError,
                Message: "Failed to parse cache data",
                Data:    nil,
            })
        }
        c.Set("X-Source", "cache")
        log.Printf("FindById from cache took %v", time.Since(startTime))
        return c.JSON(responsecodes.ResponseCode{
            Status:  fiber.StatusOK,
            Message: "Product retrieved from cache",
            Data:    data,
        })
    }
	var req models.Product

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorBadRequest(c,err)
	}

	var product models.Product
	if err := configs.Database().Where("id = ?", id).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorNotFoundProduct(c,err)
		}
		return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
			Status:  fiber.StatusInternalServerError,
			Message: "Error fetching product",
		})
	}

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 { 
		product.Price = req.Price
	}
	if req.Stock > 0 { 
		product.Stock = req.Stock
	}

	 err = configs.Database().Transaction(func(tx *gorm.DB) error {
		result := tx.Save(&product)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}) 


	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
			Status:  fiber.StatusInternalServerError,
			Message: "Error updating product",
		})
	}

	result := responses.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
	}

	err = configs.UpdatedProductToIndex(configs.ESClient, id, result)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
			Status:  fiber.StatusInternalServerError,
			Message: "Error updating product to index",
		})
	}

	cacheData, err := json.Marshal(result)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
            Status:  fiber.StatusInternalServerError,
            Message: "Failed to marshal data for caching",
            Data:    nil,
        })
    }
    err = configs.RedisClient.Set(c.Context(), cacheKey, cacheData, 5*time.Minute).Err()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
            Status:  fiber.StatusInternalServerError,
            Message: "Failed to cache product",
            Data:    nil,
        })
    }

    c.Set("X-Source", "database")
    log.Printf("FindAll from database took %v", time.Since(startTime))
	return c.Status(fiber.StatusOK).JSON(responsecodes.ResponseCode{
		Status:  fiber.StatusOK,
		Message: "Product updated successfully",
		Data:    result,
	})
}


func (controller *ProductController) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product
	if err := configs.Database().Where("id = ?", id).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.ErrorNotFoundProduct(c,err)	
		}
	}

	err := configs.Database().Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&product)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
			Status: fiber.StatusInternalServerError,
			Message: "Error deleting product",
		})
	}

	err = configs.DeleteProductFromIndex(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responsecodes.ResponseCode{
			Status: fiber.StatusInternalServerError,
			Message: "Error deleting product from index",
		})
	}

	cacheKey := os.Getenv("CACHE_KEY_PRODUCT_PREFIX") + id
    configs.RedisClient.Del(c.Context(), cacheKey)

    configs.RedisClient.Del(c.Context(), os.Getenv("CACHE_KEY_PRODUCT_ALL"))


	return c.JSON(responsecodes.ResponseCode{
		Status:  fiber.StatusOK,
		Message: "Product deleted successfully",
	})
}


func (controller *ProductController) SearchProduct(c *fiber.Ctx) error {
	var requestBody map[string]string
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Invalid request body",
		})
	}

	query, ok := requestBody["q"]
	if !ok || query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Query parameter 'q' is required in body",
		})
	}

	resp, err := esapi.SearchRequest{
		Index: []string{configs.SearchIndex},
		Body:  strings.NewReader(`{"query": {"match": {"name": "` + query + `"}}}`),
	}.Do(context.Background(), configs.ESClient)

	if err != nil {
		log.Printf("Error executing search request: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to search product",
		})
	}

	defer resp.Body.Close()

	var esResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&esResponse); err != nil {
		log.Printf("Error decoding search response: %s", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to search product",
		})
	}
	hits := esResponse["hits"].(map[string]interface{})["hits"].([]interface{})
	totalResults := int(esResponse["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	var products []responses.Product
	for _, hit := range hits {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})
		product := responses.Product{
			Id:          hit.(map[string]interface{})["_id"].(string),
			Name:        source["name"].(string),
			Description: source["description"].(string),
		}
		if price, ok := source["price"].(float64); ok {
			product.Price = int(price)
		}
		if stock, ok := source["stock"].(float64); ok {
			product.Stock = int(stock)
		}
		products = append(products, product)
	}

	result := models.SearchResult{
		Status:  fiber.StatusOK,
		Message: "Product search result",
		Data: struct {
			TotalResults int        `json:"total_results"`
			Hits         []responses.Product `json:"hits"`
		}{
			TotalResults: totalResults,
			Hits:         products,
		},
	}

	return c.JSON(result)
}

