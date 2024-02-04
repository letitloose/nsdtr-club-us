package models

import (
	"database/sql"
	"errors"
)

type Address struct {
	ID            int
	Address1      string
	Address2      sql.NullString
	City          sql.NullString
	StateProvince sql.NullString
	ZipCode       sql.NullString
	CountryCode   sql.NullString
}

type AddressModel struct {
	DB *sql.DB
}

func (m *AddressModel) Insert(address *Address) (int, error) {

	statement := `INSERT INTO address (address1, address2, city, stateProvince, zipCode, country)
    VALUES(?, ?, ?, ?, ?, ?)`

	result, err := m.DB.Exec(statement, address.Address1, address.Address2, address.City, address.StateProvince, address.ZipCode, address.CountryCode)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *AddressModel) Get(id int) (*Address, error) {

	stmt := `select id, address1, address2, city, stateProvince, zipCode, country
		from address where id = ?`

	result := m.DB.QueryRow(stmt, id)

	address := &Address{}
	err := result.Scan(&address.ID, &address.Address1, &address.Address2, &address.City,
		&address.StateProvince, &address.ZipCode, &address.CountryCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return address, nil
}
