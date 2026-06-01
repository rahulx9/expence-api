package main

import (
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine, store *Store) {
	r.POST("/data", func(c *gin.Context) {
		var item Item
		if err := c.ShouldBindJSON(&item); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		if err := store.AddItem(item); err != nil {
			if err == ErrDuplicateID {
				c.JSON(400, gin.H{"error": "Item with this ID already exists"})
			} else {
				c.JSON(500, gin.H{"error": "Failed to add item"})
			}
			return
		}

		if err := store.Save(); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save item"})
			return
		}
	})

	r.GET("/data", func(c *gin.Context) {
		items := store.GetAll()
		c.JSON(200, items)
	})

	r.DELETE("/data/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(400, gin.H{"error": "ID is required"})
			return
		}

		if !store.DeleteItem(id) {
			c.JSON(404, gin.H{"error": "Item not found"})
			return
		}

		c.Status(204)
	})
}
