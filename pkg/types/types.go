package types

import (
	"strings"
	"time"

	"github.com/mvanbrummen/website-uptime-probe/pkg/db"
)

type Website struct {
	URL                  string    `json:"url,omitempty"`
	Active               bool      `json:"active,omitempty"`
	CreatedDate          time.Time `json:"created_date,omitempty"`
	ProbeScheduleMinutes int       `json:"probe_schedule_minutes,omitempty"`
}

func MapWebsites(source []db.WebsiteProbes) []Website {
	var websites []Website
	for _, w := range source {
		websites = append(websites, Website{
			URL:                  strings.Replace(w.SK, db.WebsiteSKPrefix, "", 1),
			Active:               w.Active,
			CreatedDate:          w.CreatedDate,
			ProbeScheduleMinutes: w.ProbeScheduleMinutes,
		})
	}
	return websites
}
