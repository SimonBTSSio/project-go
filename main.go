package main

import (
	"go-project/broadcast"
	"go-project/handlerPayment"
	"go-project/handlerProduct"
	"go-project/model"
	"go-project/payment"
	"go-project/product"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "user:password@tcp(127.0.0.1:3310)/go_project?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dbURL), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&model.Product{})
	db.AutoMigrate(&model.Payment{})
	//bc.DelChan()
	bc := broadcast.NewBroadcaster(10)
	//init product
	productRepository := product.NewRepository(db)
	productService := product.NewService(productRepository)
	productHandler := handlerProduct.NewProductHandler(productService)

	//init payment
	paymentRepository := payment.NewRepository(db)
	paymentService := payment.NewService(paymentRepository)
	paymentHandler := handlerPayment.NewPaymentHandler(paymentService, bc)

	r := gin.Default()
	api := r.Group("/api")

	//ROUTES PRODUCT
	api.POST("/product", productHandler.Create)
	api.GET("/products", productHandler.GetAll)
	api.GET("/product/:id", productHandler.GetById)
	api.PUT("/product/:id", productHandler.Update)
	api.DELETE("/product/:id", productHandler.Delete)

	//ROUTES PAYMENT
	api.POST("/payment", paymentHandler.Create)
	api.GET("/payments", paymentHandler.GetAll)
	api.GET("/payment/:id", paymentHandler.GetById)
	api.PUT("/payment/:id", paymentHandler.Update)
	api.DELETE("/payment/:id", paymentHandler.Delete)
	api.GET("/payment/stream", paymentHandler.Stream)

	r.Run(":3000")
}
