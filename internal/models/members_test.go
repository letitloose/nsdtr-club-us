package models

import (
	"database/sql"
	"testing"
	"time"
)

func TestInsertMember(t *testing.T) {
	db := NewTestDB(t)

	mm := MemberModel{DB: db}

	member := &Member{}
	_, err := mm.Insert(member.FirstName, member.LastName, "", "", "", member.Region, time.Now())
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddAddress(t *testing.T) {
	db := NewTestDB(t)

	mm := MemberModel{DB: db}

	member := &Member{}
	memberID, err := mm.Insert(member.FirstName, member.LastName, "", "", "", member.Region, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	am := AddressModel{DB: db}
	member.Address = &Address{City: sql.NullString{String: "NY", Valid: true}, StateProvince: sql.NullString{String: "NY", Valid: true}, CountryCode: sql.NullString{String: "USA", Valid: true}}
	addressID, err := am.Insert(member.Address)
	if err != nil {
		t.Fatal(err)
	}

	err = mm.AddAddress(memberID, addressID)
	if err != nil {
		t.Fatal(err)
	}

	member, err = mm.Get(memberID)
	if err != nil {
		t.Fatal(err)
	}

	address, err := am.Get(member.Address.ID)
	if err != nil {
		t.Fatal(err)
	}

	if address.City.String != "NY" {
		t.Fatal("wrong address returned")
	}
}
