package domain

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrIncorrectSlugFormat = errors.New("Incorrect slug format")

type Slug struct {
	ID         int    `json:"id"`
	Text       string `json:"text"`
	NextNumber int    `json:"next_number"`
}

type UrlSlug struct {
	Text   string `json:"text"`
	Number int    `json:"number"`
}

type Url struct {
	ID   int     `json:"id"`
	Url  string  `json:"url"`
	Slug UrlSlug `json:"slug"`
}

func (u UrlSlug) String() string {
	return fmt.Sprintf("%s_%d", u.Text, u.Number)
}

func ParseUrlSlug(s string) (UrlSlug, error) {
	parts := strings.Split(s, "_")

	if len(parts) != 2 {
		return UrlSlug{}, ErrIncorrectSlugFormat
	}

	text := parts[0]
	number, err := strconv.Atoi(parts[1])
	if err != nil {
		return UrlSlug{}, ErrIncorrectSlugFormat
	}

	return UrlSlug{
		Text:   text,
		Number: number,
	}, nil
}
