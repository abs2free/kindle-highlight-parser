package main

import (
	"regexp"
	"strconv"
)

func extractNumber(s string) (int, error) {
	re := regexp.MustCompile("[0-9]+")
	numberS := re.FindString(s)
	return strconv.Atoi(numberS)
}
