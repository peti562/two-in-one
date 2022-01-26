package mocket

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	mocket "github.com/selvatico/go-mocket"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
)

type Driver struct {
	sqlite.Dialector
	DSN string
}

func Open(dsn string) gorm.Dialector {
	return &Driver{DSN: dsn}
}

func (dialector Driver) Name() string {
	return mocket.DriverName
}

func (dialector Driver) Initialize(db *gorm.DB) (err error) {
	// register callbacks
	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{
		LastInsertIDReversed: true,
	})
	db.ConnPool, err = sql.Open(dialector.Name(), dialector.DSN)
	return
}

func (dialector Driver) QuoteTo(writer clause.Writer, str string) {
	_ = writer.WriteByte('`')
	if strings.Contains(str, ".") {
		for idx, str := range strings.Split(str, ".") {
			if idx > 0 {
				_, _ = writer.WriteString(".`")
			}
			_, _ = writer.WriteString(str)
			_ = writer.WriteByte('`')
		}
	} else {
		_, _ = writer.WriteString(str)
		_ = writer.WriteByte('`')
	}
}
