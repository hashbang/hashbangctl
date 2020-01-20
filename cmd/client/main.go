package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os"
)

func getUsername() {
	// Modify input username to be unix compatible
	// if result is available, return
	// If not, append random 4 digit number then return
	//return true
}

func createAccount(server string, user string, key string) {
	// Create account
	// trigger rendering of results and exit button
	//log.Println(server, user, key)
}

func main() {
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
			[]string{"de1.hashbang.sh", "la1.hashbang.sh"}, 0, nil,
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
