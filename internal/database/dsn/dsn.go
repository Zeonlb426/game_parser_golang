package dsn

import (
	core "betassist.ru/bookmaker-game-parser/internal"
	"fmt"
)

func GetMaster() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		core.Config.Db.HostMaster,
		core.Config.Db.PortMaster,
		core.Config.Db.User,
		core.Config.Db.Password,
		core.Config.Db.NameMaster,
	)
}

func GetSlave() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		core.Config.Db.HostSlave,
		core.Config.Db.PortSlave,
		core.Config.Db.User,
		core.Config.Db.Password,
		core.Config.Db.NameSlave,
	)
}
