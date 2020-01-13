package main

import (
    "time"
    "os"
    "fmt"
    "github.com/gdamore/tcell"
    "github.com/rivo/tview"
)

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
        AddDropDown("Server", []string{"de1.hashbang.sh", "ny1.hashbang.sh"}, 0, nil).
        AddInputField("username", os.Getenv("USER"), 150, nil, nil).
        AddInputField("ssh-key", os.Getenv("KEY"), 150, nil, nil).
        AddButton("Create Account", func() {
            // actually create account here
        }).
        AddButton("Exit", func() {
            //app.Stop()
        })

    // form.SetItemPadding(2).SetTitle("Signup for hashbang").SetTitle("Signup for hashbang")

    form.SetBorder(true).SetTitle("Sign Up for #!").SetTitleAlign(tview.AlignCenter)
    // return tview.NewBox().
    //     SetBackgroundColor(tcell.ColorBlue).
    //     SetDrawFunc(drawTime)

    return form
}

func drawTime(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
    timeStr := time.Now().Format("Current time is 15:04:05")
    tview.Print(screen, timeStr, x, height/2, width, tview.AlignCenter, tcell.ColorLime)
    return 0, 0, 0, 0
}
