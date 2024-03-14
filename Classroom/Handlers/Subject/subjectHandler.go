package Subject

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"webapptrials/Classroom/Models"
)

// Assuming db is your *gorm.DB instance
var db *gorm.DB

func InitializeDB(d *gorm.DB) {
	db = d
}

func GetAllSubjects(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var subjects []Models.Subject
	if result := db.Find(&subjects); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	err := json.NewEncoder(w).Encode(subjects)
	if err != nil {
		fmt.Println(err)
		return
	}
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

func CreateSubject(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var subject Models.Subject
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if result := db.Create(&subject); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	err := json.NewEncoder(w).Encode(subject)
	if err != nil {
		fmt.Println(err)
		return
	}
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
