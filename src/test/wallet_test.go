package test

import (
	"ankasa-be/src/configs"
	"ankasa-be/src/controllers"
	"ankasa-be/src/models"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDatabase() {
	url := "postgresql://postgres:@127.0.0.1:5432/ankasa_travel"
	var err error
	configs.DB, err = gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to test database!")
	}

	configs.DB.Exec("DROP TABLE IF EXISTS wallets CASCADE")
	configs.DB.Exec("DROP TABLE IF EXISTS customers CASCADE")
	configs.DB.Exec("DROP TABLE IF EXISTS users CASCADE")

	configs.DB.AutoMigrate(
		&models.User{},
		&models.Customer{},
		&models.Wallet{},
	)
}

func TestGetWallet(t *testing.T) {
	app := fiber.New()
	setupTestDatabase()

	user := models.User{Email: "test@gmail.com", Password: "test123", Role: "customer", IsVerify: "1"}
	userCreate := configs.DB.Create(&user)
	if userCreate.Error != nil {
		log.Fatalf("Failed to create user: %v", userCreate.Error)
	}

	customer := models.Customer{UserID: user.ID, Username: "jhon doe", PhoneNumber: "123456790", City: "Bandung"}
	customerCreate := configs.DB.Create(&customer)
	if customerCreate.Error != nil {
		log.Fatalf("Failed to create customer: %v", customerCreate.Error)
	}

	var createdCustomer models.Customer
	configs.DB.First(&createdCustomer, "id = ?", customer.ID)
	// log.Printf("Created customer: %+v", createdCustomer)

	wallet := models.Wallet{CustomerID: int(createdCustomer.ID), Saldo: 100000}
	result := configs.DB.Create(&wallet)
	if result.Error != nil {
		log.Fatalf("Failed to create wallet: %v", result.Error)
	}

	app.Get("/wallets", controllers.GetAllWallet)

	req := httptest.NewRequest("GET", "/wallets", nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)

	var response map[string][]models.Wallet
	json.NewDecoder(resp.Body).Decode(&response)
	wallets, _ := response["data"]
	// log.Printf("Wallets response: %+v", wallets)
	assert.Equal(t, 1, len(wallets))
	assert.Equal(t, createdCustomer.ID, wallets[0].CustomerID)
	assert.Equal(t, float64(100000), wallets[0].Saldo)
}

func TestUpdateWallet(t *testing.T) {
	app := fiber.New()
	setupTestDatabase()

	user := models.User{Email: "test@gmail.com", Password: "test123", Role: "customer", IsVerify: "1"}
	configs.DB.Create(&user)

	customer := models.Customer{UserID: user.ID, Username: "jhon doe", PhoneNumber: "123456790", City: "Bandung"}
	configs.DB.Create(&customer)

	wallet := models.Wallet{CustomerID: int(customer.ID), Saldo: 100000}
	configs.DB.Create(&wallet)

	app.Put("/wallets/:id", controllers.UpdateWallet)

	updatedSaldo := 200000
	reqBody := fmt.Sprintf(`{"saldo": %d}`, updatedSaldo)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/wallets/%d", wallet.ID), strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)

	var updatedWallet models.Wallet
	configs.DB.First(&updatedWallet, wallet.ID)
	assert.Equal(t, float64(updatedSaldo), float64(updatedWallet.Saldo))
}

func TestDeleteWallet(t *testing.T) {
	app := fiber.New()
	setupTestDatabase()

	user := models.User{Email: "test@gmail.com", Password: "test123", Role: "customer", IsVerify: "1"}
	configs.DB.Create(&user)

	customer := models.Customer{UserID: user.ID, Username: "jhon doe", PhoneNumber: "123456790", City: "Bandung"}
	configs.DB.Create(&customer)

	wallet := models.Wallet{CustomerID: int(customer.ID), Saldo: 100000}
	configs.DB.Create(&wallet)

	app.Delete("/wallets/:id", controllers.DeleteWallet)

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/wallets/%d", wallet.ID), nil)
	resp, _ := app.Test(req)
	assert.Equal(t, 200, resp.StatusCode)

	var deletedWallet models.Wallet
	result := configs.DB.First(&deletedWallet, wallet.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
}
