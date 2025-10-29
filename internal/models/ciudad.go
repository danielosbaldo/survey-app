package models

type Ciudad struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Nombre string `json:"nombre"`
	Slug   string `json:"slug"`
	Region string `json:"region"`
	Pais   string `json:"pais"`
	Activa bool   `json:"activa"`
}

// TableName especifica el nombre de la tabla para el modelo Ciudad
func (Ciudad) TableName() string {
	return "ciudads"
}
