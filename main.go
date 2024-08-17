package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"Clone-TokoOnline/configs"
	"Clone-TokoOnline/pkg/controllers"
	"Clone-TokoOnline/pkg/models"
	"Clone-TokoOnline/pkg/routes"
)

const migrationFolder = "configs/migrations"
const migrationFile = "migrations_done.txt"

func runMigrations(db *gorm.DB) {
    migrationFilePath := filepath.Join(migrationFolder, migrationFile)

    // Cek apakah folder configs dan migrations sudah ada, jika tidak buat
    if _, err := os.Stat("configs"); os.IsNotExist(err) {
        if err := os.Mkdir("configs", os.ModePerm); err != nil {
            log.Fatalf("Error creating configs folder: %v", err)
        }
    }

    if _, err := os.Stat(migrationFolder); os.IsNotExist(err) {
        if err := os.Mkdir(migrationFolder, os.ModePerm); err != nil {
            log.Fatalf("Error creating migrations folder: %v", err)
        }
    }

    // Cek apakah migrasi sudah dilakukan
    if _, err := os.Stat(migrationFilePath); os.IsNotExist(err) {
        log.Println("Running migrations...")
        
        models.MigrateOrder(db)
        models.MigrateProduct(db)
        models.MigrateCarts(db)
        models.MigrateUsers(db)

        // Tandai bahwa migrasi sudah dilakukan
        file, err := os.Create(migrationFilePath)
        if err != nil {
            log.Fatalf("Error creating migration status file: %v", err)
        }

        // Tambahkan pesan ke file
        message := "Migrations completed successfully on " + time.Now().Format(time.RFC3339) + "\n"
        if _, err := file.WriteString(message); err != nil {
            log.Fatalf("Error writing to migration status file: %v", err)
        }
        file.Close()
    } else {
        log.Println("Migrations already run.")
    }
}

func main() {
    app := fiber.New(fiber.Config{
        IdleTimeout:  time.Second * 10,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		Prefork:      false,
    })

    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    configs.Initialize()
    configs.ESClientConnection()
    configs.ESCreateIndexIfNotExist()
    db := configs.Database()

    runMigrations(db)

    usercontroller := controllers.NewUserController()
    productcontroller := controllers.NewProductController()
    ordercontroller := controllers.NewOrderController()
    cartscontroller := controllers.NewCartsController(ordercontroller)

    routes.NewRoutes(app, usercontroller)
    routes.NewRoutesProduct(app, productcontroller)
    routes.NewRoutesOrder(app, ordercontroller)
    routes.NewRoutesCarts(app, cartscontroller)

    // Rute untuk Swagger UI
    app.Get("/swagger/*", swagger.New(swagger.Config{
        URL: "/swagger.yaml",
    }))
    app.Static("/swagger.yaml", "./docs/swagger.yaml")
    err = app.Listen(":3000")
    if err != nil {
        log.Fatal("Failed to start server: ", err)
    }
}
