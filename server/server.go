package server

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"sync"
	"syscall"
)

type MiniServer struct {
	DefaultDir    string
	StaticHostDir *sync.Map
	ProxyHostDir  *sync.Map
}

func (m *MiniServer) check() {
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
}

func (m *MiniServer) Start(port string) {
	m.check()
	if m.DefaultDir == "" {
		m.DefaultDir = "./static"
	}
	defaultServer := http.FileServer(http.Dir(m.DefaultDir))

	// 自定义处理器
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		host := r.Host
		if m.StaticHostDir != nil {
			if path, ok := m.StaticHostDir.Load(host); ok {
				server := http.FileServer(http.Dir(path.(string)))
				http.StripPrefix("/", server).ServeHTTP(w, r)
				return
			}
		}
		if m.ProxyHostDir != nil {
			if addr, ok := m.ProxyHostDir.Load(host); ok {
				proxyURL, _ := url.Parse(addr.(string)) // 替换为目标服务器2的URL
				proxyServer := httputil.NewSingleHostReverseProxy(proxyURL)
				http.StripPrefix("/", proxyServer).ServeHTTP(w, r)
				return
			}
		}
		// 默认提供静态文件服务
		http.StripPrefix("/", defaultServer).ServeHTTP(w, r)
	})

	// 启动HTTP服务器
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
