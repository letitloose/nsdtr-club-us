package models

import (
	"database/sql"
)

type LegacyModel struct {
	DB *sql.DB
}

type LegacyMember struct {
	*Member
	*Address
}

type LegacyMembership struct {
	MemberID         int
	Year             int
	CheckDate        sql.NullTime
	CheckNumber      int
	AmountPaid       float32
	MembershipAmount float32
	RosterAmount     float32
	HealthDonations  float32
	RescueDonations  float32
	DateReceived     sql.NullTime
	DateProcessed    sql.NullTime
	Notes            sql.NullString
}

func (m *LegacyModel) List() ([]*LegacyMember, error) {

	stmt := `select  FirstName1, LastName1, FirstName2, LastName2, HomePhone, EmailName, ` + "`Web Page`" + `, Region, DateJoined1, 
		HomeAddress1, HomeAddress2, HomeCity, HomeStateOrProvince, HomePostalCode, HomeCountry from Members`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []*LegacyMember{}

	for rows.Next() {
		member := &LegacyMember{Member: &Member{}, Address: &Address{}}

		err := rows.Scan(&member.FirstName, &member.LastName, &member.JointFirstName, &member.JointLastName, &member.PhoneNumber,
			&member.Email, &member.Website, &member.Region, &member.JoinedDate, &member.Address1,
			&member.Address2, &member.City, &member.StateProvince, &member.ZipCode, &member.CountryCode)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

func (m *LegacyModel) GetMemberships(memberID int) ([]*LegacyMembership, error) {

	stmt := "select MemberId, `Membership Year` , `Check Date` , `Check Number` , `Amount Paid` , `Amount due` , `Printed Roster Paid` , `Health and Genetics Amount` , `Rescue Amount` , " +
		`DateReceived , DateProcessed , NotestoTreasurer 
		from Dues 
		where MemberID = ?
		order by ` + "`Membership Year`;"

	rows, err := m.DB.Query(stmt, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberships := []*LegacyMembership{}

	for rows.Next() {
		membership := &LegacyMembership{}

		err := rows.Scan(&membership.MemberID, &membership.Year, &membership.CheckDate,
			&membership.CheckNumber, &membership.AmountPaid, &membership.MembershipAmount, &membership.RosterAmount, &membership.HealthDonations,
			&membership.RescueDonations, &membership.DateReceived, &membership.DateProcessed, &membership.Notes)
		if err != nil {
			return nil, err
		}

		memberships = append(memberships, membership)
	}

	return memberships, nil
}
