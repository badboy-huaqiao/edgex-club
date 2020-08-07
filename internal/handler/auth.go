// Copyright © 2020-2021 https://www.edgexfoundry.club. All Rights Reserved.
// SPDX-License-Identifier: GPL-2.0

package handler

import (
	"bytes"
	"edgex-club/internal/authorization"
	"edgex-club/internal/config"
	"edgex-club/internal/model"
	repo "edgex-club/internal/repository"
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

func LoginByGitHubCallback(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := r.URL.Query()
	code := vars["code"][0]
	userPrePage := vars["state"][0]

	var githubtoken string
	var err error
	var githubUserInfo GithubUserInfo
	if githubtoken, err = getGithubTokenByCode(code); err != nil {
		fmt.Printf("Access github error: %s\n", err.Error())
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	if githubUserInfo, err = getUserInfoByToken(githubtoken); err != nil {
		fmt.Printf("Access github error: %s\n", err.Error())
		http.Redirect(w, r, "/error", http.StatusTemporaryRedirect)
		return
	}

	user := model.User{
		Name:      githubUserInfo.Login,
		GitHubId:  strconv.FormatInt(githubUserInfo.Id, 10),
		AvatarUrl: githubUserInfo.Avatar_url,
	}
	u, err := repo.UserRepositoryClient().FetchOneByGitHub(user.GitHubId)
	if err != nil {
		repo.UserRepositoryClient().Add(user)
		log.Println("user not exist, it's a new user,save to db.")
	} else if u.AvatarUrl != user.AvatarUrl {
		repo.UserRepositoryClient().Update(user)
		log.Println("user existed, update user's avatar.")
	}

	u, _ = repo.UserRepositoryClient().FetchOneByName(user.Name)

	creds := model.Credentials{
		Name:      u.Name,
		AvatarUrl: u.AvatarUrl,
		Id:        u.Id.Hex(),
	}

	token, err := authorization.NewToken(creds)
	if err != nil {
		log.Printf("生成token失败！err: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("User [id=%s,githubId=%s,avatarUrl=%s] login.\n", u.Id.Hex(), u.GitHubId, u.AvatarUrl)

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
		// Domain:   "edgexfoundry.club",
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		//Path:     "/",
		Expires: time.Now().Add(7 * 24 * time.Hour),
		MaxAge:  int(time.Now().Add(7*24*time.Hour).UnixNano() / 1000000),
	})
	http.Redirect(w, r, userPrePage, http.StatusTemporaryRedirect)
}

func getGithubTokenByCode(code string) (string, error) {
	url := "https://github.com/login/oauth/access_token"
	reqBody := map[string]string{
		"client_id":     config.Conf().GitHub().ClientId,
		"client_secret": config.Conf().GitHub().Secret,
		"code":          code,
	}

	body, _ := json.Marshal(reqBody)
	request, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	request.Header.Set(ContentType, ContentTypeJSON)
	client := &http.Client{Timeout: 30 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("request github to get code faild: %s\n", err.Error())
		return "", err
	}
	defer response.Body.Close()
	respBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read github response body faild: %s\n", err.Error())
		return "", err
	}

	//tokenStr access_token=11f5681361f6448411abcb33401641fc61d9d6e4&scope=user&token_type=bearer
	tokenStr := string(respBytes)
	args := strings.Split(tokenStr, "&")
	tokenInfo := strings.Split(args[0], "=")
	token := tokenInfo[1]

	return token, nil
}

func getUserInfoByToken(token string) (GithubUserInfo, error) {
	url := "https://api.github.com/user?access_token=" + token + "&scope=user"
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	var githubUserInfo GithubUserInfo
	if err != nil {
		log.Printf("Retrieve github User Infor failed: %s", err.Error())
		return githubUserInfo, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&githubUserInfo); err != nil {
		log.Printf("Can't read  github response body: %s", err.Error())
		return githubUserInfo, err
	}

	return githubUserInfo, nil
}
