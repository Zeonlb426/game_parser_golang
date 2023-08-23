package helpers

import (
	core "betassist.ru/bookmaker-game-parser/internal"
	"fmt"
)

func WithPrefix(baseName string) string {
	return fmt.Sprintf("%s%s", core.Config.Gorm.TablePrefix, baseName)
}
