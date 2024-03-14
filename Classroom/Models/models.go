package Models

import (
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	Name        string `json:"studentname"` //It would be saved as Name in database, and required as studentname from the client
	Age         int
	ClassroomID uint
}

type Teacher struct {
	gorm.Model
	Name      string
	SubjectID uint
	Subject   Subject `gorm:"foreignKey:SubjectID"`
}

type Classroom struct {
	gorm.Model
	StudentIDs []Student `gorm:"foreignKey:ClassroomID"`
	TeacherIDs []Teacher `gorm:"many2many:classroom_teachers;"`
}

type Subject struct {
	gorm.Model
	Name      string `gorm:"uniqueIndex"`
	TeacherID uint   `gorm:"unique"`
}
