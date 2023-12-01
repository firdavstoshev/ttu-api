package models

// Projector структура представляет модель проектора
type Projector struct {
	ID       uint `gorm:"primaryKey"`
	Model    string
	Name     string
	Width    int
	Height   int
	Mode     string
	IsActive bool
}
