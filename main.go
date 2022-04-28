package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin" // web framework API
	"github.com/joho/godotenv" // env
	"github.com/nuttakarn-nj/golang-todo/auth"
	"github.com/nuttakarn-nj/golang-todo/todo"
	"golang.org/x/time/rate"
	"gorm.io/driver/sqlite" // driver
	"gorm.io/gorm"          // ORM library for Golang
)

// variable for build cmd
var (
	buildcommit = "dev"
	buildTime   = time.Now().String()
)

// rate limit
var limiter = rate.NewLimiter(5, 5) // limit/sec and burst size 

func limiterHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	c.JSON(200, gin.H{
		"message": "allowed",
	})
}

func main() {
	// liveness check
	_, err := os.Create("/tmp/live")

	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live") // remove before end program

	// load env
	err = godotenv.Load("local.env")
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

	router := gin.Default()

	// middleware
	signkey := os.Getenv("SIGN")
	protected := router.Group("", auth.Protect([]byte(signkey)))

	// routes
	// realiness check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{})
	})

	// limit chack
	router.GET("/limit", limiterHandler)

	router.GET("/x", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"buildcommit": buildcommit,
			"buildTime":   buildTime,
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ponggggggg"})
	})

	router.GET("/token", auth.AccessToken(signkey))

	handler := todo.NewTodoHandler(db)
	protected.POST("/todos", handler.NewTask)

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	// router.Run()

	// Graceful shutdown start
	server := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// handle signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// goroutine
	go func() {
		err := server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// done when signal arrive
	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown
	if server.Shutdown(timeoutCtx) != nil {
		fmt.Println(err)
	}

	// Graceful shutdown end
}
