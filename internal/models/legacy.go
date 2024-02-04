package models

import (
	"database/sql"
)

type LegacyModel struct {
	DB *sql.DB
}

func (m *LegacyModel) List() ([]*Member, error) {

	stmt := `select MemberID, FirstName1, LastName1, HomePhone, EmailName, ` + "`Web Page`" + `, Region, DateJoined1, 
		HomeAddress1, HomeAddress2, HomeCity, HomeStateOrProvince, HomePostalCode, HomeCountry from Members`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []*Member{}

	for rows.Next() {
		member := &Member{Address: &Address{}}
		err := rows.Scan(&member.ID, &member.FirstName, &member.LastName, &member.PhoneNumber,
			&member.Email, &member.Website, &member.Region, &member.JoinedDate, &member.Address.Address1,
			&member.Address.Address2, &member.Address.City, &member.Address.StateProvince, &member.Address.ZipCode, &member.Address.CountryCode)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}
