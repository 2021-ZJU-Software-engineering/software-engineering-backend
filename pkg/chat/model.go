package chat

// type Chat struct {
// 	ID        uint   `gorm:"primaryKey"`
// 	URL       string //The url of the chatting window
// 	CaseID    uint
// 	DoctorID  uint
// 	PatientID uint
// }

/*
Questions
-------------------------------------------------------------------
According to the department of the doctor,we can find a set of questions.
*/
type Category struct {
	ID           uint   `gorm:"primaryKey"`
	DepartmentID uint   `validate:"required"`
	Questions    string `validate:"required"`
}