package main

import (
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type pushFiles struct {
	PushFiles []string `yaml:"pushFiles"`
}

func readPushFiles(fileName string) ([]string, error) {
	bs, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read from file")
	}

	pf := &pushFiles{PushFiles: make([]string, 0)}
	err = yaml.Unmarshal(bs, pf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal yaml")
	}

	return pf.PushFiles, nil
}
