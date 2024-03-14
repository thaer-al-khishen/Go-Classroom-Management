package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"webapptrials/Classroom/Auth/Handlers"
	"webapptrials/Classroom/Auth/Middleware"
	Models2 "webapptrials/Classroom/Auth/Models"
	"webapptrials/Classroom/Handlers/Teacher"
	"webapptrials/Classroom/Models"
	"webapptrials/Classroom/Test"
)

func main() {
	dsn := "host=localhost user=postgres dbname=classroom_management sslmode=disable password=Admin@123"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Database connected!")

	err = db.AutoMigrate(&Models.Student{}, &Models.Teacher{}, &Models.Classroom{}, &Models.Subject{}, &Models2.User{}, &Models2.RefreshTokenModel{})
	if err != nil {
		panic("failed to migrate database")
	}

	Handlers.InitializeDB(db) // Pass the DB instance to Auth
	//Student.InitializeDB(db)   // Pass the DB instance to Student
	Teacher.InitializeDB(db) // Pass the DB instance to Teacher
	//Classroom.InitializeDB(db) // Pass the DB instance to Classroom
	//Subject.InitializeDB(db)   // Pass the DB instance to Subject

	//
	//router := httprouter.New()
	//
	////Auth
	//router.POST("/login", Handlers.Login)
	//router.POST("/register", Handlers.Register)
	//router.POST("/refresh", Handlers.RefreshToken)
	//
	////Student
	//router.GET("/students", Student.GetAllStudents)
	//router.GET("/students/:id", Student.GetStudent)
	//router.POST("/student", Student.CreateStudent)
	//router.PUT("/students/:id", Student.UpdateStudent)
	//router.PATCH("/students/:id", Student.PatchStudent)
	//router.DELETE("/students/:id", Student.DeleteStudent)
	//
	////Teacher
	//router.GET("/teachers", Teacher.GetAllTeachers)
	//router.GET("/teachers/:id", Teacher.GetTeacher)
	//router.POST("/teacher", Teacher.CreateTeacher)
	//router.PUT("/teachers/:id", Teacher.UpdateTeacher)
	//router.PATCH("/teachers/:id", Teacher.PatchTeacher)
	//router.DELETE("/teachers/:id", Teacher.DeleteTeacher)
	//
	////Subject
	//router.GET("/subjects/:id", Subject.GetSubject)
	//router.GET("/subjects", Subject.GetAllSubjects)
	//router.POST("/subject", Subject.CreateSubject)
	//router.PUT("/subjects/:id", Subject.UpdateSubject)
	//router.PATCH("/subjects/:id", Subject.PatchSubject)
	//router.DELETE("/subjects/:id", Subject.DeleteSubject)
	//
	////Classroom
	//router.GET("/classrooms", Classroom.GetAllClassrooms)
	//router.GET("/classrooms/:id", Classroom.GetClassroom)
	//router.POST("/classroom", Classroom.CreateClassroom)
	//router.PUT("/classrooms/:id", Classroom.UpdateClassroom)
	//router.PATCH("/classrooms/:id", Classroom.PatchClassroom)
	//router.DELETE("/classrooms/:id", Classroom.DeleteClassroom)
	//
	////Test
	//router.POST("/test", Middleware.Authenticate(Test.TestData))
	//
	//fmt.Println("Server started on port 8081")
	//httpConnectError := http.ListenAndServe(":8081", router)
	//if httpConnectError != nil {
	//	fmt.Println(httpConnectError)
	//}

	r := gin.Default()

	//Auth
	r.POST("/login", Handlers.Login)
	r.POST("/register", Handlers.Register)
	r.POST("/refresh", Handlers.RefreshToken)

	//Teachers
	r.GET("/teachers", Teacher.GetAllTeachers)
	r.GET("/teachers/:id", Teacher.GetTeacher)
	r.POST("/teachers", Teacher.CreateTeacher)
	r.PUT("/teachers/:id", Teacher.UpdateTeacher)
	r.PATCH("/teachers/:id", Teacher.PatchTeacher)
	r.DELETE("/teachers/:id", Teacher.DeleteTeacher)

	//Test
	r.POST("/test", Middleware.RoleCheckMiddleware(Models2.Admin), Test.TestData)

	ginError := r.Run(":8081")
	if ginError != nil {
		fmt.Println(ginError)
		return
	} // By default, it listens on :8080

}
