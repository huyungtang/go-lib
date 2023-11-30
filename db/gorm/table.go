package gorm

import (
	"github.com/huyungtang/go-lib/db"
	"github.com/huyungtang/go-lib/reflect"
	"github.com/huyungtang/go-lib/strings"
	base "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// table **********************************************************************************************************************************
type table struct {
	*base.DB
	pKey   string
	entity interface{}
}

// GetById
// ****************************************************************************************************************************************
func (o *table) GetById(ety, id interface{}) (err error) {
	return o.DB.
		Where(
			strings.Format("%s = @id", o.pKey),
			map[string]interface{}{"id": id},
		).Take(ety).Error
}

// GetPagedList
// ****************************************************************************************************************************************
func (o *table) GetPagedList(ety interface{}) (err error) {
	return o.Get(ety)
}

// LockForUpdate
// ****************************************************************************************************************************************
func (o *table) LockForUpdate() db.Table {
	o.DB = o.DB.Clauses(clause.Locking{
		Strength: "UPDATE",
	})

	return o
}

// Select
// ****************************************************************************************************************************************
func (o *table) Select(qry interface{}, cols ...interface{}) db.Table {
	o.DB = o.DB.Select(qry, cols...)

	return o
}

// Omit
// ****************************************************************************************************************************************
func (o *table) Omit(cols ...string) db.Table {
	o.DB = o.DB.Omit(cols...)

	return o
}

// Join
// ****************************************************************************************************************************************
func (o *table) Join(qry string, args ...interface{}) db.Table {
	o.DB = o.DB.Joins(qry, args...)

	return o
}

// Preload
// ****************************************************************************************************************************************
func (o *table) Preload(qry string, args ...interface{}) db.Table {
	o.DB = o.DB.Preload(qry, args...)

	return o
}

// Available
// ****************************************************************************************************************************************
func (o *table) Available() db.Table {
	if _, isOK := o.entity.(db.Deleted); isOK {
		o.DB = o.DB.Where("deleted_at = 0")
	}

	return o
}

// Where
// ****************************************************************************************************************************************
func (o *table) Where(qry interface{}, args ...interface{}) db.Table {
	o.DB = o.DB.Where(qry, args...)

	return o
}

// Having
// ****************************************************************************************************************************************
func (o *table) Having(qry interface{}, args ...interface{}) db.Table {
	o.DB = o.DB.Having(qry, args...)

	return o
}

// Order
// ****************************************************************************************************************************************
func (o *table) Order(col interface{}) db.Table {
	o.DB = o.DB.Order(col)

	return o
}

// Group
// ****************************************************************************************************************************************
func (o *table) Group(col string) db.Table {
	o.DB = o.DB.Group(col)

	return o
}

// Offset
// ****************************************************************************************************************************************
func (o *table) Offset(n int) db.Table {
	o.DB = o.DB.Offset(n)

	return o
}

// Limit
// ****************************************************************************************************************************************
func (o *table) Limit(n int) db.Table {
	o.DB = o.DB.Limit(n)

	return o
}

// Create
// ****************************************************************************************************************************************
func (o *table) Create(ety interface{}) (err error) {
	if reflect.IsSlice(ety) {
		vals := reflect.ValueOf(ety)
		for i := 0; i < vals.Len(); i++ {
			beforeCreate(vals.Index(i).Interface())
		}
	} else {
		beforeCreate(ety)
	}

	return o.DB.Create(ety).Error
}

// Get
// ****************************************************************************************************************************************
func (o *table) Get(ety interface{}) (err error) {
	tar := ety
	if p, isOK := ety.(db.Paged); isOK {
		var cnt int64
		if cnt, err = o.Count(); err != nil {
			return
		}

		p.SetCount(cnt)
		o.DB = o.DB.
			Offset(p.GetPageIndex() * p.GetPagedSize()).
			Limit(p.GetPagedSize())
		tar = p.GetDataDTO()
	}

	return o.DB.Find(tar).Error
}

// Update
// ****************************************************************************************************************************************
func (o *table) Update(ety interface{}) (err error) {
	if e, isOK := ety.(db.Updated); isOK {
		e.Update()
		o.DB.Statement.Selects = append(o.DB.Statement.Selects, "updated_at")
	}

	return o.DB.Updates(ety).Error
}

// Delete
// ****************************************************************************************************************************************
func (o *table) Delete() (err error) {
	if e, isOK := o.entity.(db.Deleted); isOK {
		e.Delete()

		return o.DB.Select("deleted_at").Updates(e).Error
	}

	return o.DB.Delete(o.entity).Error
}

// Count
// ****************************************************************************************************************************************
func (o *table) Count() (cnt int64, err error) {
	var p map[string][]interface{}
	p, o.DB.Statement.Preloads = o.DB.Statement.Preloads, make(map[string][]interface{})
	err = o.DB.Count(&cnt).Error
	o.DB.Statement.Preloads = p

	return
}

// Exec
// ****************************************************************************************************************************************
func (o *table) Exec(sqlcmd string, args ...interface{}) (err error) {
	return o.DB.Exec(sqlcmd, args...).Error
}

// Begin
// ****************************************************************************************************************************************
func (o *table) Begin() db.Transaction {
	tx := &table{
		DB:     o.DB.Begin(),
		pKey:   o.pKey,
		entity: o.entity,
	}

	return tx
}

// Rollback
// ****************************************************************************************************************************************
func (o *table) Rollback() {
	o.DB.Rollback()
}

// Commit
// ****************************************************************************************************************************************
func (o *table) Commit() (err error) {
	return o.DB.Commit().Error
}

// ResetClauses
// ****************************************************************************************************************************************
func (o *table) ResetClauses() {
	o.DB.Statement.Clauses = make(map[string]clause.Clause)
}

// CreateColumns
// ****************************************************************************************************************************************
func (o *table) CreateColumns() (cols []string) {
	tp := reflect.TypeOf(o.entity)
	cols = make([]string, 0, tp.NumField())
	for i := 0; i < tp.NumField(); i++ {
		tg := reflect.GetTags(tp.Field(i), "gorm")

		col, isOK := tg["column"]
		if !isOK {
			continue
		}

		per, isOK := tg["<-"]
		if isOK && per == "update" {
			continue
		}

		cols = append(cols, col)
	}

	return
}

// UpdateColumns
// ****************************************************************************************************************************************
func (o *table) UpdateColumns() (cols []string) {
	tp := reflect.TypeOf(o.entity)
	cols = make([]string, 0, tp.NumField())
	for i := 0; i < tp.NumField(); i++ {
		tg := reflect.GetTags(tp.Field(i), "gorm")

		col, isOK := tg["column"]
		if !isOK {
			continue
		}

		per, isOK := tg["<-"]
		if isOK && per == "create" {
			continue
		}

		cols = append(cols, col)
	}

	return
}

// RowsAffected
// ****************************************************************************************************************************************
func (o *table) RowsAffected() int64 {
	return o.DB.RowsAffected
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// beforeCreate ***************************************************************************************************************************
func beforeCreate(ety interface{}) {
	if e, isOK := ety.(db.Created); isOK {
		e.Create()
	}
	if e, isOK := ety.(db.Updated); isOK {
		e.Update()
	}
}
