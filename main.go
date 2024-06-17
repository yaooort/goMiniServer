package main

import (
	"fmt"
	"goNginx/ui"
	"os"
	"strconv"
	"syscall"
)

// 有些windows系统不支持opengl 解决方法 https://github.com/fyne-io/fyne/issues/4033#issuecomment-1652726057
func main() {
	pidFile := "app.pid"
	// 检查 PID 文件是否存在
	if _, err := os.Stat(pidFile); err == nil {
		// PID 文件存在，读取 PID
		data, err := os.ReadFile(pidFile)
		if err != nil {
			//fmt.Printf("读取 PID 文件失败: %v\n", err)
			os.Remove(pidFile)
			os.Exit(1)
		}
		pid, err := strconv.Atoi(string(data))
		if err != nil {
			//fmt.Printf("文件中的 PID 无效: %v\n", err)
			os.Remove(pidFile)
			os.Exit(1)
		}
		process, err := os.FindProcess(pid)
		if err == nil {
			// 进程存在，尝试发送信号以检查是否正在运行
			if err := process.Signal(syscall.Signal(0)); err == nil {
				// 进程正在运行，先停止它
				fmt.Println("应用程序已在运行，正在重启...")
				if err := process.Signal(syscall.SIGTERM); err != nil {
					fmt.Printf("启动应用程序失败: %v\n", err)
					os.Exit(1)
				}
				//fmt.Printf("已停止 PID 为 %d 的应用程序\n", pid)
			}
		}
		// 删除旧的 PID 文件
		if err := os.Remove(pidFile); err != nil {
			fmt.Printf("启动应用程序失败: %v\n", err)
		}
	}

	// 启动新实例
	pid := os.Getpid()
	//fmt.Printf("启动应用程序，PID 为 %d\n", pid)
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644); err != nil {
		fmt.Printf("启动应用程序失败: %v\n", err)
		os.Exit(1)
	}
	mainPage := ui.IndexPage{}
	mainPage.Show()
}
