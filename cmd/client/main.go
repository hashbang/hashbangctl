package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type RequestBody struct {
	Name  string   `json:"name"`
	Host  string   `json:"host"`
	Shell string   `json:"shell"`
	Keys  []string `json:"keys"`
}

type ResponseBody struct {
	Hint    string      `json:"hint"`
	Details string      `json:"details"`
	Message string      `json:"message"`
	Code    string      `json:"code"`
	Request RequestBody `json:"request"`
}

type Host struct {
	Name     string `json:"name"`
	MaxUsers int    `json:"maxusers"`
}

func getExisting() {
	// TODO:
	// hash public key
	// lookup pubkey via userdb/postgrest computed column
	// kick user out if their key is already in use
}

func getUsername() {
	// TODO:
	// curl -i "https://userdb.hashbang.sh/passwd?name=ilike.lrvick*&select=name&limit10"
	// Modify input username to be unix compatible
	// if result is available, return
	// If not, append random 4 digit number then return
	//return true
}

func getHosts() ([]string, error) {

	var hostResponse []Host
	hosts := []string{}

	apiUrl := fmt.Sprintf("%s/hosts", os.Getenv("API_URL"))
	//TODO: implement API tokens on prod
	// apiToken := os.Getenv("API_TOKEN")
	req, _ := http.NewRequest("GET", apiUrl, nil)
	//TODO: implement API tokens on prod
	// req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))
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

func createAccount(
	logger *log.Logger,
	host string,
	name string,
	key string,
) error {
	apiUrl := fmt.Sprintf("%s/signup", os.Getenv("API_URL"))
	//TODO: implement API tokens on prod
	//apiToken := os.Getenv("API_TOKEN")
	requestBody := RequestBody{
		Name:  name,
		Host:  host,
		Shell: "/bin/bash",
		Keys:  []string{key},
	}
	jsonData, err := json.Marshal(requestBody)
	logger.Println("[client] ??", string(jsonData))
	req, _ := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
	//TODO: implement API tokens on prod
	//req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))
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

func main() {

	fd := os.NewFile(3, "/proc/self/fd/3")
	defer fd.Close()
	logger := log.New(fd, "", log.Ldate|log.Ltime)

	if os.Getenv("KEY") == "none" {
		fmt.Fprintln(
			os.Stderr,
			"\nError: Public key authentication required\n",
			"\nFor help generating a key try:\n",
			"\n$ ssh-keygen -t ed25519 -f \"$HOME/.ssh/id_ed25519\"\n",
		)
		os.Exit(1)
	}

	hosts, err := getHosts()
	if err != nil {
		fmt.Fprintln(os.Stderr, "\nError: Unable to get host list")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	app := tview.NewApplication()

	logo := tview.NewTextView()
	logo.SetTextAlign(1)
	logo.SetText(`
     █████   █████       █████
     █████   █████       █████
     █████   █████       █████
███████████████████████  █████
███████████████████████  █████
     █████   █████       █████
     █████   █████       █████
███████████████████████  █████
███████████████████████  █████
     █████   █████　　　　　　
     █████   █████       █████
     █████   █████       █████

`)

	frame := tview.NewFrame(func() tview.Primitive {
		form := tview.NewForm()
		form.SetLabelColor(tcell.ColorWhite)
		form.SetItemPadding(2)
		form.SetFieldTextColor(tcell.ColorGray)
		form.SetButtonTextColor(tcell.ColorGray)
		form.SetFieldBackgroundColor(tcell.ColorWhite)
		form.SetButtonBackgroundColor(tcell.ColorWhite)
		form.SetBorder(false)
		form.SetButtonsAlign(1)
		form.AddDropDown("Server", hosts, 0, nil)
		// TODO: check username is available. Append numbers if needed
		form.AddInputField("User Name",
			os.Getenv("USER"), 33, tview.InputFieldMaxLength(30), nil,
		)
		form.AddInputField("Public Key",
			os.Getenv("KEY"), 33, tview.InputFieldMaxLength(800), nil,
		)
		form.AddButton("Create", func() {
			server_dropdown := form.GetFormItem(0).(*tview.DropDown)
			_, server := server_dropdown.GetCurrentOption()
			user := form.GetFormItem(1).(*tview.InputField).GetText()
			key := form.GetFormItem(2).(*tview.InputField).GetText()
			err := createAccount(logger, server, user, key)
			if err != nil {
				app.Stop()
				fmt.Fprintln(
					os.Stderr,
					"\nError: Account creation failed\n",
					fmt.Errorf("\n%v\n", err),
				)
				os.Exit(1)
			}
			app.Stop()
			fmt.Fprintln(
				os.Stdout,
				"\nAccount creation successful!\n",
				"\nYou can now connect to your account via:\n",
				fmt.Sprintf("\n$ ssh %s@%s\n", user, server),
			)
			os.Exit(1)
		})
		form.AddButton("Exit", app.Stop)
		return form
	}())
	frame.SetBorder(false)

	flex := tview.NewFlex()
	flex.AddItem(tview.NewBox(), 0, 1, false)
	flex.AddItem(
		tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(tview.NewBox(), 0, 1, false).
			AddItem(logo, 14, 1, false).
			AddItem(frame, 14, 1, true).
			AddItem(tview.NewBox(), 0, 1, false),
		50, 2, true)
	flex.AddItem(tview.NewBox(), 0, 1, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}

}
