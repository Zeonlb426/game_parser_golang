package utils

import (
	"github.com/rs/zerolog"
	"runtime"
	"strconv"
)

func FileWithLineNums() *zerolog.Array {
	arr := zerolog.Arr()

	for i := 3; i <= 15; i++ {
		_, file, line, ok := runtime.Caller(i)

		if !ok {
			return arr
		}

		arr.Str(file + ":" + strconv.FormatInt(int64(line), 10))
	}

	return arr
}
