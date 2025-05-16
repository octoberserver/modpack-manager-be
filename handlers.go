package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetModpacks(c *gin.Context) {
	var modpacks []Modpack
	db.Find(&modpacks)
	c.JSON(http.StatusOK, modpacks)
}

func CreateModpack(c *gin.Context) {
	var modpack Modpack
	if err := c.BindJSON(&modpack); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count := int64(0)
	result := db.Model(&Modpack{}).Where("id = ?", modpack.ID).Count(&count)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Modpack with this ID already exists"})
		return
	}

	result = db.Create(&modpack)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, modpack)
}

func GetModpack(c *gin.Context) {
	id := c.Param("id")
	var modpack Modpack
	if result := db.First(&modpack, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Modpack not found"})
		return
	}
	c.JSON(http.StatusOK, modpack)
}

func UpdateModpack(c *gin.Context) {
	id := c.Param("id")
	var modpack Modpack
	if result := db.First(&modpack, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Modpack not found"})
		return
	}

	if err := c.BindJSON(&modpack); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&modpack)
	c.JSON(http.StatusOK, modpack)
}

func DeleteModpack(c *gin.Context) {
	id := c.Param("id")
	db.Delete(&Modpack{}, "id = ?", id)
	c.Status(http.StatusNoContent)
}

func GetLatestModpack(c *gin.Context) {
	server := c.Param("server")
	var latest LatestModpack
	if result := db.Preload("Modpack").First(&latest, "server = ?", server); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Server not found"})
		return
	}
	c.JSON(http.StatusOK, latest.Modpack)
}

func GetLatestModpacks(c *gin.Context) {
}

func SetLatestModpack(c *gin.Context) {
	server := c.Param("server")
	modpackID := c.Param("modpack_id")

	// Verify modpack exists
	var modpack Modpack
	if result := db.First(&modpack, "id = ?", modpackID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Modpack not found"})
		return
	}

	// Create or update latest modpack entry
	latest := LatestModpack{
		Server:    server,
		ModpackID: modpackID,
	}

	// Upsert operation (Update if exists, else create)
	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "server"}},
		DoUpdates: clause.AssignmentColumns([]string{"modpack_id"}),
	}).Create(&latest)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update latest modpack"})
		return
	}

	data := getWebhookJSON(modpack.Name, modpack.MCVersion, modpack.URL, modpack.Links.CF)
	_, err := http.Post("https://discord.com/api/webhooks/1342488647936643082/kFZngwgGKR2YYCI53QcZKV--_GjQbkKCBuzsD46fBxDhaZEp6z_Bu3H9IcHXOCG9O70V", "application/json", strings.NewReader(data))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send webhook"})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Latest modpack updated successfully",
		"server":  server,
		"modpack": modpackID,
	})
}

func DeleteLatestModpack(c *gin.Context) {
}
