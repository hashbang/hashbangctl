package main

import (
	"time"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	_ = `
We are a diverse community of people who love teaching and learning.
Putting a #! at the beginning of a "script" style program tells a
computer that it needs to "do something" or "execute" this file.
Likewise, we are a community of people that like to "do stuff".

If you like technology, and you want to learn to write your first
program, learn to use Linux, or even take on interesting challenges
with some of the best in the industry, you are in the right place.

The following will set you up with a "shell" account on one of our
shared systems. From here you can run IRC chat clients to talk to us,
access to personal file storage and web hosting, and a wide range of
development tools.

Everything should work perfectly, unless it doesn't

Please report any issues here:
	-> https://github.com/hashbang/hashbang.sh/issues/

	`
	frame := tview.NewFrame(drawForm()).
		SetBorders(2, 2, 2, 2, 4, 4).
		AddText("Welcome to #!", true, tview.AlignLeft, tcell.ColorWhite).
		AddText("This network has three rules:", true, tview.AlignLeft, tcell.ColorRed).
		AddText("1. When people need help, teach. Don't do it for them", true, tview.AlignLeft, tcell.ColorRed).
		AddText("2. Don't use our resources for closed source projects", true, tview.AlignLeft, tcell.ColorRed).
		AddText("3. Be excellent to each other", true, tview.AlignLeft, tcell.ColorRed).
		// AddText(intro, true, tview.AlignLeft, tcell.ColorRed).
		AddText("open source everything", false, tview.AlignCenter, tcell.ColorGreen).
		AddText("for help, join irc at irc.hashbang.sh", false, tview.AlignCenter, tcell.ColorGreen)

	if err := app.SetRoot(frame, true).Run(); err != nil {
		panic(err)
	}
}

func drawForm() tview.Primitive {
	form := tview.NewForm().
		AddDropDown("Server", []string{"de1.hashbang.sh", "ny1.hashbang.sh"}, 150, nil).
		AddInputField("username", "", 150, nil, nil).
		AddInputField("ssh-key", "", 150, nil, nil).
		AddButton("Save", func() {
		}).
		AddButton("Quit", func() {
			// app.Stop()
		})

	// form.SetItemPadding(2).SetTitle("Signup for hashbang").SetTitle("Signup for hashbang")

	form.SetBorder(true).SetTitle("Sign Up for #!").SetTitleAlign(tview.AlignCenter)
	// return tview.NewBox().
	// 	SetBackgroundColor(tcell.ColorBlue).
	// 	SetDrawFunc(drawTime)

	return form
}

func drawTime(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	timeStr := time.Now().Format("Current time is 15:04:05")
	tview.Print(screen, timeStr, x, height/2, width, tview.AlignCenter, tcell.ColorLime)
	return 0, 0, 0, 0
}
