package migrations

import (
	"hospital-management/models"
)

func (d PostgreDB) MigrateJobGroup() {
	err := d.DB.AutoMigrate(&models.JobGroup{})
	if err != nil {
		return
	}

	defaultJobGroups := []models.JobGroup{
		{Name: "Doktor", Titles: []models.Title{
			{Name: "Asistan"},
			{Name: "Uzman"},
		}},
		{Name: "İdari Personel", Titles: []models.Title{
			{Name: "Başhekim"},
			{Name: "Müdür"},
		}},
		{Name: "Hizmet Personeli", Titles: []models.Title{
			{Name: "Danışman"},
			{Name: "Temizlik Görevlisi"},
			{Name: "Güvenlik Görevlisi"},
		}},
	}
	for _, jobGroup := range defaultJobGroups {
		d.DB.Where(models.JobGroup{Name: jobGroup.Name}).Attrs(jobGroup).FirstOrCreate(&jobGroup)
	}
}
