package conf

import (
	"github.com/bagaking/vaga.go/localVideos"
	"io/ioutil"
	"log"
	"os"
	"path"
	"github.com/json-iterator/go"
)

func init() {
	loadConfig()
}

const configName = "vaga.json"

var Static VagaConfig

func loadConfig() {
	Static = VagaConfig{[]localVideos.VideoBlobConf{} }

	runningPth, err := os.Getwd()
	if err != nil {
		log.Println(err)
		runningPth = "."
	}

	path := path.Join(runningPth, configName)

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("cannot find config file", runningPth)
		b, err = ioutil.ReadFile(path + ".sample")
		panic(err)
	}

	err = jsoniter.Unmarshal(b, &Static)
	if err != nil {
		panic(err)
	}
}

