package mongo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Identity
// ****************************************************************************************************************************************
type Identity struct {
	Id string `bson:"id"`
}

// SetId
// ****************************************************************************************************************************************
func (o *Identity) SetId(id any) {
	if id, isMatched := id.(primitive.ObjectID); isMatched {
		o.Id = id.Hex()
	}
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
