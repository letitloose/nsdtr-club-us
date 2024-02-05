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
	PhoneNumber sql.NullString
	Email       sql.NullString
	Website     sql.NullString
	Region      int
	CreatedDate time.Time
	JoinedDate  sql.NullTime
	AddressID   sql.NullInt16
}

type MemberProfile struct {
	*Member
	*Address
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

	stmt := `select id, firstname, lastname, phonenumber, email, website, region, created, joined, addressID 
		from members 
    	where id = ?`

	result := m.DB.QueryRow(stmt, id)

	member := &Member{}
	err := result.Scan(&member.ID, &member.FirstName, &member.LastName, &member.PhoneNumber,
		&member.Email, &member.Website, &member.Region, &member.CreatedDate, &member.JoinedDate, &member.AddressID)
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

func (m *MemberModel) AddAddress(memberID, addressID int) error {

	stmt := `UPDATE members set addressID = ? where id = ?`

	_, err := m.DB.Exec(stmt, addressID, memberID)

	return err
}

func (m *MemberModel) GetMemberProfile(id int) (*MemberProfile, error) {

	stmt := `select m.id, m.firstname, m.lastname, m.phonenumber, m.email, m.website, m.region, m.created, m.joined, 
				a.address1, a.address2, a.city, a.stateProvince, a.zipCode, a.country
		from members m
		join address a on a.id = m.addressID
    	where m.id = ?`

	result := m.DB.QueryRow(stmt, id)

	member := &MemberProfile{Member: &Member{}, Address: &Address{}}

	err := result.Scan(&member.Member.ID, &member.FirstName, &member.LastName, &member.PhoneNumber,
		&member.Email, &member.Website, &member.Region, &member.JoinedDate, &member.CreatedDate, &member.Address1,
		&member.Address2, &member.City, &member.StateProvince, &member.ZipCode, &member.CountryCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return member, nil
}
