package gorm

import (
	"github.com/huyungtang/go-lib/db"
	"github.com/huyungtang/go-lib/reflect"
	"github.com/huyungtang/go-lib/slices"
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
	entity any
}

// GetById
// ****************************************************************************************************************************************
func (o *table) GetById(ety, id any) (err error) {
	return o.DB.
		Where(strings.Format("%s = @id", o.pKey), map[string]any{"id": id}).Take(ety).Error
}

// GetPagedList
// ****************************************************************************************************************************************
func (o *table) GetPagedList(ety any) (err error) {
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
func (o *table) Select(qry any, cols ...any) db.Table {
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
func (o *table) Join(qry string, args ...any) db.Table {
	o.DB = o.DB.Joins(qry, args...)

	return o
}

// Preload
// ****************************************************************************************************************************************
func (o *table) Preload(qry string, args ...any) db.Table {
	o.DB = o.DB.Preload(qry, args...)

	return o
}

// Available
// ****************************************************************************************************************************************
func (o *table) Available() db.Table {
	if _, isMatched := o.entity.(db.Deleted); isMatched {
		o.DB = o.DB.Where("deleted_at = 0")
	}

	return o
}

// Where
// ****************************************************************************************************************************************
func (o *table) Where(qry any, args ...any) db.Table {
	o.DB = o.DB.Where(qry, args...)

	return o
}

// Having
// ****************************************************************************************************************************************
func (o *table) Having(qry any, args ...any) db.Table {
	o.DB = o.DB.Having(qry, args...)

	return o
}

// Order
// ****************************************************************************************************************************************
func (o *table) Order(col any) db.Table {
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
func (o *table) Create(ety any) (err error) {
	if reflect.IsSlice(ety) {
		vals := reflect.ValueOf(ety)
		for i := 0; i < vals.Len(); i++ {
			o.beforeCreate(vals.Index(i).Interface())
		}
	} else {
		o.beforeCreate(ety)
	}

	return o.DB.Create(ety).Error
}

// Get
// ****************************************************************************************************************************************
func (o *table) Get(ety any) (err error) {
	tar := ety
	if p, isMatched := ety.(db.Paged); isMatched {
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
func (o *table) Update(ety any) (err error) {
	if e, isMatched := ety.(db.Updated); isMatched {
		e.Update()
		o.DB.Statement.Selects = append(o.DB.Statement.Selects, "updated_at")
	}

	return o.DB.Updates(ety).Error
}

// Delete
// ****************************************************************************************************************************************
func (o *table) Delete(ety any) (err error) {
	if e, isMatched := ety.(db.Deleted); isMatched {
		e.Delete()
		o.DB.Statement.Selects = append(o.DB.Statement.Selects, "deleted_at")

		return o.DB.Updates(e).Error
	}

	return o.DB.Delete(ety).Error
}

// Count
// ****************************************************************************************************************************************
func (o *table) Count() (cnt int64, err error) {
	var p map[string][]any
	p, o.DB.Statement.Preloads = o.DB.Statement.Preloads, make(map[string][]any)
	err = o.DB.Count(&cnt).Error
	o.DB.Statement.Preloads = p

	return
}

// Exec
// ****************************************************************************************************************************************
func (o *table) Exec(sqlcmd string, args ...any) (err error) {
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

	return o.getColumns(o.entity, "update")
}

// UpdateColumns
// ****************************************************************************************************************************************
func (o *table) UpdateColumns() (cols []string) {

	return o.getColumns(o.entity, "create")
}

// RowsAffected
// ****************************************************************************************************************************************
func (o *table) RowsAffected() int64 {
	return o.DB.RowsAffected
}

// SubQuery
// ****************************************************************************************************************************************
func (o *table) SubQuery() any {
	return o.DB
}

// beforeCreate ***************************************************************************************************************************
func (o *table) beforeCreate(ety any) {
	if e, isMatched := ety.(db.Created); isMatched {
		e.Create()
	}

	if e, isMatched := ety.(db.Updated); isMatched {
		e.Update()
	}
}

// getColumns *****************************************************************************************************************************
func (o *table) getColumns(entity any, exp string) (cols []string) {
	tp := reflect.TypeOf(entity)
	cs := slices.New()
	for i := 0; i < tp.NumField(); i++ {
		tags := reflect.GetTags(tp.Field(i), "gorm")
		_, isIgnore := tags["ignore"]
		_, isEmbed := tags["embedded"]
		_, isPrimary := tags["primaryKey"]
		col, isCol := tags["column"]
		_, isRead := tags["->"]
		write, isWrite := tags["<-"]
		if write == "false" {
			isWrite = false
		}

		if isIgnore ||
			isPrimary ||
			(isRead && !isWrite) ||
			(!isCol && !isEmbed) ||
			(isWrite && write == exp) {
			continue
		}

		if isEmbed {
			d := reflect.New(tp.Field(i).Type)
			ss := o.getColumns(d.Interface(), exp)
			for _, s := range ss {
				cs.Append(s)
			}

			continue
		}

		cs.Append(col)
	}

	return cs.GetString()
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
