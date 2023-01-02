package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"log"
	"os"
)

var logoText = `
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

`

func createForm(
	logger *log.Logger,
	hosts []string,
) {
	app := tview.NewApplication()
	logo := tview.NewTextView()
	logo.SetTextAlign(1)
	logo.SetText(logoText)
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
			err := createUser(logger, server, user, key)
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

func editForm(
	logger *log.Logger,
	hosts []string,
    user User,
    sshPublicKeys []SshPublicKey,
) {
	app := tview.NewApplication()
	logo := tview.NewTextView()
	logo.SetTextAlign(1)
	logo.SetText(logoText)
	title := tview.NewTextView()
	title.SetTextAlign(1)
    title.SetText(fmt.Sprintf("Editing User: \"%s\"",user.Name))
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
		form.AddInputField("Shell",
			user.Shell, 33, tview.InputFieldMaxLength(800), nil,
		)
        keyNum := len(sshPublicKeys)
        for i:=0; i < len(sshPublicKeys); i++{
            key := sshPublicKeys[i]
            keyString := fmt.Sprintf("%s %s",key.Type, key.Key)
		    form.AddInputField(fmt.Sprintf("Public Key %d",i+1),
		    	keyString, 33, tview.InputFieldMaxLength(800), nil,
		    )
        }
		form.AddButton("Add Key", func(){
            keyNum = keyNum + 1
		    form.AddInputField(fmt.Sprintf("Public Key %d",keyNum),
		    	"", 33, tview.InputFieldMaxLength(800), nil,
		    )
        })
		form.AddButton("Update", func() {
			//server_dropdown := form.GetFormItem(0).(*tview.DropDown)
			//_, server := server_dropdown.GetCurrentOption()
			//shell := form.GetFormItem(1).(*tview.InputField).GetText()
			err := editUser(logger, user, sshPublicKeys)
			if err != nil {
				app.Stop()
				fmt.Fprintln(
					os.Stderr,
					"\nError: Account update failed\n",
					fmt.Errorf("\n%v\n", err),
				)
                //TODO: update User and sshPublicKeys structs based on input
			    fmt.Fprintln(os.Stdout,"\nUser: ",user, sshPublicKeys)
				os.Exit(1)
			}
			app.Stop()
			fmt.Fprintln(
				os.Stdout,
				"\nAccount update successful!\n",
				"\nYou can connect to your account via:\n",
				fmt.Sprintf("\n$ ssh %s@%s\n", user.Name, user.Host),
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
			AddItem(title, 1, 1, false).
			AddItem(frame, 14, 1, true).
			AddItem(tview.NewBox(), 0, 1, false),
		50, 2, true)
	flex.AddItem(tview.NewBox(), 0, 1, false)
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
