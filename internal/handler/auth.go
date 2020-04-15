// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package handler

import (
	"bytes"
	"edgex-club/internal/authorization"
	"edgex-club/internal/config"
	"edgex-club/internal/model"
	"edgex-club/internal/repository"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	mux "github.com/gorilla/mux"
)

type ReturnLoginUserToPageData struct {
	UserInfo    string
	Token       string
	UserPrePage string
}
type GithubUserInfo struct {
	Id         int64
	Login      string
	Avatar_url string
}

func ValidToken(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	token := vars["token"]

	var isVaild string
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return config.Config.Service.JWTKey, nil
	})

	//包括超时、被篡改等，都会无效
	if err != nil || !jwtToken.Valid {
		isVaild = "0" //无效
	} else {
		isVaild = "1" //有效
	}

	w.Header().Set("Content-Type", "text/plain;charset=utf-8")
	w.Write([]byte(isVaild))
}

func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	m := make(map[string]string)
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	name := m["name"]
	pwd := m["password"]

	u := model.User{Name: name, Password: pwd}
	ok, err := repository.UserRepos.Exists(u)

	if err != nil {
		log.Println("User: " + name + " login failed : " + err.Error())
		w.Write([]byte("log failed : " + err.Error()))
		return
	}

	if ok {
		var creds model.Credentials

		log.Println("User: " + name + " login.")
		creds.AvatarUrl = "https://avatars1.githubusercontent.com/u/42457890?v=4"
		creds.Id = "5bc0081dcedad5121dccebff"
		creds.Name = name

		token, err := authorization.NewToken(creds)
		if err != nil {
			log.Println("生成token失败！")
			w.WriteHeader(http.StatusBadRequest)
		}

		mm := make(map[string]interface{})
		mm["token"] = token
		mm["userInfo"] = creds
		result, _ := json.Marshal(mm)
		http.SetCookie(w, &http.Cookie{
			Name:  "edgex-club-token",
			Value: token,
			//Expires: expirationTime,
		})
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write(result)
	} else {
		log.Println("User: " + name + " login failed : ")
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		w.Write([]byte("no user "))
	}

}

func LoginByGitHub(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := r.URL.Query()
	code := vars["code"][0]
	userPrePage := vars["state"][0]

	githubtoken := getGithubTokenByCode(code)
	userInfo := getUserInfoByToken(githubtoken)

	var githubUserInfo = GithubUserInfo{}
	jsonStr := bytes.NewReader([]byte(userInfo))
	json.NewDecoder(jsonStr).Decode(&githubUserInfo)

	userName := githubUserInfo.Login
	userId := strconv.FormatInt(githubUserInfo.Id, 10)
	avatarUrl := githubUserInfo.Avatar_url

	u := model.User{Name: userName, GitHubId: userId, AvatarUrl: avatarUrl}
	ok, _ := repository.UserRepos.ExistsByGitHub(u)
	if !ok {
		repository.UserRepos.Insert(u)
		log.Println("user not exist, it's a new user,save to db.")
	}
	u = repository.UserRepos.FindOneByName(u.Name)

	var creds model.Credentials

	token, err := authorization.NewToken(creds)
	if err != nil {
		log.Println("生成token失败！")
		w.WriteHeader(http.StatusBadRequest)
	}

	log.Printf("User login: %v", u)

	creds.Name = u.Name
	creds.AvatarUrl = u.AvatarUrl
	creds.Id = u.Id.Hex()
	credsByte, err := json.Marshal(creds)

	t, _ := template.ParseFiles("static/redirect.html")
	data := ReturnLoginUserToPageData{
		UserInfo:    string(credsByte),
		Token:       token,
		UserPrePage: userPrePage,
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "edgex-club-token",
		Value: token,
		//Expires: expirationTime,
	})

	t.Execute(w, data)
}

func getGithubTokenByCode(code string) string {
	url := "https://github.com/login/oauth/access_token"
	param := make(map[string]string, 10)
	param["client_id"] = "173d78b242d4fc35aca9"
	param["client_secret"] = "a35c510a599f19c6041325bcb7b3579072eb9228"
	param["code"] = code

	bytesData, err := json.Marshal(param)
	if err != nil {
		log.Println("param json faild!")
	}
	jsonStr := bytes.NewReader(bytesData)
	request, _ := http.NewRequest("POST", url, jsonStr)

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("request github to get code faild!")
	}
	// defer response.Body.Close()
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("respBytes  faild!")
	}
	//access_token=11f5681361f6448411abcb33401641fc61d9d6e4&scope=user&token_type=bearer
	tokenStr := string(respBytes)
	args := strings.Split(tokenStr, "&")
	tokenInfo := strings.Split(args[0], "=")
	token := tokenInfo[1]

	return token
}

func getUserInfoByToken(token string) string {
	url := "https://api.github.com/user?access_token=" + token + "&scope=user"
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, _ := client.Do(req)
	// defer resp.Body.Close()
	userData, _ := ioutil.ReadAll(resp.Body)
	userInfo := string(userData)
	//log.Println(userInfo)
	return userInfo
}
