package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Membership struct {
	ID               int
	MemberID         int
	Year             int
	MembershipType   string
	MembershipAmount sql.NullFloat64
	PrAmount         sql.NullFloat64
	HealthAmount     sql.NullFloat64
	ResueAmount      sql.NullFloat64
	TotalPaid        sql.NullFloat64
	Items            []*MembershipItem
}

type MembershipItem struct {
	ID           int
	MembershipID int
	ItemCode     string
	AmountPaid   float32
}

type DueScheduleItem struct {
	Code    string
	Display string
	Cost    float32
	Year    int
}

type DueSchedule struct {
	Items []*DueScheduleItem
}

type MembershipModel struct {
	DB *sql.DB
}

func (m *Membership) String() string {
	return fmt.Sprint(m.Year)
}

func (m *MembershipModel) Insert(membership *Membership) (int, error) {

	statement := `INSERT INTO membership (memberID, year)
    VALUES(?, ?)`

	result, err := m.DB.Exec(statement, membership.MemberID, membership.Year)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *MembershipModel) GetMemberships(memberID int) ([]*Membership, error) {

	stmt := `select m.id,
		m.year,
		mtype.display,
		mtype.amountPaid as typeAmount,
		printedRoster.amountPaid as prAmount,
		health.amountPaid as healthAmount,
		rescue.amountPaid as rescueAmount,
		IFNULL(mtype.amountPaid, 0) + IFNULL(printedRoster.amountPaid, 0) + IFNULL(health.amountPaid, 0) + IFNULL(rescue.amountPaid, 0) as totalPaid
	from membership m
	left join (select item.display, mi.amountPaid, mi.membershipID  from membershipItem mi join dueSchedule item on item.code = mi.itemCode and mi.itemCode in ('SI', 'JT', 'JR', 'CM', 'FM'))  mtype on mtype.membershipID = m.id 
	left join (select item.display, mi.amountPaid, mi.membershipID  from membershipItem mi join dueSchedule item on item.code = mi.itemCode and mi.itemCode in ('PR'))  printedRoster on printedRoster.membershipID = m.id 
	left join (select item.display, mi.amountPaid, mi.membershipID  from membershipItem mi join dueSchedule item on item.code = mi.itemCode and mi.itemCode in ('HG'))  health on health.membershipID = m.id 
	left join (select item.display, mi.amountPaid, mi.membershipID  from membershipItem mi join dueSchedule item on item.code = mi.itemCode and mi.itemCode in ('RS'))  rescue on rescue.membershipID = m.id 
	where m.memberID = ?;`

	rows, err := m.DB.Query(stmt, memberID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberships := []*Membership{}
	for rows.Next() {
		membership := &Membership{}
		err := rows.Scan(&membership.ID, &membership.Year, &membership.MembershipType, &membership.MembershipAmount, &membership.PrAmount, &membership.HealthAmount, &membership.ResueAmount, &membership.TotalPaid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNoRecord
			} else {
				return nil, err
			}
		}

		membership.Items, err = m.GetMembershipItems(membership.ID)
		if err != nil {
			return nil, err
		}
		memberships = append(memberships, membership)
	}

	return memberships, nil
}

func (m *MembershipModel) InsertMembershipItem(item *MembershipItem) (int, error) {

	statement := `INSERT INTO membershipItem (membershipID, itemCode, amountPaid)
    VALUES(?, ?, ?)`

	result, err := m.DB.Exec(statement, item.MembershipID, item.ItemCode, item.AmountPaid)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *MembershipModel) GetMembershipItems(membershipID int) ([]*MembershipItem, error) {

	stmt := `select id, membershipID, itemCode, amountPaid
		from membershipItem where membershipID = ?`

	rows, err := m.DB.Query(stmt, membershipID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*MembershipItem{}
	for rows.Next() {
		item := &MembershipItem{}
		err := rows.Scan(&item.ID, &item.MembershipID, &item.ItemCode, &item.AmountPaid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNoRecord
			} else {
				return nil, err
			}
		}
		items = append(items, item)
	}

	return items, nil
}

func (m *MembershipModel) LookupDueScheduleItem(display string) (*DueScheduleItem, error) {

	stmt := `select code, display, cost, year from dueSchedule where display = ?`

	result := m.DB.QueryRow(stmt, display)

	item := &DueScheduleItem{}

	err := result.Scan(&item.Code, &item.Display, &item.Cost, &item.Year)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return item, nil
}

func (m *MembershipModel) GetDueSchedule() (*DueSchedule, error) {

	stmt := `select code, display, cost, year from dueSchedule`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []*DueScheduleItem{}
	for rows.Next() {
		item := &DueScheduleItem{}
		err := rows.Scan(&item.Code, &item.Display, &item.Cost, &item.Year)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNoRecord
			} else {
				return nil, err
			}
		}
		items = append(items, item)
	}

	return &DueSchedule{Items: items}, nil
}
