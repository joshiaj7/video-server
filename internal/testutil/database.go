package testutil

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabase() (*gorm.DB, sqlmock.Sqlmock) {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      mockDB,
		SkipInitializeWithVersion: true,
	})
	gdb, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return gdb, sqlMock
}
