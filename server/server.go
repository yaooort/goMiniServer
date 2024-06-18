package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type MiniServer struct {
	DefaultDir    string
	StaticHostDir *sync.Map
	ProxyHostDir  *sync.Map
	Server        *http.Server
}

func (m *MiniServer) Start(ctx context.Context, port string, message chan string) {
	if m.DefaultDir == "" {
		m.DefaultDir = "./static"
	}
	defaultServer := http.FileServer(http.Dir(m.DefaultDir))

	// 创建一个新的 ServeMux
	mux := http.NewServeMux()

	// 自定义处理器
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
				proxyURL, _ := url.Parse(addr.(string)) // 替换为目标服务器的URL
				proxyServer := httputil.NewSingleHostReverseProxy(proxyURL)
				http.StripPrefix("/", proxyServer).ServeHTTP(w, r)
				return
			}
		}
		// 默认提供静态文件服务
		http.StripPrefix("/", defaultServer).ServeHTTP(w, r)
	})

	m.Server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		if err := m.Server.Shutdown(ctx); err != nil {
			fmt.Println("Server shutdown failed:", err)
		} else {
			fmt.Println("Server gracefully stopped")
		}
	}()

	// 启动HTTP服务器
	if err := m.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		message <- err.Error()
	}
}

func (m *MiniServer) Stop(ctx context.Context) {
	if m.Server != nil {
		if err := m.Server.Shutdown(ctx); err != nil {
			fmt.Println("Server shutdown failed:", err)
		} else {
			fmt.Println("Server gracefully stopped")
		}
	}
}
