package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"projector-api/models"
	"strconv"
)

var (
	db *gorm.DB
)

func main() {
	var err error
	dsn := "user=admin password=admin dbname=mydb sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	err = db.AutoMigrate(&models.Projector{})
	if err != nil {
		panic("Failed to migrate the database")
	}

	r := gin.Default()

	r.GET("/projectors", getAllProjectors)
	r.POST("/projectors", createProjector)
	r.PUT("/projectors/:id", updateProjector)
	r.GET("/projectors/:id", getProjector)
	r.PUT("/projectors/:id/turn", turnOnProjector)
	r.PUT("/projectors/:id/change-mode", changeMode)
	r.PUT("/projectors/:id/change-resolution", changeResolution)

	r.Run(":8080")
}

// Создание записи в базе данных
func createProjector(c *gin.Context) {
	var projector models.Projector
	if err := c.ShouldBindJSON(&projector); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&projector).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, projector)
}

// Обновление параметров проектора
func updateProjector(c *gin.Context) {
	var projector models.Projector
	id := c.Params.ByName("id")

	if err := db.First(&projector, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Projector not found"})
		return
	}

	if err := c.ShouldBindJSON(&projector); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.Save(&projector)

	c.JSON(200, projector)
}

// Извлечение всех проекторов из базы данных
func getAllProjectors(c *gin.Context) {
	var projectors []models.Projector

	if err := db.Find(&projectors).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch projectors"})
		return
	}

	c.JSON(200, projectors)
}

// Поиск проектора по ID
func getProjector(c *gin.Context) {
	var projector models.Projector
	id := c.Params.ByName("id")

	if err := db.First(&projector, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Projector not found"})
		return
	}

	c.JSON(200, projector)
}

// Включение проектора
func turnOnProjector(c *gin.Context) {
	var projector models.Projector
	id := c.Params.ByName("id")

	if err := db.First(&projector, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Projector not found"})
		return
	}

	projector.IsActive = !projector.IsActive
	db.Save(&projector)

	c.JSON(200, projector)
}

// Смена режима
func changeMode(c *gin.Context) {
	var projector models.Projector
	id := c.Params.ByName("id")

	if err := db.First(&projector, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Projector not found"})
		return
	}

	newMode := c.Query("mode")
	projector.Mode = newMode
	db.Save(&projector)

	c.JSON(200, projector)
}

// Изменение разрешения
func changeResolution(c *gin.Context) {
	var projector models.Projector
	id := c.Params.ByName("id")

	if err := db.First(&projector, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Projector not found"})
		return
	}

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
