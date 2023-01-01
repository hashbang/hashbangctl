package main

import (
	"bytes"
	"encoding/json"
	"encoding/base64"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
    "strings"
)

type RequestBody struct {
	Name                string      `json:"name"`
	Host                string      `json:"host"`
	Shell               string      `json:"shell"`
	Keys                []string    `json:"keys"`
}

type ResponseBody struct {
	Hint                string      `json:"hint"`
	Details             string      `json:"details"`
	Message             string      `json:"message"`
	Code                string      `json:"code"`
	Request             RequestBody `json:"request"`
}

type SshPublicKey struct {
	Fingerprint         string      `json:"fingerprint"`
    Base64Fingerprint   string      `json:"base64_fingerprint"`
	Type                string      `json:"type"`
	Key                 string      `json:"key"`
	Comment             string      `json:"comment"`
	Uid                 int         `json:"uid"`
}

type User struct {
	Uid                 int         `json:"uid"`
	Name                string      `json:"name"`
    Host                string      `json:"host"`
	Type                string      `json:"type"`
	Key                 string      `json:"key"`
	Comment             string      `json:"comment"`
    Shell               string      `json:"shell"`
}

type Host struct {
	Name     string `json:"name"`
	MaxUsers int    `json:"maxusers"`
}

func getHosts() ([]string, error) {

	var hostResponse []Host
	hosts := []string{}

	apiUrl := fmt.Sprintf("%s/hosts", os.Getenv("API_URL"))
	apiToken := os.Getenv("API_TOKEN")
	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return hosts, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return hosts, err
	}
	err = json.Unmarshal(body, &hostResponse)
	if err != nil {
		return hosts, err
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))
	for _, i := range r.Perm(len(hostResponse)) {
		hosts = append(hosts, hostResponse[i].Name)
	}
	return hosts, nil
}

func getKeys(
	key string,
) ([]SshPublicKey, error) {
	var sshPublicKeys []SshPublicKey
    keyStripped := strings.Fields(key)[1]
    keyDecoded, err := base64.StdEncoding.DecodeString(keyStripped)
	if err != nil {
		return sshPublicKeys, err
	}
    fingerprint := sha256.Sum256(keyDecoded)
	apiUrl := fmt.Sprintf(
        "%s/ssh_public_key?fingerprint=ilike.%x",
        os.Getenv("API_URL"),
        fingerprint,
    )
	apiToken := os.Getenv("API_TOKEN")
	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return sshPublicKeys, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return sshPublicKeys, err
	}
	err = json.Unmarshal([]byte(body), &sshPublicKeys)
	if err != nil {
		return sshPublicKeys, err
    }
    return sshPublicKeys, nil
}

func getUsersById(
	uid int,
) ([]User, error) {
	var users []User
	apiUrl := fmt.Sprintf(
        "%s/passwd?uid=eq.%d",
        os.Getenv("API_URL"),
        uid,
    )
	apiToken := os.Getenv("API_TOKEN")
	req, _ := http.NewRequest("GET", apiUrl, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return users, err
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return users, err
	}
	err = json.Unmarshal([]byte(body), &users)
	if err != nil {
		return users, err
    }
    return users, nil
}

func createAccount(
	logger *log.Logger,
	host string,
	name string,
	key string,
) error {
	apiUrl := fmt.Sprintf("%s/signup", os.Getenv("API_URL"))
	apiToken := os.Getenv("API_TOKEN")
	requestBody := RequestBody{
		Name:  name,
		Host:  host,
		Shell: "/bin/bash",
		Keys:  []string{key},
	}
	jsonData, err := json.Marshal(requestBody)
	logger.Println("[client] ??", string(jsonData))
	req, _ := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 201 {
		logger.Println("[client] ++", string(jsonData))
		return nil
	}
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	var responseBody ResponseBody
	err = json.Unmarshal(body, &responseBody)
	if err != nil {
		return err
	}
	responseBody.Request = requestBody
	jsonError, err := json.Marshal(responseBody)
	logger.Println("[client] !!", string(jsonError))
	return errors.New(responseBody.Message)
}
