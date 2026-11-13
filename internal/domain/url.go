package domain

type Slug struct {
	ID         int    `json:"id"`
	Text       string `json:"slug"`
	NextNumber int    `json:"next_number"`
}

type UrlSlug struct {
	Text   string `json:"slug"`
	Number int    `json:"number"`
}

func (u UrlSlug) String() string {
	return ""
}

func ParseUrlSlug(s string) (UrlSlug, error) {
	return UrlSlug{}, nil
}

type Url struct {
	ID   int     `json:"id"`
	Url  string  `json:"url"`
	Slug UrlSlug `json:"slug"`
}
