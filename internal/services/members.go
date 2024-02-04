package services

import (
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
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
	Legacy *models.LegacyModel
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

func (service *MemberService) MigrateLegacyMembers() error {

	legacyMembers, err := service.Legacy.List()
	if err != nil {
		return err
	}

	for _, member := range legacyMembers {
		memberID, err := service.Insert(member.FirstName, member.LastName, member.PhoneNumber.String, member.Email.String, member.Website.String, member.Region, member.JoinedDate.Time)
		if err != nil {
			var mySQLError *mysql.MySQLError
			if errors.As(err, &mySQLError) {
				if mySQLError.Number == 1062 {
					continue
				} else {
					return err
				}
			}
		}

		//massage countries
		if member.Address.CountryCode.String == "CANADA" {
			member.Address.CountryCode.String = "CAN"
		}
		if strings.ToUpper(member.Address.CountryCode.String) == "GERMANY" {
			member.Address.CountryCode.String = "GER"
		}
		if member.Address.CountryCode.String == "UNITED KINGDOM" {
			member.Address.CountryCode.String = "UK"
		}
		if member.Address.CountryCode.String == "AUSTRALIA" {
			member.Address.CountryCode.String = "AUS"
		}
		if member.Address.CountryCode.String == "SWEDEN" {
			member.Address.CountryCode.String = "SWE"
		}
		if member.Address.CountryCode.String == "SWITZERLAND" {
			member.Address.CountryCode.String = "SWI"
		}
		if member.Address.CountryCode.String == "FINLAND" {
			member.Address.CountryCode.String = "FIN"
		}
		if strings.ToUpper(member.Address.CountryCode.String) == "NEW ZEALAND" {
			member.Address.CountryCode.String = "NZ"
		}
		err = service.AddMemberAddress(memberID, member.Address)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *MemberService) AddMemberAddress(memberID int, address *models.Address) error {

	am := &models.AddressModel{DB: service.DB}
	addressID, err := am.Insert(address)
	if err != nil {
		return err
	}

	return service.AddAddress(memberID, addressID)
}
