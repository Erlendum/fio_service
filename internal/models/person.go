package models

type PersonField int
type PersonFieldsToUpdate map[PersonField]any

const (
	UserFieldLogin = PersonField(iota)
	UserFieldFio
	UserFieldDateBirth
	UserFieldGender
	IsAdmin
)

type PersonGender string

const (
	MaleUserGender   = PersonGender("Male")
	FemaleUserGender = PersonGender("Female")
)

type Person struct {
	Id         uint64
	Name       string
	Surname    string
	Patronymic string
	Age        uint64
	Gender     PersonGender
}
