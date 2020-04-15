// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package core

import (
	"bytes"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Visitor struct {
	Limiter    *rate.Limiter
	LastAccess time.Time
	//ForceTryCount int
}

//设定IP的结构，用于判断IP是合法格式
type IpRange struct {
	Start net.IP
	End   net.IP
}

var Visitors = make(map[string]*Visitor)
var mtx sync.Mutex

//枚举常规局域网私有IP
var privateRanges = []IpRange{
	IpRange{
		Start: net.ParseIP("10.0.0.0"),
		End:   net.ParseIP("10.255,255,255"),
	},
	IpRange{
		Start: net.ParseIP("100.64.0.0"),
		End:   net.ParseIP("100.127.255.255"),
	},
	IpRange{
		Start: net.ParseIP("172.16.0.0"),
		End:   net.ParseIP("172.31.255.255"),
	},
	IpRange{
		Start: net.ParseIP("192.0.0.0"),
		End:   net.ParseIP("192.0.0.255"),
	},
	IpRange{
		Start: net.ParseIP("192.168.0.0"),
		End:   net.ParseIP("192.168.255.255"),
	},
	IpRange{
		Start: net.ParseIP("192.18.0.0"),
		End:   net.ParseIP("192.19.255.255"),
	},
}

//判断ip是否在合法的常规局域网IP段范围内
func inRange(r IpRange, ipAddress net.IP) bool {
	if bytes.Compare(ipAddress, r.Start) >= 0 &&
		bytes.Compare(ipAddress, r.End) < 0 {
		return true
	}
	return false
}

//判断IP是否是局域网常规的IP
func isPrivateSubnet(ipAddress net.IP) bool {
	//当前阶段只考虑IPv4的情况，IPv6以后再更新
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		for _, r := range privateRanges {
			if inRange(r, ipAddress) {
				return true
			}
		}
	}

	return false
}

func getRealIP(r *http.Request) string {
	for _, h := range []string{"X-Forwarded-For", "X-Real-Ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		//采用从数组的右边开始判断，知道得到一个公网IP，
		//防止有人借助请求头恶意添加一个虚拟的公网IP
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				continue
			}
			return ip
		}
	}

	return ""
}

//每个请求每秒允许消费2个令牌，总共最大为5个令牌
//var limiter = rate.NewLimiter(2,5)
func Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//log.Println("limit=====")
		v := getVisitor(r.RemoteAddr)
		log.Println("user [ " + r.RemoteAddr + " ] access.")

		// log.Println("user [ " + getRealIP(r) + " ] access.")
		if !v.Limiter.Allow() {
			time.Sleep(2 * time.Second)
			log.Println("user [ " + r.RemoteAddr + " ] access are too frequent.")
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func addVisitor(ip string) *Visitor {
	limiter := rate.NewLimiter(5, 8)
	mtx.Lock()
	v := &Visitor{Limiter: limiter, LastAccess: time.Now()}
	Visitors[ip] = v
	mtx.Unlock()
	return v
}

func getVisitor(ip string) *Visitor {
	mtx.Lock()
	v, exist := Visitors[ip]
	if !exist {
		mtx.Unlock()
		return addVisitor(ip)
	}
	//更新一次访问时间
	v.LastAccess = time.Now()
	mtx.Unlock()
	return v
}

//后台运行，防止map越来越大
//每分钟检查一次当前visitors缓存，并删除超过3分钟都没有访问过的用户
func CleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		mtx.Lock()
		for ip, v := range Visitors {
			if time.Now().Sub(v.LastAccess) > 3*time.Minute {
				delete(Visitors, ip)
			}
		}
		mtx.Unlock()
	}
}
