package ui

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	theme2 "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"goNginx/resource"
	"goNginx/server"
	"goNginx/ui/theme"
	"net"
	"strconv"
	"strings"
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
	if desk, ok := myApp.(desktop.App); ok {
		m := fyne.NewMenu("服务器",
			fyne.NewMenuItem("控制台", func() {
				myWindow.Show()
			}))
		desk.SetSystemTrayIcon(fyne.NewStaticResource("logo", resource.Logo))
		desk.SetSystemTrayMenu(m)
	}
	myWindow.SetContent(p.mainUI(myWindow, myApp))
	myWindow.Resize(fyne.NewSize(280, 230))
	myWindow.SetFixedSize(true)
	myWindow.CenterOnScreen()
	myWindow.SetCloseIntercept(func() {
		myWindow.Hide()
	})
	myWindow.ShowAndRun()
}

// mainUI 界面
func (p *IndexPage) mainUI(myWindow fyne.Window, myApp fyne.App) *fyne.Container {

	localIp := widget.NewLabel("本机IP:" + showLocalIp())

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
	message := make(chan string)
	go func() {
		for {
			select {
			case msg := <-message:
				showError(myWindow, "启动失败", msg)
				started = false
				startButton.SetText("启动服务器")
				//startButton.TextStyle = fyne.TextStyle{}
				fmt.Println("Server stopped")
				if cancelFunc != nil {
					cancelFunc()
				}
				ms.Stop(ctx)
				startButton.Importance = widget.HighImportance
				startButton.Refresh()
			}
		}
	}()
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
			fmt.Println("Server stopped")
			if cancelFunc != nil {
				cancelFunc()
			}
			ms.Stop(ctx)
			startButton.Importance = widget.HighImportance
			startButton.Refresh()
			// Add server stop logic here
		} else {
			ctx, cancel = context.WithCancel(context.Background())
			started = true
			startButton.SetText("停止服务器")
			//startButton.TextStyle = fyne.TextStyle{Bold: true}
			fmt.Println("Server started on port:", port, "with root directory:", root)
			cancelFunc = cancel
			go ms.Start(ctx, port, message)
			startButton.Importance = widget.SuccessImportance
			startButton.Refresh()
		}
	}
	folderButton := widget.NewButtonWithIcon("选择目录", theme2.FolderOpenIcon(), func() {
		w2 := myApp.NewWindow("Select")
		w2.Resize(fyne.NewSize(800, 600))
		w2.Show()
		folderDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			defer w2.Close()
			if err != nil {
				dialog.ShowError(err, w2)
				return
			}
			if uri == nil {
				return
			}
			fmt.Println("Selected folder:", uri.Path())
			rootEntry.Text = uri.Path()
			rootEntry.Refresh()
		}, w2)
		folderDialog.SetFilter(storage.NewExtensionFileFilter([]string{""})) // Allow all folders
		folderDialog.Resize(fyne.NewSize(800, 600))
		folderDialog.Show()
	})
	return container.NewVBox(
		localIp,
		widget.NewLabel("端口:"),
		portEntry,
		widget.NewLabel("网站根目录:"),
		folderButton,
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
	if strings.Contains(message, "bind: address already in use") {
		message = strings.Replace(message, "bind: address already in use", "端口已被占用", 1)
	}
	if strings.Contains(message, "listen tcp") {
		message = strings.Replace(message, "listen tcp", "监听TCP", 1)
	}
	var dialogShow *widget.PopUp
	label := widget.NewRichTextWithText(message)
	// Optional: Set text style and wrapping
	label.Wrapping = fyne.TextWrapWord
	dialogShow = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabelWithStyle(title, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			label,
			widget.NewButton("关闭", func() {
				dialogShow.Hide()
			}),
		),
		win.Canvas(),
	)
	dialogShow.Resize(fyne.NewSize(250, 130))
	dialogShow.Show()
}

func showLocalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return "127.0.0.1"
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				fmt.Println("IP address:", ipNet.IP.String())
				return ipNet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
