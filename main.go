package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Task struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

func main() {
	db, err := gorm.Open(sqlite.Open("app.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Task{}, &Category{})

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
			"https://to-do-list-tau-taupe-69.vercel.app",
		},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/api/tasks", func(c *gin.Context) {
		var tasks []Task
		db.Order("created_at desc").Find(&tasks)
		c.JSON(http.StatusOK, tasks)
	})

	r.POST("/api/tasks", func(c *gin.Context) {
		var task Task
		if err := c.ShouldBindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task.CreatedAt = time.Now()
		db.Create(&task)
		c.JSON(http.StatusCreated, task)
	})

	r.DELETE("/api/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		db.Delete(&Task{}, id)
		c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
	})

	r.GET("/api/categories", func(c *gin.Context) {
		var categories []Category
		db.Order("created_at desc").Find(&categories)
		c.JSON(http.StatusOK, categories)
	})

	r.POST("/api/categories", func(c *gin.Context) {
		var category Category
		if err := c.ShouldBindJSON(&category); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		category.CreatedAt = time.Now()
		db.Create(&category)
		c.JSON(http.StatusCreated, category)
	})

	r.DELETE("/api/categories/:id", func(c *gin.Context) {
		id := c.Param("id")
		db.Delete(&Category{}, id)
		c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "5001"
	}
	r.Run(":" + port)

}
