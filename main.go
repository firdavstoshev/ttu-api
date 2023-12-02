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
	r.PUT("/projectors/:id/turnon", turnOnProjector)
	r.PUT("/projectors/:id/changemode", changeMode)
	r.PUT("/projectors/:id/changeresolution", changeResolution)

	r.Run(":8080")
}

func createProjector(c *gin.Context) {
	var projector models.Projector
	if err := c.ShouldBindJSON(&projector); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Создание записи в базе данных
	if err := db.Create(&projector).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, projector)
}

func updateProjector(c *gin.Context) {
	var projector models.Projector
	id := c.Params.ByName("id")

	// Поиск проектора по ID
	if err := db.First(&projector, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Projector not found"})
		return
	}

	// Обновление параметров проектора
	if err := c.ShouldBindJSON(&projector); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.Save(&projector)

	c.JSON(200, projector)
}

func getAllProjectors(c *gin.Context) {
	var projectors []models.Projector

	// Извлечение всех проекторов из базы данных
	if err := db.Find(&projectors).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch projectors"})
		return
	}

	c.JSON(200, projectors)
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
	projector.IsActive = !projector.IsActive
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
