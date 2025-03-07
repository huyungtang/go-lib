package mongo

import (
	"go.mongodb.org/mongo-driver/bson"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

const (
	aggAddFields aggEnum = "$addFields"
	aggAnd       aggEnum = "$and"
	aggCount     aggEnum = "$count"
	aggLimit     aggEnum = "$limit"
	aggMatch     aggEnum = "$match"
	aggSkip      aggEnum = "$skip"
)

var (
	aggregateCount = &aggregate{agg: aggCount, cmd: "n"}
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// aggregate
// ****************************************************************************************************************************************
type aggregate struct {
	agg aggEnum
	cmd any
}

// getCmd *********************************************************************************************************************************
func (o *aggregate) getCmd() any {
	m := bson.M{}
	m[o.agg] = o.cmd

	return m
}

// aggEnum ********************************************************************************************************************************
type aggEnum = string

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
