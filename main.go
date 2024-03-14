package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"webapptrials/Test_HttpRouter"
)

func main() {
	router := httprouter.New()
	router.GET("/", Test_HttpRouter.Welcome)
	router.GET("/students", Test_HttpRouter.GetAllStudents)
	router.GET("/students/:id", Test_HttpRouter.GetStudentByID)
	router.POST("/students", Test_HttpRouter.CreateStudent)
	router.PUT("/students/:id", Test_HttpRouter.UpdateStudent)
	router.PATCH("/students/:id", Test_HttpRouter.PatchStudent)
	router.DELETE("/students/:id", Test_HttpRouter.DeleteStudent)

	fmt.Println("Server is running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err)
	}
}
