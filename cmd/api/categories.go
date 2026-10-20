package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) getAllCategories(c *gin.Context) {
	tenantID := c.GetInt("tenantID")
	categories, err := app.models.Categories.GetAll(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve categories",
		})
		return
	}

	c.JSON(http.StatusOK, categories)
}
