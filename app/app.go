package app

import (
	"kong/config"
	"kong/resources"
)

type Application struct {
	Config    config.Config
	Resources Resources
}

type Resources struct {
	ElasticSearch resources.ElasticSearch
	Kafka         resources.Kafka
}

func NewApplication(config config.Config) Application {
	return Application{
		Config: config,
		Resources: Resources{
			ElasticSearch: resources.ElasticSearch{},
			Kafka:         resources.Kafka{},
		},
	}
}

func (app *Application) Boot() {}

func (app *Application) Run() {}
