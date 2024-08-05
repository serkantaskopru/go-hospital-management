package app

import (
	"hospital-management/controllers"
	"hospital-management/middleware"
	"hospital-management/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(c *fiber.App, db *gorm.DB) {
	locationController := &controllers.LocationController{DB: db}
	polyclinicController := &controllers.PolyclinicController{DB: db}
	personController := &controllers.PersonController{DB: db}
	jobGroupController := &controllers.JobGroupController{DB: db}
	titleController := &controllers.TitleController{DB: db}
	hospitalClinicController := &controllers.HospitalClinicController{DB: db}
	staffController := &controllers.StaffController{DB: db}
	c.Get("/cities", locationController.GetCities)
	c.Get("/cities/:cityId/districts", locationController.GetDistrictsByCity)
	api := c.Group("/api", middleware.AuthMiddleware)
	{
		api.Get("/protected", func(c *fiber.Ctx) error {
			user := c.Locals("user").(models.User)

			return c.JSON(fiber.Map{"user": user})
		})
		api.Get("/jobgroups", jobGroupController.GetJobGroups)
		api.Get("/titles/:jobGroupId", titleController.GetTitles)

		api.Get("/staff", staffController.ListStaff)
		api.Post("/staff", middleware.AuthorizedMiddleware, staffController.CreateStaff)
		api.Delete("/staff/:id", middleware.AuthorizedMiddleware, staffController.DeleteStaff)
		api.Put("/staff/:id", middleware.AuthorizedMiddleware, staffController.UpdateStaff)

		api.Get("/persons", personController.ListPersons)
		api.Post("/persons", middleware.AuthorizedMiddleware, personController.CreatePerson)
		api.Delete("/persons/:id", middleware.AuthorizedMiddleware, personController.DeletePerson)
		api.Put("/persons/:id", middleware.AuthorizedMiddleware, personController.UpdatePerson)

		api.Get("/polyclinics", polyclinicController.GetPolyclinics)
		api.Get("/hospitalclinics", hospitalClinicController.GetHospitalClinics)
		api.Post("/hospitalclinics/create", middleware.AuthorizedMiddleware, hospitalClinicController.CreateHospitalClinic)
	}
}
