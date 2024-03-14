package Test_HttpRouter

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"` //You can do something like Name string `json:"name,omitempty"
	Age  int    `json:"age"`
}

//In situations where you're dealing with multiple versions of an API or trying to maintain backward compatibility,
//omitempty allows you to easily introduce new fields without breaking existing clients.
//Clients that are unaware of the new fields will simply ignore them if they're not present in the response.

var studentsDB = make(map[string]Student)
var lastStudentID int

func Welcome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := fmt.Fprint(w, "Welcome to the Students API!\n")
	if err != nil {
		return
	}
}

func GetAllStudents(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var students []Student
	for _, student := range studentsDB {
		students = append(students, student)
	}
	err := json.NewEncoder(w).Encode(students)
	if err != nil {
		return
	}
}

func GetStudentByID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	student, ok := studentsDB[id]
	if !ok {
		http.NotFound(w, r)
		return
	}
	err := json.NewEncoder(w).Encode(student)
	if err != nil {
		return
	}
}

func CreateStudent(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var student Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lastStudentID++
	student.ID = lastStudentID
	studentsDB[strconv.Itoa(student.ID)] = student

	w.WriteHeader(http.StatusCreated)
	err := json.NewEncoder(w).Encode(student)
	if err != nil {
		return
	}
}

func UpdateStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if _, ok := studentsDB[id]; !ok {
		http.NotFound(w, r)
		return
	}

	var studentUpdate Student
	if err := json.NewDecoder(r.Body).Decode(&studentUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	studentUpdate.ID, _ = strconv.Atoi(id) // Keep the ID from the URL
	studentsDB[id] = studentUpdate
	err := json.NewEncoder(w).Encode(studentUpdate)
	if err != nil {
		return
	}
}

func PatchStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	student, ok := studentsDB[id]
	if !ok {
		http.NotFound(w, r)
		return
	}

	// Apply updates to the student
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if name, ok := updates["name"]; ok {
		student.Name = name.(string)
	}
	if age, ok := updates["age"]; ok {
		student.Age = int(age.(float64)) // JSON numbers are float64
	}

	studentsDB[id] = student
	err := json.NewEncoder(w).Encode(student)
	if err != nil {
		return
	}
}

func DeleteStudent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if _, ok := studentsDB[id]; !ok {
		http.NotFound(w, r)
		return
	}

	delete(studentsDB, id)
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "Student with ID %s deleted\n", id)
	if err != nil {
		return
	}
}
