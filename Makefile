
# 使用goReleaser 编译 goreleaser release --snapshot --rm-dist

# 输出目录
OUTPUT=./build

# 编译android 可执行文件

ANDROID_OUT=$(OUTPUT)/android
ANDROID_SDK=$(HOME)/Library/Android/sdk
NDK_BIN=$(ANDROIDSDKROOT)/ndk/25.1.8937393/toolchains/llvm/prebuilt/darwin-x86_64/bin


android:
	CGO_ENABLED=1 \
	ANDROID_NDK_HOME=$(ANDROIDSDKROOT)/ndk/25.1.8937393 \
	fyne package -os android -appID com.oortk.server -icon Icon.png

# 编译windows 可执行文件
WINDOWS_OUT=$(OUTPUT)/windows

windows-x86:
	CGO_ENABLED=1 \
    GOOS=windows \
    GOARCH=386 \
    CC=i686-w64-mingw32-gcc \
    CXX=i686-w64-mingw32-g++ \
    fyne package -os windows -icon Icon.png --name server-86.exe

# upx 压缩 upx -9 -o ll.exe server.exe
# 去除黑窗口 -ldflags "-s -w -H=windowsgui"
windows-x86_64:
	CGO_ENABLED=1 \
    GOOS=windows \
    GOARCH=amd64 \
    CC=x86_64-w64-mingw32-gcc \
    CXX=x86_64-w64-mingw32-g++ \
    fyne package -os windows -icon Icon.png --name server-64.exe


windows: windows-x86 windows-x86_64

# 编译linux可执行文件
LINUX_OUT=$(OUTPUT)/linux


# x86 64-bit
linux-amd64:
	CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64 \
    CC=x86_64-linux-musl-gcc \
    CXX=x86_64-linux-musl-g++ \
    fyne package -os linux -icon Icon.png

# x86 32-bit
linux-i486:
	CGO_ENABLED=1 \
  	GOOS=linux \
  	GOARCH=386 \
  	CC=i486-linux-musl-gcc \
  	CXX=i486-linux-musl-g++ \
  	fyne package -os linux -icon Icon.png

# x86 arm64
linux-arm64:
	CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=arm64 \
    CC=aarch64-linux-musl-gcc \
    CXX=aarch64-linux-musl-g++ \
    fyne package -os linux -icon Icon.png

linux: linux-amd64 linux-arm64



# 编译macos
MACOS_OUT=$(OUTPUT)/macos

macos-amd64:
	CGO_ENABLED=1 \
    GOOS=darwin \
    GOARCH=amd64 \
  	fyne package -os darwin -icon Icon.png

macos-arm64:
	CGO_ENABLED=1 \
    GOOS=darwin \
    GOARCH=arm64 \
    fyne package -os darwin -icon Icon.png

mac:macos-amd64 macos-arm64


# 编译web
WEB_OUT=$(OUTPUT)/web

web:
	CGO_ENABLED=1 \
	GOPHERJS_GOROOT=/opt/homebrew/Cellar/go@1.18/1.18.10/libexec \
	fyne package -os web


# 编译所有
all: android windows linux mac web