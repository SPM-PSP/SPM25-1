package handler

import (
	"UnoBackend/DB"
	"UnoBackend/internal/model/suop"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllSuops(c *gin.Context) {
	var suops []suop.Suop
	if err := DB.DB.Find(&suops).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, suops)
}

func CreateSuop(c *gin.Context) {
	var newSuop suop.Suop
	if err := c.ShouldBindJSON(&newSuop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := DB.DB.Create(&newSuop).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, newSuop)
}

func UpdateSuop(c *gin.Context) {
	id := c.Param("id")
	var suopData suop.Suop
	if err := DB.DB.First(&suopData, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录未找到"})
		return
	}

	if err := c.ShouldBindJSON(&suopData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.DB.Save(&suopData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, suopData)
}

// 删除 Suop
func DeleteSuop(c *gin.Context) {
	id := c.Param("id")
	if err := DB.DB.Delete(&suop.Suop{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
