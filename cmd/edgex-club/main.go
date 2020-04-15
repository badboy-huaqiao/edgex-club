// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package main

import (
	"edgex-club/internal"
	"edgex-club/internal/config"
	"edgex-club/internal/core"
	"edgex-club/internal/repository"
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {

	confpath := flag.String("confpath", "res/configuration.toml", "配置文件路径")
	isProd := flag.Bool("prod", false, "如果不是生产环境，默认不会启动TLS服务")
	flag.Parse()

	err := config.InitConfig(*confpath)
	if err != nil {
		log.Println("Cann't parse config file, exit!")
		return
	}

	r := internal.InitRouter()

	b := repository.DBConnect()
	if b {
		log.Println("connect to db success!")
	} else {
		log.Println("failed connect to db!")
	}

	//用户访问限制功能，定时清除3分钟内已经被锁定的用户，
	//防止map缓存越过内存边界
	go core.CleanupVisitors()

	if *isProd {
		go func() {
			// cer, err := tls.LoadX509KeyPair("./env/edgex-club-nginx.crt", "./env/edgex-club-nginx.key")
			// if err != nil {
			// 	log.Println(err)
			// 	return
			// }
			// tlsConfig := &tls.Config{Certificates: []tls.Certificate{cer}}
			server := &http.Server{
				Handler: core.GeneralFilter(core.Limit(r)),
				Addr:    ":443",
				//TLSConfig:    tlsConfig,
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}
			log.Println("TLS Server Listen On Port: 443")
			log.Fatal(server.ListenAndServeTLS(config.Config.Certificate.Crt, config.Config.Certificate.Key))
		}()
		log.Println("Server Listen On Port: 8080")
		if err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "https://www.edgexfoundry.club"+r.RequestURI, http.StatusMovedPermanently)
		})); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	} else {

		server := &http.Server{
			Handler:      core.GeneralFilter(core.Limit(r)),
			Addr:         ":8080",
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}
		log.Println("Dev Server Listen On Port: 8080")
		log.Fatal(server.ListenAndServe())
	}

}
