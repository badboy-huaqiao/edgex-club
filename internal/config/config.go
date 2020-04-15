// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package config

import (
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Service struct {
	Host                      string   `toml:"Host"`
	Port                      int      `toml:"Port"`
	Labels                    []string `toml:"Labels"`
	StartupMsg                string   `toml:"StartupMsg"`
	DefaultStaticResourcePath string   `toml:"DefaultStaticResourcePath"`
	JWTKey                    string   `toml:"JWTKey"`
}

type Database struct {
	Host         string `toml:"Host"`
	Port         int    `toml:"Port"`
	DatabaseName string `toml:"DatabaseName"`
	Username     string `toml:"Username"`
	Password     string `toml:"Password"`
}

type Configuration struct {
	Service     Service     `toml:"Service"`
	Database    Database    `toml:"Database"`
	Certificate Certificate `toml:"Certificate"`
}

type Certificate struct {
	CertPath string `toml:"CertPath"`
	Crt      string `toml:"Crt"`
	Key      string `toml:"Key"`
}

var Config *Configuration

func InitConfig(path string) error {
	c := &Configuration{}
	absPath, _ := filepath.Abs(path)
	if _, err := toml.DecodeFile(absPath, c); err != nil {
		log.Printf("Decode Config File Error:%v", err)
		return err
	}

	Config = c

	return nil
}
