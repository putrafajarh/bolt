package controller

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/putrafajarh/bolt/internal/infra"
	"github.com/putrafajarh/bolt/internal/middlewares"
	"github.com/putrafajarh/bolt/internal/njir"
	"github.com/putrafajarh/bolt/internal/services"
)

func GetIssue(c *fiber.Ctx) error {
	yesterday := time.Now().Add(-time.Hour * 24)
	dateQuery := c.Query("date", yesterday.Format("2006-01-02"))

	njir, err := njir.NewNjirClient()
	if err != nil {
		return err
	}

	fs, err := infra.NewFirestoreClient()
	if err != nil {
		return err
	}
	defer fs.Close()

	date, err := time.Parse("2006-01-02", dateQuery)
	if err != nil {
		return err
	}

	db, err := middlewares.GetTx(c)
	if err != nil {
		return err
	}
	issues, err := services.NewJiraServiceImpl(services.JiraServiceOptions{
		NjirClient:      njir,
		FirestoreClient: fs,
		DB:              db,
	}).SyncDailyIssues(c.UserContext(), date)
	if err != nil {
		return c.JSON(fiber.Map{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    issues,
	})
}

func Me(c *fiber.Ctx) error {
	client, err := njir.NewClient()
	if err != nil {
		return err
	}
	me, _, err := client.User.GetSelf()
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"me": me,
	})
}
