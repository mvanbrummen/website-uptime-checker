package db

import (
	"time"

	"github.com/guregu/dynamo"
)

type WebsiteProbes struct {
	PK string `dynamo:",hash"`
	SK string `dynamo:",range"`

	ProbeAttributes
	WebiteAttributes
}

type ProbeAttributes struct {
	ResponseBody       string `dynamo:"ResponseBody"`
	ResponseTimeMillis int    `dynamo:"ResponseTimeMillis"`
	ResponseStatus     int    `dynamo:"ResponseStatus"`
}

type WebiteAttributes struct {
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
	// TODO
	return []WebsiteProbes{}
}
