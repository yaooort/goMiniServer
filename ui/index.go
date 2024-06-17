package ui

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"goNginx/server"
	"goNginx/ui/theme"
	"strconv"
)

func logLifecycle(a fyne.App) {
	//a.Lifecycle().SetOnStarted(func() {
	//	fmt.Println("Lifecycle: Started")
	//})
	//a.Lifecycle().SetOnStopped(func() {
	//	fmt.Println("Lifecycle: Stopped")
	//})
	//a.Lifecycle().SetOnEnteredForeground(func() {
	//	fmt.Println("Lifecycle: Entered Foreground")
	//})
	//a.Lifecycle().SetOnExitedForeground(func() {
	//	fmt.Println("Lifecycle: Exited Foreground")
	//})
}

type IndexPage struct {
}

func (p *IndexPage) Show() {
	myApp := app.New()
	logLifecycle(myApp)
	myApp.Settings().SetTheme(&theme.MyTheme{})
	myWindow := myApp.NewWindow("老莫服务器")

	myWindow.SetContent(p.mainUI(myWindow))
	myWindow.Resize(fyne.NewSize(280, 200))
	myWindow.SetFixedSize(true)
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

// mainUI 界面
func (p *IndexPage) mainUI(myWindow fyne.Window) *fyne.Container {

	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("请输入端口号:")
	portEntry.Text = "80"

	rootEntry := widget.NewEntry()
	rootEntry.SetPlaceHolder("请输入网站根目录:")
	rootEntry.Text = "./static"

	startButton := widget.NewButton("启动服务器", nil)
	startButton.Importance = widget.HighImportance
	started := false
	ms := server.MiniServer{}
	//
	var cancelFunc context.CancelFunc
	ctx, cancel := context.WithCancel(context.Background())

	startButton.OnTapped = func() {
		port := portEntry.Text
		root := rootEntry.Text
		if !isValidPort(port) {
			showError(myWindow, "端口无效", "请输入有效的端口号 (0-65535).")
			return
		}
		ms.DefaultDir = root
		if started {
			started = false
			startButton.SetText("启动服务器")
			//startButton.TextStyle = fyne.TextStyle{}
			startButton.Refresh()
			fmt.Println("Server stopped")
			if cancelFunc != nil {
				cancelFunc()
			}
			ms.Stop(ctx)
			startButton.Importance = widget.HighImportance
			// Add server stop logic here
		} else {
			ctx, cancel = context.WithCancel(context.Background())
			started = true
			startButton.SetText("停止服务器")
			//startButton.TextStyle = fyne.TextStyle{Bold: true}
			startButton.Refresh()
			fmt.Println("Server started on port:", port, "with root directory:", root)
			cancelFunc = cancel
			go ms.Start(ctx, port)
			startButton.Importance = widget.SuccessImportance
		}
	}

	return container.NewVBox(
		widget.NewLabel("端口:"),
		portEntry,
		widget.NewLabel("网站根目录:"),
		rootEntry,
		layout.NewSpacer(), // 添加间隔
		startButton,
	)

}
func isValidPort(port string) bool {
	if port == "" {
		return false
	}
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 0 || portNum > 65535 {
		return false
	}
	return true
}

func showError(win fyne.Window, title, message string) {
	var dialogShow *widget.PopUp
	dialogShow = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabel(message),
			widget.NewButton("关闭", func() {
				dialogShow.Hide()
			}),
		),
		win.Canvas(),
	)
	dialogShow.Show()
}
