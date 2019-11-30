package utils

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

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []int32, val int32) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func Exists(slice []int32, val int32) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func Intersection(a, b []int32) (c []int32) {
	m := make(map[int32]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}

func Minus(a, b []int32) (c []int32) {
	m := make(map[int32]uint8)
	for _, k := range a {
		m[k] |= 1 << 0
	}
	for _, k := range b {
		m[k] |= 1 << 1
	}
	var result []int32
	for k, v := range m {
		a := v&(1<<0) != 0
		b := v&(1<<1) != 0
		if a && !b {
			result = append(result, k)
		}
	}
	return result
}
