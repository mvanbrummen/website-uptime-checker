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

type Probe struct {
	URL string `json:"url,omitempty"`

	ResponseBody       string `json:"response_body,omitempty"`
	ResponseTimeMillis int    `json:"response_time_millis,omitempty"`
	ResponseStatus     int    `json:"response_status,omitempty"`
}

type CreateWebsiteRequest struct {
	URL                  string `json:"url,omitempty" binding:"required,url"`
	ProbeScheduleMinutes int    `json:"probe_schedule_minutes,omitempty" binding:"required,gte=1"`
}

type DeleteWebsiteRequest struct {
	URL string `json:"url,omitempty" binding:"required,url"`
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

func MapProbes(source []db.WebsiteProbes) []*Probe {
	probes := make([]*Probe, 0)
	for _, p := range source {
		probes = append(probes, MapProbe(&p))
	}
	return probes
}

func MapProbe(source *db.WebsiteProbes) *Probe {
	return &Probe{
		URL:                strings.Replace(source.SK, db.WebsiteSKPrefix, "", 1),
		ResponseBody:       source.ResponseBody,
		ResponseStatus:     source.ResponseStatus,
		ResponseTimeMillis: source.ResponseTimeMillis,
	}
}

func MapWebsite(source *db.WebsiteProbes) *Website {
	return &Website{
		URL:                  strings.Replace(source.SK, db.WebsiteSKPrefix, "", 1),
		Active:               source.Active,
		CreatedDate:          source.CreatedDate,
		ProbeScheduleMinutes: source.ProbeScheduleMinutes,
	}
}
