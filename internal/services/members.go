package services

import (
	"time"

	"github.com/letitloose/nsdtr-club-us/internal/models"
	"github.com/letitloose/nsdtr-club-us/internal/validator"
)

type MemberForm struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Website     string
	Region      int
	JoinedDate  string
	validator.Validator
}

type MemberService struct {
	*models.MemberModel
	*Email
}

func (service *MemberService) CreateMember(mf *MemberForm) (int, error) {

	//validate
	mf.CheckField(validator.NotBlank(mf.FirstName), "firstname", "You must enter a first name.")
	mf.CheckField(validator.NotBlank(mf.LastName), "lastname", "You must enter a last name.")
	mf.CheckField(validator.ValidEmail(mf.Email), "email", "You must enter a valid email: name@domain.ext")

	if !mf.Valid() {
		return 0, models.ErrBadData
	}

	joined, err := time.Parse("2006-01-02", mf.JoinedDate)
	if err != nil {
		return 0, err
	}

	return service.Insert(mf.FirstName, mf.LastName, mf.PhoneNumber, mf.Email, mf.Website, mf.Region, joined)
}
