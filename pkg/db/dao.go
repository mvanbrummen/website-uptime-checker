package db

import (
	"time"

	"github.com/guregu/dynamo"
)

const (
	WebsiteProbesTable = "WebsiteProbes"
	WebsitePK          = "WEBSITE#"
	WebsiteSKPrefix    = "WEBSITE#"
)

type WebsiteProbes struct {
	PK string `dynamo:",hash"`
	SK string `dynamo:",range"`

	ProbeAttributes
	WebsiteAttributes
}

type ProbeAttributes struct {
	ResponseBody       string `dynamo:"ResponseBody"`
	ResponseTimeMillis int    `dynamo:"ResponseTimeMillis"`
	ResponseStatus     int    `dynamo:"ResponseStatus"`
}

type WebsiteAttributes struct {
	ProbeScheduleMinutes int       `dynamo:"ProbeScheduleMinutes"`
	CreatedDate          time.Time `dynamo:"CreatedDate"`
	Active               bool      `dynamo:"Active"`
}

type ProbesDao struct {
	DB *dynamo.DB
}

func NewProbesDao(db *dynamo.DB) *ProbesDao {
	return &ProbesDao{
		DB: db,
	}
}

func (d *ProbesDao) GetWebsites() []WebsiteProbes {
	t := d.DB.Table(WebsiteProbesTable)

	var websites []WebsiteProbes
	t.Get("PK", WebsitePK).
		Range("SK", dynamo.BeginsWith, WebsiteSKPrefix).
		All(&websites)

	return websites
}
