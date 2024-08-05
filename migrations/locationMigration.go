package migrations

import (
	"hospital-management/models"
)

func (d PostgreDB) MigrateLocations() {
	// Modelleri migrate et
	d.DB.AutoMigrate(&models.City{}, &models.District{})

	// Varsayılan iller ve ilçeler
	defaultCities := []models.City{
		{Name: "Adana", Plate: 1, Districts: []models.District{
			{Name: "Seyhan"},
			{Name: "Yüreğir"},
			{Name: "Çukurova"},
			{Name: "Sarıçam"},
		}},
		{Name: "Adıyaman", Plate: 2, Districts: []models.District{
			{Name: "Merkez"},
			{Name: "Kahta"},
			{Name: "Samsat"},
			{Name: "Besni"},
		}},
		{Name: "Afyonkarahisar", Plate: 3, Districts: []models.District{
			{Name: "Merkez"},
			{Name: "Bolvadin"},
			{Name: "Dinar"},
			{Name: "Emirdağ"},
		}},
		{Name: "Ağrı", Plate: 4, Districts: []models.District{
			{Name: "Merkez"},
			{Name: "Doğubayazıt"},
			{Name: "Patnos"},
			{Name: "Eleşkirt"},
		}},
		{Name: "Amasya", Plate: 5, Districts: []models.District{
			{Name: "Merkez"},
			{Name: "Merzifon"},
			{Name: "Suluova"},
			{Name: "Taşova"},
		}},
		{Name: "Ankara", Plate: 6, Districts: []models.District{
			{Name: "Çankaya"},
			{Name: "Keçiören"},
			{Name: "Mamak"},
			{Name: "Sincan"},
		}},
	}

	// Varsayılan verileri ekleyin
	for _, city := range defaultCities {
		d.DB.Where(models.City{Name: city.Name}).Attrs(city).FirstOrCreate(&city)
	}
}
