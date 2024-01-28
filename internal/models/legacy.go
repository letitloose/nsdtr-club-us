package models

import (
	"database/sql"
)

type LegacyModel struct {
	DB *sql.DB
}

func (m *LegacyModel) List() ([]*Member, error) {

	stmt := "select MemberID, FirstName1, LastName1, HomePhone, EmailName, `Web Page`, Region, DateJoined1 from Members"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []*Member{}

	for rows.Next() {
		member := &Member{}
		err := rows.Scan(&member.ID, &member.FirstName, &member.LastName, &member.PhoneNumber,
			&member.Email, &member.Website, &member.Region, &member.JoinedDate)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}
