// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package main

import (
	"edgex-club/internal"
	"edgex-club/internal/config"
	"edgex-club/internal/core"
	"edgex-club/internal/repository"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	confpath := flag.String("confpath", "res/configuration.toml", "配置文件路径")
	flag.Parse()
	if err := config.InitConfig(*confpath); err != nil {
		log.Println("Cann't parse config file, exit!")
		return
	}
	if err := repository.DBConnect(); err != nil {
		log.Fatalf("connect to db error: %s\n", err.Error())
		return
	}
	log.Println("Success connect to db!")

	// authorization.InitAuth()

	//用户访问限制功能，定时清除3分钟内已经被锁定的用户，
	//防止map缓存越过内存边界
	// go core.CleanupVisitors()

	r := internal.InitRouter()

	if config.Conf().Env().Prod {
		go func() {
			server := &http.Server{
				Handler:      core.GeneralFilter(core.Limit(r)),
				Addr:         fmt.Sprintf(":%d", config.Conf().Service().Port),
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}
			log.Println("TLS Server Listen On Port: 443")
			log.Fatal(server.ListenAndServeTLS(config.Conf().Crt().Crt, config.Conf().Crt().Key))
		}()
		log.Println("Server Listen On Port: 80")
		if err := http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dest := fmt.Sprintf("https://www.edgexfoundry.club%s", r.RequestURI)
			http.Redirect(w, r, dest, http.StatusMovedPermanently)
		})); err != nil {
			log.Fatalf("ListenAndServe error: %s\n", err.Error())
		}
	} else {
		server := &http.Server{
			Handler:      core.GeneralFilter(r),
			Addr:         fmt.Sprintf(":%d", config.Conf().Service().Port),
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.Printf("Dev Server Listen On Port: %d\n", config.Conf().Service().Port)
		log.Fatal(server.ListenAndServe())
	}

}
