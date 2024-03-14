1- CRUD Operations:
Check the handler.go files

2- Input validation:
Check GetAllStudents inside studentHandler

3- Changing json field names:
type Student struct {
	gorm.Model
	Name        string `json:"studentname"` //It would be saved as Name in database, and required as name from the client
	Age         int
	ClassroomID uint
}

4- Generic return types:
Check ApiResponse
And GetAllStudents inside studentHandler

5- DTO Return type:
ClassroomResponse inside classroomHandler
