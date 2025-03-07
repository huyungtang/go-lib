package mongo

import (
	"context"
	"errors"

	"github.com/huyungtang/go-lib/db"
	"github.com/huyungtang/go-lib/reflect"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	base "go.mongodb.org/mongo-driver/mongo"
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

// collection *****************************************************************************************************************************
type collection struct {
	*base.Collection
	aggs []*aggregate
}

// GetById
// ****************************************************************************************************************************************
func (o *collection) GetById(ety, id any) (err error) {
	sid, isOK := id.(string)
	if !isOK || !primitive.IsValidObjectID(sid) {
		return errors.New("id is not a valid hex string")
	}

	oid, _ := primitive.ObjectIDFromHex(sid)
	res := o.Collection.FindOne(context.Background(), bson.M{"_id": oid})
	if err = res.Err(); err == nil {
		err = res.Decode(ety)
	}

	return
}

// GetPagedList
// ****************************************************************************************************************************************
func (o *collection) GetPagedList(ety any) (err error) {
	return o.Get(ety)
}

// AddFields
// ****************************************************************************************************************************************
func (o *collection) AddFields(m bson.M) db.Collection {
	o.aggs = append(o.aggs, &aggregate{agg: aggAddFields, cmd: m})

	return o
}

// Match
// ****************************************************************************************************************************************
func (o *collection) Match(m bson.M) db.Collection {
	o.aggs = append(o.aggs, &aggregate{agg: aggMatch, cmd: m})

	return o
}

// Offset
// ****************************************************************************************************************************************
func (o *collection) Offset(n int) db.Collection {
	o.aggs = append(o.aggs, &aggregate{agg: aggSkip, cmd: n})

	return o
}

// Limit
// ****************************************************************************************************************************************
func (o *collection) Limit(n int) db.Collection {
	o.aggs = append(o.aggs, &aggregate{agg: aggLimit, cmd: n})

	return o
}

// Create
// ****************************************************************************************************************************************
func (o *collection) Create(ety any) (err error) {
	if reflect.IsSlice(ety) {
		vals := reflect.ValueOf(ety)
		lens := vals.Len()
		dtos := make([]any, lens)
		for i := 0; i < lens; i++ {
			dtos[i] = vals.Index(i).Interface()
			beforeCreate(dtos[i])
		}

		var res *base.InsertManyResult
		if res, err = o.Collection.InsertMany(context.Background(), dtos); err == nil {
			for i := 0; i < lens; i++ {
				afterCreate(dtos[i], res.InsertedIDs[i])
			}
		}

	} else {
		beforeCreate(ety)

		var res *base.InsertOneResult
		if res, err = o.Collection.InsertOne(context.Background(), ety); err == nil {
			afterCreate(ety, res.InsertedID)
		}
	}

	return
}

// Get
// ****************************************************************************************************************************************
func (o *collection) Get(ety any) (err error) {
	var cur *base.Cursor
	tar := ety
	if p, isOK := ety.(db.Paged); isOK {
		var cnt int64
		if cnt, err = o.Count(); err != nil {
			return
		}
		p.SetCount(cnt)

		o.Offset(p.GetPageIndex() * p.GetPagedSize()).
			Limit(p.GetPagedSize())
		tar = p.GetDataDTO()
	}

	if cur, err = o.Collection.Aggregate(context.Background(), o.aggregates()); err != nil {
		return
	}
	defer cur.Close(context.Background())

	if reflect.IsSlice(tar) {
		err = cur.All(context.Background(), tar)
	} else if cur.TryNext(context.Background()) {
		err = cur.Decode(ety)
	}

	return
}

// Count
// ****************************************************************************************************************************************
func (o *collection) Count() (cnt int64, err error) {
	var cur *base.Cursor
	if cur, err = o.Collection.Aggregate(context.Background(), o.aggregates(aggregateCount)); err != nil {
		return
	}
	defer cur.Close(context.Background())

	if cur.TryNext(context.Background()) {
		ety := &struct {
			N int64 `bson:"n"`
		}{}
		if err = cur.Decode(ety); err != nil {
			return
		}

		cnt = ety.N
	}

	return
}

// aggregates *****************************************************************************************************************************
func (o *collection) aggregates(aggs ...*aggregate) any {
	cmds := append(o.aggs, aggs...)
	if lens := len(cmds); lens > 0 {
		ms := make([]any, lens)
		for i := 0; i < lens; i++ {
			ms[i] = cmds[i].getCmd()
		}

		return ms
	}

	return bson.D{}
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// beforeCreate ***************************************************************************************************************************
func beforeCreate(ety any) {
	if e, isOK := ety.(db.Created); isOK {
		e.Create()
	}
	if e, isOK := ety.(db.Updated); isOK {
		e.Update()
	}
}

// afterCreate ****************************************************************************************************************************
func afterCreate(ety, id any) {
	if etyId, isOK := ety.(db.Identity); isOK {
		etyId.SetId(id)
	}
}
