package Classroom

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http"
	"time"
	"webapptrials/Classroom/Models"
)

var db *gorm.DB

func InitializeDB(d *gorm.DB) {
	db = d
}

func GetAllClassrooms(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var classrooms []Models.Classroom
	var response []ClassroomResponse

	if result := db.Preload("StudentIDs").Preload("TeacherIDs").Find(&classrooms); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	for _, classroom := range classrooms {
		cr := ClassroomResponse{
			ID:        classroom.ID,
			CreatedAt: classroom.CreatedAt,
			UpdatedAt: classroom.UpdatedAt,
			// Populate StudentIDs and TeacherIDs from the classroom associations
		}
		for _, student := range classroom.StudentIDs {
			cr.StudentIDs = append(cr.StudentIDs, student.ID)
		}
		for _, teacher := range classroom.TeacherIDs {
			cr.TeacherIDs = append(cr.TeacherIDs, teacher.ID)
		}
		response = append(response, cr)
	}

	json.NewEncoder(w).Encode(response)
}

func GetClassroom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	var classroom Models.Classroom
	var response ClassroomResponse

	if result := db.Preload("StudentIDs").Preload("TeacherIDs").First(&classroom, id); result.Error != nil {
		http.NotFound(w, r)
		return
	}

	response.ID = classroom.ID
	response.CreatedAt = classroom.CreatedAt
	response.UpdatedAt = classroom.UpdatedAt
	// Populate StudentIDs and TeacherIDs
	for _, student := range classroom.StudentIDs {
		response.StudentIDs = append(response.StudentIDs, student.ID)
	}
	for _, teacher := range classroom.TeacherIDs {
		response.TeacherIDs = append(response.TeacherIDs, teacher.ID)
	}

	json.NewEncoder(w).Encode(response)
}

func CreateClassroom(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var classroom Models.Classroom
	if err := json.NewDecoder(r.Body).Decode(&classroom); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result := db.Create(&classroom); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(classroom)
}

func UpdateClassroom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Perform the update
	if result := db.Model(&Models.Classroom{}).Where("id = ?", id).Updates(updates); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the updated classroom
	var updatedClassroom Models.Classroom
	if result := db.Preload("StudentIDs").Preload("TeacherIDs").First(&updatedClassroom, id); result.Error != nil {
		http.NotFound(w, r)
		return
	}

	// Construct the response
	var response ClassroomResponse
	response.ID = updatedClassroom.ID
	response.CreatedAt = updatedClassroom.CreatedAt
	response.UpdatedAt = updatedClassroom.UpdatedAt
	for _, student := range updatedClassroom.StudentIDs {
		response.StudentIDs = append(response.StudentIDs, student.ID)
	}
	for _, teacher := range updatedClassroom.TeacherIDs {
		response.TeacherIDs = append(response.TeacherIDs, teacher.ID)
	}

	// Return the updated classroom
	json.NewEncoder(w).Encode(response)
}

func PatchClassroom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	var input struct {
		StudentIDs []uint `json:"StudentIDs"`
		// Add other fields here if needed
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the classroom to update its associations
	var classroom Models.Classroom
	if err := db.First(&classroom, id).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	// If StudentIDs are provided, update the association
	if len(input.StudentIDs) > 0 {
		var students []Models.Student
		if err := db.Find(&students, input.StudentIDs).Error; err != nil {
			http.Error(w, "Failed to find students", http.StatusBadRequest)
			return
		}
		// Replace the classroom's students with the provided list
		if err := db.Model(&classroom).Association("StudentIDs").Replace(&students); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Optionally, handle other fields and update them here if needed

	// Respond with the updated classroom information
	db.Preload("StudentIDs").First(&classroom, id)
	json.NewEncoder(w).Encode(classroom)
}

func DeleteClassroom(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if result := db.Delete(&Models.Classroom{}, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Classroom Deleted")
}

type ClassroomResponse struct {
	ID         uint      `json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	StudentIDs []uint    `json:"studentIds"`
	TeacherIDs []uint    `json:"teacherIds"`
}
