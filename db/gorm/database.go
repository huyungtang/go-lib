package gorm

import (
	"context"
	"database/sql"
	"errors"
	reflect_ "reflect"

	"github.com/huyungtang/go-lib/db"
	"github.com/huyungtang/go-lib/reflect"
	"github.com/huyungtang/go-lib/strings"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	base "gorm.io/gorm"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(dsn string, opts ...db.Options) (sqlDB db.SqlDB, err error) {
	var dial base.Dialector
	d := strings.ToLower(strings.Find(dsn, `^(?i)(mssql|mysql|postgres)://`))
	switch d {
	case "mssql://":
		dial = sqlserver.Open(dsn[len(d):])
	case "mysql://":
		dial = mysql.Open(dsn[len(d):])
	case "postgres://":
		dial = postgres.Open(dsn[len(d):])
	default:
		return nil, errors.New("cannot identify the database driver")
	}

	cfg := new(db.Option).
		ApplyOptions(opts,
			db.SkipDefaultTransactionOption(true),
		)

	var conn *base.DB
	if conn, err = base.Open(dial, &base.Config{
		SkipDefaultTransaction: cfg.SkipDefaultTransaction,
	}); err != nil {
		return
	}

	if cfg.Debug {
		conn = conn.Debug()
	}

	return &database{conn}, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// database
// ****************************************************************************************************************************************
type database struct {
	*base.DB
}

// Table
// ****************************************************************************************************************************************
func (o *database) Table(ety interface{}) db.Table {
	return &table{
		DB:     o.DB.WithContext(context.Background()).Table(o.getTableName(ety)),
		pKey:   o.getPrimaryKey(ety),
		entity: ety,
	}
}

// Close
// ****************************************************************************************************************************************
func (o *database) Close() (err error) {
	if o.DB == nil {
		return
	}

	var conn *sql.DB
	if conn, err = o.DB.DB(); err != nil {
		return
	}

	return conn.Close()
}

// getTableName ***************************************************************************************************************************
func (o *database) getTableName(ety interface{}) (name string) {
	if t, isOK := ety.(db.TableName); isOK {
		name = t.TableName()
	} else {
		name = strings.ToLower(reflect.TypeOf(ety).Name())
		if strings.HasSuffix(name, "entity") {
			name = name[:len(name)-6]
		}
	}

	return
}

// getPrimaryKey **************************************************************************************************************************
func (o *database) getPrimaryKey(ety interface{}) (pk string) {
	t := reflect.TypeOf(ety)
	v := reflect.ValueOf(ety)
	for i := 0; i < t.NumField(); i++ {
		if _, isMap := v.Field(i).Interface().(Identity); isMap {
			pk = "id"
			break
		}

		tags := reflect.GetTags(t.Field(i), "gorm")
		if _, isMap := tags["ignore"]; isMap {
			continue
		}

		if _, isMap := tags["embedded"]; isMap {
			continue
		}

		if _, isKey := tags["primaryKey"]; pk == "" || isKey {
			pk = o.getFieldName(t.Field(i))

			if isKey {
				break
			}
		}
	}

	return
}

// getFieldName ***************************************************************************************************************************
func (o *database) getFieldName(f reflect_.StructField) (nm string) {
	var isMap bool
	tags := reflect.GetTags(f, "gorm")
	if nm, isMap = tags["column"]; !isMap {
		nm = f.Name
	}

	return
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
