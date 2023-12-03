package testdb

import "github.com/jaevor/go-nanoid"

const (
	nanoIDCharacters = "abcdefghijklmnopqrstuvwxyz"
	nanoIDLength     = 16
)

func generateRandomString() string {
	return must(nanoid.CustomASCII(nanoIDCharacters, nanoIDLength))()
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
