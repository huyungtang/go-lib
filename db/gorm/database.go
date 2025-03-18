package gorm

import (
	"context"
	"database/sql"
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
func Init(dsn string, opts ...db.Option) (sqlDB db.SqlDB, err error) {
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
		return nil, errInvalidDriver
	}

	cfg := new(db.Context).
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
func (o *database) Table(ety any) db.Table {
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

// Ping
// ****************************************************************************************************************************************
func (o *database) Ping() (err error) {
	var db *sql.DB
	if db, err = o.DB.DB(); err != nil {
		return
	}

	return db.Ping()
}

// getTableName ***************************************************************************************************************************
func (o *database) getTableName(ety any) (name string) {
	if t, isMatched := ety.(db.TableName); isMatched {
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
func (o *database) getPrimaryKey(ety any) (pk string) {
	t := reflect.TypeOf(ety)
	v := reflect.ValueOf(ety)
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).IsExported() {
			continue
		}

		if _, isMatched := v.Field(i).Interface().(IdEntity); isMatched {
			pk = "id"
			break
		}

		tags := reflect.GetTags(t.Field(i), "gorm")
		if _, isMatched := tags["ignore"]; isMatched {
			continue
		}

		if _, isMatched := tags["embedded"]; isMatched {
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
	var isMatched bool
	tags := reflect.GetTags(f, "gorm")
	if nm, isMatched = tags["column"]; !isMatched {
		nm = f.Name
	}

	return
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
