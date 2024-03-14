package Student

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
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

func GetAllStudents(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var students []Models.Student
	if result := db.Find(&students); result.Error != nil {
		Shared.SendApiResponse[any](w, http.StatusInternalServerError, "Failed to retrieve students", nil, result.Error.Error())
		return
	}

	Shared.SendApiResponse(w, http.StatusOK, "Students retrieved successfully", students, "")
}

func GetStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	var student Models.Student
	if result := db.First(&student, id); result.Error != nil {
		Shared.SendApiResponse[any](w, http.StatusNotFound, "This student doesn't exist", nil, result.Error.Error())
		return
	}

	Shared.SendApiResponse(w, http.StatusOK, "Student retrieved successfully", student, "")
}

func CreateStudent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var student Models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Here you could add manual validation for `student` fields
	// For example, checking if `student.Name` is not empty since it's mandatory
	if student.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if result := db.Create(&student); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	err := json.NewEncoder(w).Encode(student)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func UpdateStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	var student Models.Student
	if err := db.First(&student, id).Error; err != nil {
		http.NotFound(w, r)
		return
	}

	var updates Models.Student
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result := db.Model(&student).Updates(updates); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	err := json.NewEncoder(w).Encode(student)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func PatchStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	u64, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	studentID := uint(u64)

	// Decode the request body for partial updates
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch the existing student
	var student Models.Student
	if result := db.First(&student, studentID); result.Error != nil {
		http.NotFound(w, r)
		return
	}

	// Perform the update using the map
	if result := db.Model(&student).Updates(updates); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func DeleteStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")

	if result := db.Delete(&Models.Student{}, id); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	_, err := fmt.Fprintf(w, "Student Deleted")
	if err != nil {
		fmt.Println(err)
		return
	}
}
