package models

import (
	"database/sql"
	"errors"
	"time"
)

type Member struct {
	ID          int
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Website     string
	Region      int
	CreatedDate time.Time
	JoinedDate  sql.NullTime
}

type MemberModel struct {
	DB *sql.DB
}

func (m *MemberModel) Insert(firstname, lastname, phonenumber, email, website string, region int, joined time.Time) (int, error) {

	stmt := `INSERT INTO members (firstname, lastname, phonenumber, email, website, region, created, joined)
    VALUES(?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), ?)`

	result, err := m.DB.Exec(stmt, firstname, lastname, phonenumber, email, website, region, joined)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

func (m *MemberModel) Get(id int) (*Member, error) {

	stmt := `select id, firstname, lastname, phonenumber, email, website, region, created, joined 
		from members 
    	where id = ?`

	result := m.DB.QueryRow(stmt, id)

	member := &Member{}
	err := result.Scan(&member.ID, &member.FirstName, &member.LastName, &member.PhoneNumber,
		&member.Email, &member.Website, &member.Region, &member.CreatedDate, &member.JoinedDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return member, nil
}

func (m *MemberModel) List() ([]*Member, error) {

	stmt := `select id, firstname, lastname, phonenumber, email, website, region, created, joined 
		from members`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []*Member{}

	for rows.Next() {
		member := &Member{}
		err := rows.Scan(&member.ID, &member.FirstName, &member.LastName, &member.PhoneNumber,
			&member.Email, &member.Website, &member.Region, &member.CreatedDate, &member.JoinedDate)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}
