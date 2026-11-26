package util

import (
	"math/rand/v2"
	"regexp"
	"strings"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(text string) string {
	return strings.Trim(nonAlphanumeric.ReplaceAllString(strings.ToLower(text), "-"), "-")
}

func UniqueSlug(text string) string {
	const letters = "abcdefghjkmnpqrstuvwxyz23456789"
	code := make([]byte, 4)
	for index := range code {
		code[index] = letters[rand.IntN(len(letters))]
	}
	return Slugify(text) + "-" + string(code)
}
