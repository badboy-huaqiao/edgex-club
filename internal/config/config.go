// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package config

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type service struct {
	Host                      string   `toml:"Host"`
	Port                      int      `toml:"Port"`
	Labels                    []string `toml:"Labels"`
	StartupMsg                string   `toml:"StartupMsg"`
	DefaultStaticResourcePath string   `toml:"DefaultStaticResourcePath"`
}

type database struct {
	Host         string `toml:"Host"`
	Port         int    `toml:"Port"`
	DatabaseName string `toml:"DatabaseName"`
	Username     string `toml:"Username"`
	Password     string `toml:"Password"`
}

type certificate struct {
	Crt string `toml:"Crt"`
	Key string `toml:"Key"`
}

type github struct {
	ClientId string `toml:"ClientId"`
	Secret   string `toml:"Secret"`
}

type configuration struct {
	ServiceConf     service     `toml:"Service"`
	DatabaseConf    database    `toml:"Database"`
	CertificateConf certificate `toml:"Certificate"`
	EnvConf         env         `toml:"Env"`
}

type env struct {
	Prod bool
}

var config *configuration

func Conf() *configuration {
	return config
}

func (c *configuration) Service() service {
	return c.ServiceConf
}

func (c *configuration) Database() database {
	return c.DatabaseConf
}

func (c *configuration) DBAddr() string {
	return fmt.Sprintf("%s:%d", c.Database().Host, c.Database().Port)
}

func (c *configuration) Crt() certificate {
	return c.CertificateConf
}

func (c *configuration) Env() env {
	return c.EnvConf
}

func InitConfig(path string) error {
	c := &configuration{}
	absPath, _ := filepath.Abs(path)
	if _, err := toml.DecodeFile(absPath, c); err != nil {
		log.Printf("Decode config file error: %s\n", err.Error())
		return err
	}
	config = c
	return nil
}
