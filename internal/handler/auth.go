// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package handler

import (
	"bytes"
	"edgex-club/internal/authorization"
	"edgex-club/internal/model"
	"edgex-club/internal/repository"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ReturnLoginUserToPageData struct {
	UserInfo    string
	Token       string
	UserPrePage string
}
type GithubUserInfo struct {
	Id         int64
	Login      string //user name
	Avatar_url string
}

func genCredsUser(r *http.Request) *model.Credentials {
	credsUserStr := r.Header.Get(CredsUser)
	if credsUserStr == "" {
		return nil
	}
	var credsUser model.Credentials
	json.Unmarshal([]byte(credsUserStr), &credsUser)
	return &credsUser
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
			Name:     "edgex-club-token",
			Value:    token,
			HttpOnly: true,
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

	var githubtoken, userInfo string
	var err error
	if githubtoken, err = getGithubTokenByCode(code); err != nil {
		fmt.Printf("Access github error: %s\n", err.Error())
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}
	if userInfo, err = getUserInfoByToken(githubtoken); err != nil {
		fmt.Printf("Access github error: %s\n", err.Error())
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

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

	creds := model.Credentials{
		Name:      u.Name,
		AvatarUrl: u.AvatarUrl,
		Id:        u.Id.Hex(),
	}
	token, err := authorization.NewToken(creds)
	if err != nil {
		log.Println("生成token失败！")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("User login: %v\n", u)

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
	})
	http.Redirect(w, r, userPrePage, http.StatusTemporaryRedirect)
}

func getGithubTokenByCode(code string) (string, error) {
	url := "https://github.com/login/oauth/access_token"
	param := make(map[string]string, 10)
	param["client_id"] = "8dc598397ad0cc13bed8"
	param["client_secret"] = "5d4ceb59b837b846cfd5ce7416af5dbc8a89b241"
	param["code"] = code

	bytesData, err := json.Marshal(param)
	if err != nil {
		log.Println("param json faild!")
		return "", err
	}
	jsonStr := bytes.NewReader(bytesData)
	request, _ := http.NewRequest("POST", url, jsonStr)
	request.Header.Set(ContentType, ContentTypeJSON)
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		log.Println("request github to get code faild!")
		return "", err
	}
	defer response.Body.Close()
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("respBytes  faild!")
		return "", err
	}

	//tokenStr access_token=11f5681361f6448411abcb33401641fc61d9d6e4&scope=user&token_type=bearer
	tokenStr := string(respBytes)
	args := strings.Split(tokenStr, "&")
	tokenInfo := strings.Split(args[0], "=")
	token := tokenInfo[1]

	return token, nil
}

func getUserInfoByToken(token string) (string, error) {
	url := "https://api.github.com/user?access_token=" + token + "&scope=user"
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Retrieve User Infor from github failed: %s", err.Error())
		return "", err
	}
	defer resp.Body.Close()
	userInfoData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Can't read response body: %s", err.Error())
		return "", err
	}
	return string(userInfoData), nil
}
