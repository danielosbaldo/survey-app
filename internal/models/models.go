package models

import (
	"time"
)

type Shop struct {
	ID        uint        `json:"id" gorm:"primaryKey"`
	Name      string      `json:"name"`
	Slug      string      `json:"slug"`
	CiudadID  uint        `json:"ciudad_id"`
	Ciudad    Ciudad      `json:"ciudad" gorm:"foreignKey:CiudadID"`
	Employees []*Employee `gorm:"many2many:employee_shops;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Employee struct {
	ID        uint    `gorm:"primaryKey"`
	Name      string  `gorm:"not null"`
	Active    bool    `gorm:"default:true"`
	Shops     []*Shop `gorm:"many2many:employee_shops;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Tabla intermedia para la relación muchos-a-muchos
type EmployeeShop struct {
	EmployeeID uint `gorm:"primaryKey"`
	ShopID     uint `gorm:"primaryKey"`
}

type Question struct {
	ID        uint   `gorm:"primaryKey"`
	Prompt    string `gorm:"not null"`
	Type      string `gorm:"not null"` // radio, scale, text
	OrderNum  int    `gorm:"index"`
	Choices   []Choice
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Choice struct {
	ID         uint   `gorm:"primaryKey"`
	QuestionID uint   `gorm:"index;not null"`
	Label      string `gorm:"not null"`
	Value      string `gorm:"not null"`
	OrderNum   int
}

type Response struct {
	ID         uint  `gorm:"primaryKey"`
	ShopID     uint  `gorm:"index;not null"`
	EmployeeID uint  `gorm:"index;not null"`
	Answers    JSONB `gorm:"type:jsonb"`
	UserAgent  string
	CreatedAt  time.Time
}

type JSONB map[string]any
