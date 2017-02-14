package plugins

import (
	"encoding/json"
	"log"

	"gopkg.in/olivere/elastic.v3"

	"github.com/InnovaCo/serve/manifest"
)

func init() {
	manifest.PluginRegestry.Add("dashboard.create.kibana", DashboardCreateKibana{})

}

type DashboardSettings struct {
	User      string `json:"user"`
	Group     string `json:"group"`
	Title     string `json:"title"`
	Dashboard string `json:"dashboard"`
}

type DashboardCreateKibana struct {
}

func (p DashboardCreateKibana) Run(data manifest.Manifest) error {
	if data.GetBool("purge") {
		return p.Drop(data)
	}
	return p.Create(data)
}

func (p DashboardCreateKibana) Create(data manifest.Manifest) error {
	log.Println(data.GetString("url"))

	client, err := elastic.NewClient(elastic.SetURL(data.GetString("url")))
	if err != nil {
		return err
	}

	ds := DashboardSettings{
		User:      data.GetStringOr("user", "guest"),
		Group:     data.GetStringOr("group", "guest"),
		Title:     data.GetString("title"),
		Dashboard: data.GetTree("dashboard").String()}

	_, err = client.Index().
		Index("kibana-int").
		Type("dashboard").
		Id(data.GetString("title")).
		BodyJson(&ds).
		Do()
	return err
}

func (p DashboardCreateKibana) Drop(data manifest.Manifest) error {
	client, err := elastic.NewClient(elastic.SetURL(data.GetString("url")))
	if err != nil {
		return err
	}
	_, err = client.Delete().
		Index("kibana-int").
		Type("dashboard").
		Id(data.GetString("title")).
		Do()
	return err
}
