package todo

import (
	"log"
	"net/http" // http.StatusOK

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// A struct is a collection of fields.
type Todo struct {
	Title string `json:"text" binding:"required"`
	gorm.Model
	// type Model struct {
	// 	ID        uint           `gorm:"primaryKey"`
	// 	CreatedAt time.Time
	// 	UpdatedAt time.Time
	// 	DeletedAt gorm.DeletedAt `gorm:"index"`
	//   }
}

func (Todo) TableName() string {
	return "todos"
}

type TodoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

func (t *TodoHandler) NewTask(c *gin.Context) {
	var todo Todo
	// parse req form context into todo
	err := c.ShouldBindJSON(&todo)

	// have error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// test logging
	if todo.Title == "sleep" {
		transactionID := c.Request.Header.Get("transactionID")
		aud, _ := c.Get("aud")
		log.Println(transactionID, aud, "Not allowed")

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Not allowed this task",
		})
		return
	}

	result := t.db.Create(&todo)
	err = result.Error

	// have error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// response
	c.JSON(http.StatusCreated, gin.H{"ID": todo.Model.ID})
}
