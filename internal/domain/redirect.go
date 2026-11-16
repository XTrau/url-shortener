package domain

import "time"

type RedirectInfo struct {
	CreatedAt time.Time
	IP        string
	UserAgent string
	UrlID     int
}
