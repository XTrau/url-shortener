package usecases

import (
	"math/rand"
	"strings"
)

const charSet = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM134567890"

func generateSlug(size int) string {
	var b strings.Builder
	b.Grow(size)

	for range size {
		ind := rand.Intn(len(charSet))
		b.Write([]byte(charSet[ind : ind+1]))
	}
	return b.String()
}
