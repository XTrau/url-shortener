package domain

type SummaryStatistics struct {
	Url string `json:"url"`

	Total     int `json:"total"`
	UniqueIPs int `json:"unique"`

	Total24H int `json:"total_24h"`
	Total7D  int `json:"total_7d"`
	Total30D int `json:"total_30d"`

	Unique24H int `json:"unique_24h"`
	Unique7D  int `json:"unique_7d"`
	Unique30D int `json:"unique_30d"`
}

type PeriodStats struct {
	Label     string `json:"labels"`
	Redirects int    `json:"redirects"`
	Unique    int    `json:"unique"`
}

type TimelineStatistics struct {
	Period string        `json:"period"`
	Data   []PeriodStats `json:"data"`
}
