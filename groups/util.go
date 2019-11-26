package groups

import (
	"strconv"
	"strings"
)

func GetIntID(name string) int32 {
	idStr, err := strconv.ParseInt(strings.Split(name, "::")[1], 10, 32)
	if err != nil {
		panic(err)
	}
	id := int32(int(idStr))
	return id
}
