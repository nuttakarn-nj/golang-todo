package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin" // web framework API
	"github.com/joho/godotenv"
	"github.com/nuttakarn-nj/golang-todo/auth"
	"github.com/nuttakarn-nj/golang-todo/todo"
	"gorm.io/driver/sqlite" // driver
	"gorm.io/gorm"          // ORM library for Golang
)

func main() {
	// load env
	err := godotenv.Load("local.env")
	if err != nil {
		log.Printf("please consider environment variables: %s\n", err)
	}

	// open connection
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		// abort if a function returns an error value that we don’t know how to (or want to) handle
		panic(("failed to connect database"))
	}

	// Migrate the schema
	db.AutoMigrate((&todo.Todo{}))

	// Create
	// db.Create(&Product{Code: "D42", Price: 100})

	// Read
	// var product Product
	// db.First(&product, 1)                 // find product with integer primary key
	// db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	// db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	// db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	// db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - delete product
	// db.Delete(&product, 1)

	// routes
	router := gin.Default()

	// middleware
	signkey := os.Getenv("SIGN")
	protected := router.Group("", auth.Protect([]byte(signkey)))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ponggggggg"})
	})

	router.GET("/token", auth.AccessToken(signkey))

	handler := todo.NewTodoHandler(db)
	protected.POST("/todos", handler.NewTask)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	router.Run()
}
