package routes

import (
	"Clone-TokoOnline/pkg/controllers"
	middleware "Clone-TokoOnline/pkg/middlewares"

	"github.com/gofiber/fiber/v2"
)

func NewRoutes(router fiber.Router, usercontroller *controllers.UserController){
	app := router.Group("/api/users")
	app.Post("/register", usercontroller.Register)
	app.Post("/login", usercontroller.Login)

}

func NewRoutesProduct(router fiber.Router, productcontroller *controllers.ProductController){
	app := router.Group("/api/products", middleware.AuthMiddleware)
	app.Post("/admin/insert-product", productcontroller.InsertProduct)
	app.Get("/user/find-by-id/:id", productcontroller.FindByIdProduct)
	app.Get("/user/find-all", productcontroller.FindAllProduct)
	app.Post("/user/search-product", productcontroller.SearchProduct)
	app.Put("/admin/update-product/:id", productcontroller.UpdatedProduct)
	app.Delete("/admin/delete-product/:id", productcontroller.DeleteProduct)
}

func NewRoutesOrder(router fiber.Router, ordercontroller *controllers.OrderController){
	app := router.Group("/api/orders", middleware.AuthMiddleware)
	app.Post("/user/create-order", ordercontroller.NewOrder )
	app.Get("/user/view-order-all", ordercontroller.ViewOrderAll)
	app.Delete("/user/delete-order/:id", ordercontroller.DeleteOrder)
	app.Delete("/user/cancel-order/:id", ordercontroller.CancelOrder)
}

func NewRoutesCarts(router fiber.Router, cartscontroller *controllers.CartsController){
	app := router.Group("/api/carts", middleware.AuthMiddleware)
	app.Post("/user/insert-item-to-cart", cartscontroller.InsertItemToCarts)
	app.Get("/user/view-cart", cartscontroller.ViewCarts)
	app.Post("/user/checkout-cart", cartscontroller.CheckOutCarts)
}