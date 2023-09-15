package models

type PersonField int
type PersonFieldsToUpdate map[PersonField]any

const (
	PersonFieldName = PersonField(iota)
	PersonFieldSurname
	PersonFieldPatronymic
	PersonFieldAge
	PersonFieldGender
	PersonFieldNationality
)

type PersonGender string

const (
	MaleUserGender   = PersonGender("Male")
	FemaleUserGender = PersonGender("Female")
)

type Person struct {
	Id          uint64
	Name        string
	Surname     string
	Patronymic  string
	Age         uint64
	Gender      PersonGender
	Nationality string
}
