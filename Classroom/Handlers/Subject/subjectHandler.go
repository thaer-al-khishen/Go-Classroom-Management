package Subject

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"webapptrials/Classroom/Models"
	"webapptrials/Classroom/Shared"
)

// Assuming db is your *gorm.DB instance
var db *gorm.DB

func InitializeDB(d *gorm.DB) {
	db = d
}

func GetAllSubjects(c *gin.Context) {
	// Helper function to get the first non-empty query param value
	//getQueryParam := func(keys ...string) string {
	//	for _, key := range keys {
	//		if value, exists := c.GetQuery(key); exists && value != "" {
	//			return value
	//		}
	//	}
	//	return ""
	//}

	// Default values
	defaultPage := 1
	defaultPageSize := 10

	// Try to get 'page' parameter (ignore case)
	pageStr := Shared.GetQueryParam(c, "page", "Page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = defaultPage // Use default if not specified or invalid
	}

	// Try to get 'pageSize' parameter (ignore case)
	pageSizeStr := Shared.GetQueryParam(c, "pageSize", "Pagesize", "pagesize", "PageSize")
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		pageSize = defaultPageSize // Use default if not specified or invalid
	}

	nameFilter := Shared.GetQueryParam(c, "name", "Name")

	// Apply pagination and filtering
	var subjects []Models.Subject
	query := db.Offset((page - 1) * pageSize).Limit(pageSize)

	if nameFilter != "" {
		query = query.Where("name LIKE ?", "%"+nameFilter+"%")
	}

	if result := query.Find(&subjects); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Calculate total number of records for pagination metadata
	var totalRecords int64
	db.Model(&Models.Subject{}).Where("name LIKE ?", "%"+nameFilter+"%").Count(&totalRecords)

	// Return paginated and filtered results along with pagination metadata
	c.JSON(http.StatusOK, gin.H{
		"data":          subjects,
		"total_records": totalRecords,
		"page":          page,
		"page_size":     pageSize,
		"total_pages":   (totalRecords + int64(pageSize) - 1) / int64(pageSize), // Ceiling division
	})
}

func GetSubject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var subject Models.Subject
	if result := db.First(&subject, id); result.Error != nil {
		http.NotFound(w, r)
		return
	}

	err := json.NewEncoder(w).Encode(subject)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func CreateSubject(c *gin.Context) {
	var subject Models.Subject
	if err := c.ShouldBindJSON(&subject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := db.Create(&subject); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, subject)
}

func UpdateSubject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var subject Models.Subject
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//subject.ID = uint(id) // Convert id to uint and assign it to the subject
	u64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	subject.ID = uint(u64)

	if result := db.Save(&subject); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	jsonError := json.NewEncoder(w).Encode(subject)
	if jsonError != nil {
		fmt.Println(jsonError)
		return
	}
}

func PatchSubject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	u64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	subjectID := uint(u64)

	// Fetch the existing subject
	var subject Models.Subject
	if result := db.First(&subject, subjectID); result.Error != nil {
		http.NotFound(w, r)
		return
	}

	// Decode the request body for partial updates
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform the update using the map
	if result := db.Model(&subject).Updates(updates); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(subject)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func DeleteSubject(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var subject Models.Subject
	if result := db.Delete(&subject, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	_, err := fmt.Fprintf(w, "Subject Deleted")
	if err != nil {
		fmt.Println(err)
		return
	}
}
