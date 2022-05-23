package common

import (
	"testing"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

func Test_PlatformConfig_Init(t *testing.T) {
	text := `
name: local
config:
  clientRepo:
    host: localpost
    port: 3306
    username: root
    password: redhat
    dbName: auth_user_client`

	var (
		m   map[string]interface{}
		err error
	)
	err = yaml.Unmarshal([]byte(text), &m)
	if err != nil {
		t.Errorf("unmarshalling should not return error: %v", err)
		return
	}

	var p PlatformConfig
	if err = mapstructure.Decode(m, &p); err != nil {
		t.Errorf("platform config decoding should not return err: %v", err)
		return
	}

	if err = p.Init(); err != nil {
		t.Errorf("platform config initialization should not return error: %v", err)
		return
	}
}

type TestConfig struct {
	Url string `mapstructure:"url"`
}
