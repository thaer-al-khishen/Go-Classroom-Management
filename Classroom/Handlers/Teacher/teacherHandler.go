package Teacher

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"webapptrials/Classroom/Models"
	"webapptrials/Classroom/Shared"
)

var db *gorm.DB

func InitializeDB(d *gorm.DB) {
	db = d
}

func GetAllTeachers(c *gin.Context) {
	var teachers []Models.Teacher
	if result := db.Preload("Subject").Find(&teachers); result.Error != nil {
		Shared.SendGinGenericApiResponse(c, http.StatusInternalServerError, "Can't find teachers", "", result.Error.Error())
		return
	}

	Shared.SendGinGenericApiResponse[any](c, http.StatusOK, "Teachers retrieved", teachers, "")
}

func GetTeacher(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var teacher Models.Teacher
	if result := db.Preload("Subject").First(&teacher, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Teacher not found"})
		return
	}

	c.JSON(http.StatusOK, teacher)
}

func CreateTeacher(c *gin.Context) {
	var teacher Models.Teacher
	if err := c.ShouldBindJSON(&teacher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := db.Create(&teacher); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	db.Preload("Subject").First(&teacher, teacher.ID)
	c.JSON(http.StatusCreated, teacher)
}

func UpdateTeacher(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := db.Model(&Models.Teacher{}).Where("id = ?", id).Updates(updates); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Teacher updated successfully"})
}

func PatchTeacher(c *gin.Context) {
	id := c.Param("id")
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := db.Model(&Models.Teacher{}).Where("id = ?", id).Updates(updates); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var teacher Models.Teacher
	db.Preload("Subject").First(&teacher, id)
	c.JSON(http.StatusOK, teacher)
}

func DeleteTeacher(c *gin.Context) {
	id := c.Param("id")
	if result := db.Delete(&Models.Teacher{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Teacher deleted"})
}
