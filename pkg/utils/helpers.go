package utils

import "github.com/gofiber/fiber/v2"

func ErrorCecker(err error) error {
	if err != nil {
		panic(err)
	}

	return nil
}

func ErrorMigrate(err error) error {
	if err != nil {
		panic("Failed to migrate database schema: " + err.Error())
	}
	return nil
}

func ErrorNotFoundProduct(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Product not found",
			"message": err.Error(),
		})
	}
	return nil
}

func ErrorInternalServerError(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	return nil
}

func ErrorNotFoundOrder(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Order not found",
			"message": err.Error(),
		})
	}
	return nil
}

func ErrorNotFoundUser(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "User not found",
			"message": err.Error(),
		})
	}
	return nil
}

func ErrorValidation(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation error",
			"message": err.Error(),
		})
	}
	return nil
}

func ErrorBadRequest(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad request",
			"message": err.Error(),
		})
	}
	return nil
}

func ErrorAlreadyCanceled(c *fiber.Ctx) error {
	
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Already canceled",
		})
	
}