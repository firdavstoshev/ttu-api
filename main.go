package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"projector-api/models"
	"strconv"
)

var (
	db *gorm.DB
)

func main() {
	// Инициализация базы данных
	var err error
	dsn := "user=admin password=admin dbname=mydb sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	// Создание таблицы проекторов
	err = db.AutoMigrate(&models.Projector{})
	if err != nil {
		panic("Failed to migrate the database")
	}

	// Создание экземпляра Gin
	r := gin.Default()

	// Определение маршрутов
	r.POST("/projectors", createProjector)
	r.GET("/projectors/:id", getProjector)
	r.PUT("/projectors/:id/turnon", turnOnProjector)
	r.PUT("/projectors/:id/changemode", changeMode)
	r.PUT("/projectors/:id/changeresolution", changeResolution)

	// Запуск сервера
	r.Run(":8080")
}

func createProjector(c *gin.Context) {
	var projector models.Projector
	if err := c.ShouldBindJSON(&projector); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Создание записи в базе данных
	db.Create(&projector)

	c.JSON(200, projector)
}

func getProjector(c *gin.Context) {
	var projector models.Projector
	id := c.Params.ByName("id")

	// Поиск проектора по ID
	if err := db.First(&projector, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Projector not found"})
		return
	}

	c.JSON(200, projector)
}

func turnOnProjector(c *gin.Context) {
	var projector models.Projector
	id := c.Params.ByName("id")

	// Поиск проектора по ID
	if err := db.First(&projector, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Projector not found"})
		return
	}

	// Включение проектора
	projector.IsActive = true
	db.Save(&projector)

	c.JSON(200, projector)
}

func changeMode(c *gin.Context) {
	var projector models.Projector
	id := c.Params.ByName("id")

	// Поиск проектора по ID
	if err := db.First(&projector, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Projector not found"})
		return
	}

	// Смена режима
	newMode := c.Query("mode")
	projector.Mode = newMode
	db.Save(&projector)

	c.JSON(200, projector)
}

func changeResolution(c *gin.Context) {
	var projector models.Projector
	id := c.Params.ByName("id")

	// Поиск проектора по ID
	if err := db.First(&projector, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Projector not found"})
		return
	}

	// Изменение разрешения
	newWidth := c.Query("width")
	newHeight := c.Query("height")

	if newWidth != "" {
		projector.Width, _ = strconv.Atoi(newWidth)
	}

	if newHeight != "" {
		projector.Height, _ = strconv.Atoi(newHeight)
	}

	db.Save(&projector)

	c.JSON(200, projector)
}
