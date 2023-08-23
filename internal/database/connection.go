package database

import (
	core "betassist.ru/bookmaker-game-parser/internal"
	"betassist.ru/bookmaker-game-parser/internal/database/dsn"
	"betassist.ru/bookmaker-game-parser/internal/database/logger"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	master *gorm.DB
	slave  *gorm.DB
)

func Open() {
	var err error

	if master, err = gorm.Open(postgres.Open(dsn.GetMaster()), getConfig()); nil != err {
		panic(err)
	}

	log.Debug().Msg("database => master connection established")

	if slave, err = gorm.Open(postgres.Open(dsn.GetSlave()), getConfig()); nil != err {
		panic(err)
	}

	log.Debug().Msg("database => slave connection established")

	prepare()
}

func GetMaster() *gorm.DB {
	return master
}

func GetSlave() *gorm.DB {
	return slave
}

func Close() {
	masterDB, _ := master.DB()
	_ = masterDB.Close()

	slaveDB, _ := slave.DB()
	_ = slaveDB.Close()

	log.Debug().Msg("database => connection closed")
}

func getConfig() *gorm.Config {
	return &gorm.Config{
		SkipDefaultTransaction:                   core.Config.Gorm.SkipDefaultTransaction,
		FullSaveAssociations:                     core.Config.Gorm.FullSaveAssociations,
		DryRun:                                   core.Config.Gorm.DryRun,
		PrepareStmt:                              core.Config.Gorm.PrepareStatement,
		DisableAutomaticPing:                     core.Config.Gorm.DisableAutomaticPing,
		DisableForeignKeyConstraintWhenMigrating: core.Config.Gorm.DisableForeignKeyConstraintWhenMigrating,
		DisableNestedTransaction:                 core.Config.Gorm.DisableNestedTransaction,
		AllowGlobalUpdate:                        core.Config.Gorm.AllowGlobalUpdate,
		QueryFields:                              core.Config.Gorm.QueryFields,
		CreateBatchSize:                          core.Config.Gorm.CreateBatchSize,
		Logger:                                   logger.New(),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   core.Config.Gorm.TablePrefix,
			SingularTable: core.Config.Gorm.SingularTableNames,
		},
	}
}

func prepare() {
	masterDB, _ := master.DB()
	masterDB.SetConnMaxIdleTime(core.Config.Gorm.MaxConnIdleTime)
	masterDB.SetConnMaxLifetime(core.Config.Gorm.MaxConnLifeTime)
	masterDB.SetMaxIdleConns(core.Config.Gorm.MaxIdleConnections)
	masterDB.SetMaxOpenConns(core.Config.Gorm.MaxOpenConnections)

	slaveDB, _ := slave.DB()
	slaveDB.SetConnMaxIdleTime(core.Config.Gorm.MaxConnIdleTime)
	slaveDB.SetConnMaxLifetime(core.Config.Gorm.MaxConnLifeTime)
	slaveDB.SetMaxIdleConns(core.Config.Gorm.MaxIdleConnections)
	slaveDB.SetMaxOpenConns(core.Config.Gorm.MaxOpenConnections)
}
