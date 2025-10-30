package seed

import (
	"github.com/danielosbaldo/survey-app/internal/models"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	var count int64
	db.Model(&models.Ciudad{}).Count(&count)
	if count == 0 {
		ciudades := []models.Ciudad{
			{Nombre: "Hermosillo"},
			{Nombre: "Navojoa"},
			{Nombre: "Jilotepec"},
		}
		db.Create(&ciudades)

		var hermosillo, navojoa, jilotepec models.Ciudad
		db.Where("nombre = ?", "Hermosillo").First(&hermosillo)
		db.Where("nombre = ?", "Navojoa").First(&navojoa)
		db.Where("nombre = ?", "Jilotepec").First(&jilotepec)

		// Crear una shop por ciudad
		shopjilotepec := models.Shop{Name: "Sucursal presidencia", Slug: "presidencia", CiudadID: jilotepec.ID}
		shopjilotepec2 := models.Shop{Name: "Sucursal Jilotepec 2", Slug: "hidalgo", CiudadID: jilotepec.ID}
		shopjilotepec3 := models.Shop{Name: "Sucursal Jilotepec 3", Slug: "guerrero-1", CiudadID: jilotepec.ID}
		shopjilotepec4 := models.Shop{Name: "Sucursal Jilotepec 4", Slug: "guerrero-2", CiudadID: jilotepec.ID}
		shopjilotepec5 := models.Shop{Name: "Sucursal Jilotepec 5", Slug: "isidro", CiudadID: jilotepec.ID}

		db.Create(&shopjilotepec)
		db.Create(&shopjilotepec2)
		db.Create(&shopjilotepec3)
		db.Create(&shopjilotepec4)
		db.Create(&shopjilotepec5)

		shop3 := models.Shop{Name: "Sucursal Jilotepec", Slug: "don-nico", CiudadID: jilotepec.ID}
		shop1 := models.Shop{Name: "Sucursal Hermosillo", Slug: "hermosillo", CiudadID: hermosillo.ID}
		shop2 := models.Shop{Name: "Sucursal Navojoa", Slug: "navojoa", CiudadID: navojoa.ID}

		db.Create(&shop1)
		db.Create(&shop2)
		db.Create(&shop3)

		// Crear empleados y asignar a una sucursal cada uno
		maria := models.Employee{Name: "María"}
		luis := models.Employee{Name: "Luis"}
		karla := models.Employee{Name: "Karla"}
		db.Create(&maria)
		db.Create(&luis)
		db.Create(&karla)

		// Relacionar empleados con sucursales (muchos-a-muchos)
		db.Model(&maria).Association("Shops").Append(&shop1)
		db.Model(&luis).Association("Shops").Append(&shop2)
		db.Model(&karla).Association("Shops").Append(&shop3)
	}

	db.Model(&models.Question{}).Count(&count)
	if count == 0 {
		q1 := models.Question{Prompt: "¿Cómo calificarías tu nivel de satisfacción con tu visita reciente a nuestra sucursal?", Type: "radio", OrderNum: 1}
		q2 := models.Question{Prompt: "Considerando tu experiencia en nuestros productos, ¿qué probabilidades hay de que nos recomiendes a un amigo o familiar?", Type: "scale", OrderNum: 2}
		q3 := models.Question{Prompt: "¿Cómo describirías nuestros productos?", Type: "radio", OrderNum: 3}
		q4 := models.Question{Prompt: "Nombre del empleado que atendió su servicio", Type: "text", OrderNum: 4}
		q5 := models.Question{Prompt: "¿Cómo calificarías la amabilidad de nuestro representante de atención al cliente?", Type: "scale", OrderNum: 5}
		db.Create(&[]models.Question{q1, q2, q3, q4, q5})

		var cq1, cq3 models.Question
		db.Where("order_num = ?", 1).First(&cq1)
		db.Where("order_num = ?", 3).First(&cq3)
		db.Create(&[]models.Choice{
			{QuestionID: cq1.ID, Label: "Malo", Value: "malo", OrderNum: 1},
			{QuestionID: cq1.ID, Label: "Regular", Value: "regular", OrderNum: 2},
			{QuestionID: cq1.ID, Label: "Bueno", Value: "bueno", OrderNum: 3},
		})
		db.Create(&[]models.Choice{
			{QuestionID: cq3.ID, Label: "Ni buenos ni malos", Value: "ni_buenos_ni_malos", OrderNum: 1},
			{QuestionID: cq3.ID, Label: "Están bien", Value: "estan_bien", OrderNum: 2},
			{QuestionID: cq3.ID, Label: "Son buenísimos", Value: "son_buenisimos", OrderNum: 3},
			{QuestionID: cq3.ID, Label: "Malos", Value: "malos", OrderNum: 4},
		})
	}
	return nil
}
