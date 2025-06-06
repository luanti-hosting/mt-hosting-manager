package mtui

type Stats struct {
	Uptime      int     `json:"uptime"`
	MaxLag      float64 `json:"max_lag"`
	TimeOfDay   float64 `json:"time_of_day"`
	PlayerCount int     `json:"player_count"`
	Maintenance bool    `json:"maintenance"`
}
