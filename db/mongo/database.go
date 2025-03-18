package mongo

import (
	"context"
	"errors"
	"net/url"

	"github.com/huyungtang/go-lib/db"
	"github.com/huyungtang/go-lib/reflect"
	"github.com/huyungtang/go-lib/strings"
	"go.mongodb.org/mongo-driver/bson"
	base "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(dsn string, opts ...db.Option) (nosql db.NoSqlDB, err error) {
	var u *url.URL
	if u, err = url.Parse(dsn); err != nil {
		return
	}

	var dbName string
	if dbName = u.Path; len(dbName) < 2 {
		return nil, errors.New("database name is not defined")
	}

	var client *base.Client
	if client, err = base.Connect(context.Background(), options.Client().ApplyURI(dsn)); err != nil {
		return
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return
	}

	return &database{client, client.Database(dbName[1:])}, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// database *******************************************************************************************************************************
type database struct {
	*base.Client
	*base.Database
}

// Collection
// ****************************************************************************************************************************************
func (o *database) Collection(ety any) db.Collection {
	return &collection{
		Collection: o.Database.Collection(o.getCollectionName(ety)),
		aggs:       []*aggregate{{agg: aggAddFields, cmd: bson.M{"id": bson.M{"$toString": "$_id"}}}},
	}
}

// Close
// ****************************************************************************************************************************************
func (o *database) Close() (err error) {
	if o.Client == nil {
		return
	}

	return o.Client.Disconnect(context.Background())
}

// getCollectionName **********************************************************************************************************************
func (o *database) getCollectionName(ety any) (name string) {
	if s, isMatched := ety.(string); isMatched {
		name = s
	} else if c, isMatched := ety.(db.CollectionName); isMatched {
		name = c.CollectionName()
	} else {
		name = reflect.TypeOf(ety).Name()
		if strings.HasSuffix(name, "entity") || strings.HasSuffix(name, "Entity") {
			name = strings.ToLower(name[:len(name)-6])
		}
	}

	return
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
