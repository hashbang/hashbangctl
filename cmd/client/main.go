package main

import (
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
        form.AddDropDown("Server", []string{"de1.hashbang.sh"}, 0, nil)
        form.AddInputField("User Name", os.Getenv("USER"), 25, nil, nil)
        form.AddInputField("Public Key", os.Getenv("KEY"), 150, nil, nil)
        form.AddButton("Create", func() {
            // actually create account here
        })
        form.AddButton("Exit", func() { app.Stop() })
        form.SetLabelColor(tcell.ColorWhite)
        form.SetItemPadding(2)
        form.SetFieldBackgroundColor(tcell.ColorGray)
        form.SetButtonTextColor(tcell.ColorWhite)
        form.SetButtonBackgroundColor(tcell.ColorGray)
        form.SetButtonsAlign(1)
        form.SetBorder(true)
        return form
    }())
    frame.SetBorder(false)
    frame.SetBorders(0, 2, 0, 0, 10, 10)

    grid := tview.NewGrid()
    grid.AddItem(logo, 0, 0, 1, 1, 0, 0, false)
    grid.AddItem(frame, 1, 0, 1, 1, 0, 0, true)


    if err := app.SetRoot(grid, true).Run(); err != nil { panic(err) }
}
