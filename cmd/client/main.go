package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func getUsername() {
	// Modify input username to be unix compatible
	// if result is available, return
	// If not, append random 4 digit number then return
	//return true
}

func getHosts() {
	// Modify input username to be unix compatible
	// if result is available, return
	// If not, append random 4 digit number then return
	//return true
}

func createAccount(logger *log.Logger, host string, name string, key string) {

	type AccountData struct {
		Shell   string   `json:"shell"`
		SshKeys []string `json:"ssh_keys"`
	}

	type AccountBody struct {
		Name string      `json:"name"`
		Host string      `json:"host"`
		Data AccountData `json:"data"`
	}

	api_url := fmt.Sprintf("%s/passwd", os.Getenv("API_URL"))
	api_token := os.Getenv("API_TOKEN")
	jsonData, err := json.Marshal(AccountBody{
		Name: name,
		Host: host,
		Data: AccountData{
			Shell:   "/bin/bash",
			SshKeys: []string{key},
		},
	})

	logger.Println("<- ", string(jsonData))

	req, _ := http.NewRequest("POST", api_url, bytes.NewBuffer(jsonData))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", api_token))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	logger.Println("-> ", string(body))
	resp.Body.Close()

	// TODO: trigger rendering of results and exit button
}

type writer struct {
	io.Writer
	timeFormat string
}

func (w writer) Write(b []byte) (n int, err error) {
	return w.Writer.Write(append([]byte(time.Now().Format(w.timeFormat)), b...))
}

func main() {

	fd := os.NewFile(3, "/proc/self/fd/3")
	defer fd.Close()

	logger := log.New(&writer{fd, "2006/01/02 15:04:05 "}, "[client] ", 0)

	if os.Getenv("KEY") == "none" {
		fmt.Fprintln(
			os.Stderr,
			"\nError: Public key authentication required\n",
			"\nFor help generating a key try:\n",
			"\n$ ssh-keygen -t ed25519 -f \"$HOME/.ssh/id_ed25519\"\n",
		)
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
		form.AddDropDown("Server",
			[]string{"te1.hashbang.sh", "te2.hashbang.sh"}, 0, nil,
		)
		form.AddInputField("User Name",
			os.Getenv("USER"), 33, tview.InputFieldMaxLength(30), nil,
		)
		form.AddInputField("Public Key",
			os.Getenv("KEY"), 33, tview.InputFieldMaxLength(800), nil,
		)
		form.AddButton("Create", func() {
			server_dropdown := form.GetFormItem(0).(*tview.DropDown)
			_, server := server_dropdown.GetCurrentOption()
			createAccount(
				logger,
				server,
				form.GetFormItem(1).(*tview.InputField).GetText(),
				form.GetFormItem(2).(*tview.InputField).GetText(),
			)
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
