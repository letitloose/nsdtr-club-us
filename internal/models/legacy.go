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

func (m *LegacyModel) List() ([]*LegacyMember, error) {

	stmt := `select  FirstName1, LastName1, HomePhone, EmailName, ` + "`Web Page`" + `, Region, DateJoined1, 
		HomeAddress1, HomeAddress2, HomeCity, HomeStateOrProvince, HomePostalCode, HomeCountry from Members`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []*LegacyMember{}

	for rows.Next() {
		member := &LegacyMember{Member: &Member{}, Address: &Address{}}

		err := rows.Scan(&member.FirstName, &member.LastName, &member.PhoneNumber,
			&member.Email, &member.Website, &member.Region, &member.JoinedDate, &member.Address1,
			&member.Address2, &member.City, &member.StateProvince, &member.ZipCode, &member.CountryCode)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}
