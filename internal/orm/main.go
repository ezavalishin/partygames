package orm

import (
	log "github.com/ezavalishin/partygames/internal/logger"
	"github.com/ezavalishin/partygames/internal/orm/migration"
	"github.com/ezavalishin/partygames/pkg/utils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var autoMigrate, logMode, seedDB bool
var dsn, dialect string

type ORM struct {
	DB *gorm.DB
}

func init() {
	dialect = utils.MustGet("GORM_DIALECT")
	dsn = utils.MustGet("GORM_CONNECTION_DSN")
	seedDB = utils.MustGetBool("GORM_SEED_DB")
	logMode = utils.MustGetBool("GORM_LOGMODE")
	autoMigrate = utils.MustGetBool("GORM_AUTOMIGRATE")
}

func Factory() (*ORM, error) {
	db, err := gorm.Open(dialect, dsn)

	if err != nil {
		log.Panic("[ORM] err: ", err)
	}

	orm := &ORM{
		DB: db,
	}

	db.LogMode(logMode)

	if autoMigrate {
		err = migration.ServiceAutoMigration(orm.DB)
	}

	log.Info("[ORM] Database connection initialized.")
	return orm, err
}
