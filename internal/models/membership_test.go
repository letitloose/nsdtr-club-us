package models

import (
	"database/sql"
	"testing"
	"time"
)

func TestInsertMembership(t *testing.T) {
	db := NewTestDB(t)

	mmod := MemberModel{DB: db}

	member := &Member{FirstName: "Lou", LastName: "Garwood", JointFirstName: sql.NullString{String: "Annie"}, JointLastName: sql.NullString{String: "Garwood"}, Region: 1}
	memberID, err := mmod.Insert(member.FirstName, member.LastName, member.JointFirstName.String, member.JointLastName.String, "", "", "", member.Region, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	mm := MembershipModel{DB: db}

	membership := &Membership{MemberID: memberID, Year: 2003}
	_, err = mm.Insert(membership)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMemberships(t *testing.T) {
	db := NewTestDB(t)

	mmod := MemberModel{DB: db}

	member := &Member{FirstName: "Lou", LastName: "Garwood", JointFirstName: sql.NullString{String: "Annie"}, JointLastName: sql.NullString{String: "Garwood"}, Region: 1}
	memberID, err := mmod.Insert(member.FirstName, member.LastName, member.JointFirstName.String, member.JointLastName.String, "", "", "", member.Region, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	mm := MembershipModel{DB: db}

	membership := &Membership{MemberID: memberID, Year: 2003}
	membershipID, err := mm.Insert(membership)
	if err != nil {
		t.Fatal(err)
	}

	item := &MembershipItem{MembershipID: membershipID, ItemCode: "SI", AmountPaid: 0.0}
	_, err = mm.InsertMembershipItem(item)
	if err != nil {
		t.Fatal(err)
	}

	memberships, err := mm.GetMemberships(memberID)
	if err != nil {
		t.Fatal(err)
	}

	expected := 1
	if len(memberships) != expected {
		t.Fatalf("wrong number of results returned. expecting %d, got %d", expected, len(memberships))
	}

	if len(memberships[0].Items) != expected {
		t.Fatalf("wrong number of results returned. expecting %d, got %d", expected, len(memberships[0].Items))
	}
}

func TestInsertMembershipItem(t *testing.T) {
	db := NewTestDB(t)

	mmod := MemberModel{DB: db}

	member := &Member{FirstName: "Lou", LastName: "Garwood", JointFirstName: sql.NullString{String: "Annie"}, JointLastName: sql.NullString{String: "Garwood"}, Region: 1}
	memberID, err := mmod.Insert(member.FirstName, member.LastName, member.JointFirstName.String, member.JointLastName.String, "", "", "", member.Region, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	mm := MembershipModel{DB: db}

	membership := &Membership{MemberID: memberID, Year: 2003}
	membershipID, err := mm.Insert(membership)
	if err != nil {
		t.Fatal(err)
	}

	item := &MembershipItem{MembershipID: membershipID, ItemCode: "SI", AmountPaid: 0.0}
	_, err = mm.InsertMembershipItem(item)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetMembershipItem(t *testing.T) {
	db := NewTestDB(t)

	mmod := MemberModel{DB: db}

	member := &Member{FirstName: "Lou", LastName: "Garwood", JointFirstName: sql.NullString{String: "Annie"}, JointLastName: sql.NullString{String: "Garwood"}, Region: 1}
	memberID, err := mmod.Insert(member.FirstName, member.LastName, member.JointFirstName.String, member.JointLastName.String, "", "", "", member.Region, time.Now())
	if err != nil {
		t.Fatal(err)
	}

	mm := MembershipModel{DB: db}

	membership := &Membership{MemberID: memberID, Year: 2003}
	membershipID, err := mm.Insert(membership)
	if err != nil {
		t.Fatal(err)
	}

	item := &MembershipItem{MembershipID: membershipID, ItemCode: "SI", AmountPaid: 0.0}
	_, err = mm.InsertMembershipItem(item)
	if err != nil {
		t.Fatal(err)
	}

	items, err := mm.GetMembershipItems(membershipID)
	if err != nil {
		t.Fatal(err)
	}

	expected := 1
	if len(items) != expected {
		t.Fatalf("wrong number of results returned. expecting %d, got %d", expected, len(items))
	}
}

func TestLookupDueScheduleItem(t *testing.T) {
	db := NewTestDB(t)

	mm := MembershipModel{DB: db}

	item, err := mm.LookupDueScheduleItem("Single")
	if err != nil {
		t.Fatal(err)
	}

	expected := "SI"
	if item.Code != expected {
		t.Fatalf("wrong number of results returned. expecting %s, got %s", expected, item.Code)
	}
}

func TestGetDueScheduleItem(t *testing.T) {
	db := NewTestDB(t)

	mm := MembershipModel{DB: db}

	doeSchedule, err := mm.GetDueSchedule()
	if err != nil {
		t.Fatal(err)
	}

	expected := 9
	got := len(doeSchedule.Items)
	if got != expected {
		t.Fatalf("wrong number of results returned. expecting %d, got %d", expected, got)
	}
}
