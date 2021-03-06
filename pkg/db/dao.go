package db

import (
	"fmt"
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

func (d *ProbesDao) DeleteWebsite(url string) error {
	t := d.DB.Table(WebsiteProbesTable)

	error := t.Delete("PK", WebsitePK).Range("SK", WebsiteSKPrefix+url).Run()

	return error
}

func (d *ProbesDao) GetWebsites() []WebsiteProbes {
	t := d.DB.Table(WebsiteProbesTable)

	var websites []WebsiteProbes
	t.Get("PK", WebsitePK).
		Range("SK", dynamo.BeginsWith, WebsiteSKPrefix).
		All(&websites)

	return websites
}

func (d *ProbesDao) GetWebsite(url string) *WebsiteProbes {
	t := d.DB.Table(WebsiteProbesTable)

	var website *WebsiteProbes
	t.Get("PK", WebsitePK).
		Range("SK", dynamo.Equal, WebsiteSKPrefix+url).
		One(&website)

	return website
}

func (d *ProbesDao) PutWebsite(url string, probeScheduleMinutes int) (*WebsiteProbes, error) {
	t := d.DB.Table(WebsiteProbesTable)

	newWebsite := &WebsiteProbes{
		PK: WebsitePK,
		SK: WebsiteSKPrefix + url,

		WebsiteAttributes: WebsiteAttributes{
			ProbeScheduleMinutes: probeScheduleMinutes,
			Active:               true,
			CreatedDate:          time.Now(),
		},
	}

	err := t.Put(newWebsite).Run()

	if err != nil {
		return nil, err
	}

	return newWebsite, nil
}

func (d *ProbesDao) GetWebsiteProbes(url string) []WebsiteProbes {
	t := d.DB.Table(WebsiteProbesTable)

	var websitesProbes []WebsiteProbes
	t.Get("PK", fmt.Sprintf("%s%s#PROBE", WebsitePK, url)).
		All(&websitesProbes)

	return websitesProbes
}

// TODO fix this as it times out
func (d *ProbesDao) PutWebsiteProbe(url, responseBody string, responseStatus, responseTimeMillis int) (*WebsiteProbes, error) {
	t := d.DB.Table(WebsiteProbesTable)

	newWebsiteProbe := &WebsiteProbes{
		PK: fmt.Sprintf("%s%s#PROBE", WebsitePK, url),
		SK: time.Now().Format(time.RFC3339),

		ProbeAttributes: ProbeAttributes{
			ResponseBody:       responseBody,
			ResponseTimeMillis: responseTimeMillis,
			ResponseStatus:     responseStatus,
		},
	}

	err := t.Put(newWebsiteProbe).Run()

	if err != nil {
		return nil, err
	}

	return newWebsiteProbe, nil
}
