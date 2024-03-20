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
	FirstName      string
	LastName       string
	JointFirstName string
	JointLastName  string
	PhoneNumber    string
	Email          string
	Website        string
	Region         int
	JoinedDate     string
	validator.Validator
}

type MembershipForm struct {
	ID               int
	MemberID         int
	Year             int
	MembershipType   string
	MembershipAmount float32
	RosterAmount     float32
	HealthDonations  float32
	RescueDonations  float32
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

	return service.Insert(mf.FirstName, mf.LastName, mf.JointFirstName, mf.JointLastName, mf.PhoneNumber, mf.Email, mf.Website, mf.Region, joined)
}

func (service *MemberService) MigrateLegacyMembers() error {

	legacyMembers, err := service.Legacy.List()
	if err != nil {
		return err
	}

	for _, member := range legacyMembers {
		memberID, err := service.Insert(member.FirstName, member.LastName, member.JointFirstName.String, member.JointLastName.String, member.PhoneNumber.String, member.Email.String, member.Website.String, member.Region, member.JoinedDate.Time)
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

		err = service.AddMemberships(member, memberID)
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

func (service *MemberService) AddMembership(membershipForm *MembershipForm) error {

	mm := models.MembershipModel{DB: service.DB}
	newMembership := &models.Membership{MemberID: membershipForm.MemberID, Year: membershipForm.Year}
	membershipID, err := mm.Insert(newMembership)
	if err != nil {
		return err
	}

	typeItem := &models.MembershipItem{MembershipID: membershipID, ItemCode: membershipForm.MembershipType, AmountPaid: membershipForm.MembershipAmount}
	_, err = mm.InsertMembershipItem(typeItem)
	if err != nil {
		return err
	}

	//roster
	if membershipForm.RosterAmount != 0.0 {
		rosterItem := &models.MembershipItem{MembershipID: membershipID, ItemCode: "PR", AmountPaid: membershipForm.RosterAmount}
		_, err = mm.InsertMembershipItem(rosterItem)
		if err != nil {
			return err
		}
	}

	//health and genetics donation
	if membershipForm.HealthDonations != 0.0 {
		healthItem := &models.MembershipItem{MembershipID: membershipID, ItemCode: "HG", AmountPaid: membershipForm.HealthDonations}
		_, err = mm.InsertMembershipItem(healthItem)
		if err != nil {
			return err
		}
	}

	//rescue donation
	if membershipForm.RescueDonations != 0.0 {
		rescueItem := &models.MembershipItem{MembershipID: membershipID, ItemCode: "RS", AmountPaid: membershipForm.RescueDonations}
		_, err = mm.InsertMembershipItem(rescueItem)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *MemberService) AddMemberships(legacyMember *models.LegacyMember, memberID int) error {

	legacyMemberships, err := service.Legacy.GetMemberships(legacyMember.Member.ID)
	if err != nil {
		return err
	}

	mm := models.MembershipModel{DB: service.DB}
	dueSchedule, err := mm.GetDueSchedule()
	if err != nil {
		return err
	}

	for _, membership := range legacyMemberships {
		newMembership := &models.Membership{MemberID: memberID, Year: membership.Year}
		membershipID, err := mm.Insert(newMembership)
		if err != nil {
			return err
		}

		membershipTypeCode := lookupItem(dueSchedule.Items , legacyMember.MembershipType)
		if legacyMember.CountryCode.String != "USA" {
			if legacyMember.CountryCode.String == "CAN" {
				membershipTypeCode = "CM"
			} else {
				membershipTypeCode = "FM"
			}
		}

		typeItem := &models.MembershipItem{MembershipID: membershipID, ItemCode: membershipTypeCode, AmountPaid: membership.MembershipAmount}
		_, err = mm.InsertMembershipItem(typeItem)
		if err != nil {
			return err
		}

		//roster
		if membership.RosterAmount != 0.0 {
			rosterItem := &models.MembershipItem{MembershipID: membershipID, ItemCode: "PR", AmountPaid: membership.RosterAmount}
			_, err = mm.InsertMembershipItem(rosterItem)
			if err != nil {
				return err
			}
		}

		//health and genetics donation
		if membership.HealthDonations != 0.0 {
			healthItem := &models.MembershipItem{MembershipID: membershipID, ItemCode: "HG", AmountPaid: membership.HealthDonations}
			_, err = mm.InsertMembershipItem(healthItem)
			if err != nil {
				return err
			}
		}

		//rescue donation
		if membership.RescueDonations != 0.0 {
			rescueItem := &models.MembershipItem{MembershipID: membershipID, ItemCode: "RS", AmountPaid: membership.RescueDonations}
			_, err = mm.InsertMembershipItem(rescueItem)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func lookupItem(items []*models.DueScheduleItem, display string) string {
	for _, item := range items {
		if item.Display == display {
			return item.Code
		}
	}

	return ""
}
