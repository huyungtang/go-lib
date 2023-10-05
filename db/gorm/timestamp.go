package gorm

import (
	"github.com/huyungtang/go-lib/times"
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

// Created
// ****************************************************************************************************************************************
type Created struct {
	CreatedAt int64 `gorm:"column:created_at;autoUpdateTime:false" json:",omitempty"`
}

// Create
// ****************************************************************************************************************************************
func (o *Created) Create() {
	o.CreatedAt = currnetTime()
}

// Updated
// ****************************************************************************************************************************************
type Updated struct {
	UpdatedAt int64 `gorm:"column:updated_at;autoUpdateTime:false" json:",omitempty"`
}

// Update
// ****************************************************************************************************************************************
func (o *Updated) Update() {
	o.UpdatedAt = currnetTime()
}

// Deleted
// ****************************************************************************************************************************************
type Deleted struct {
	DeletedAt int64 `gorm:"column:deleted_at;autoUpdateTime:false" json:",omitempty"`
}

// Delete
// ****************************************************************************************************************************************
func (o *Deleted) Delete() {
	o.DeletedAt = currnetTime()
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// currnetTime ****************************************************************************************************************************
func currnetTime() int64 {
	return times.Now().UnixMilli()
}
