package graceful

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	GRACEFUL_KEY = "GRACEFUL_RESTART"
)

var (
	listener      *net.TCPListener
	stopServer    = make(chan bool)
	signalChannel = make(chan os.Signal)
)

func StartGin(engine *gin.Engine, addr ...string) (err error) {
	defer func() {
		if err != nil && gin.IsDebugging() {
			fmt.Println("[GIN-debug] [ERROR] %v\n"+err.Error())
		}
	}()
	address := resolveAddress(addr)
	debugPrint("Listening and serving HTTP on %s\n", address)

	var (
		_addr *net.TCPAddr
	)
	_addr, err = net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return
	}
	if isGraceful() {
		fmt.Println("Graceful Restart")
		listener, err = getListenerFd()
		if err != nil {
			fmt.Println("Graceful restart failed, " + err.Error())
		}
		if listener != nil && !isSameAddr(listener.Addr(), _addr) {
			listener = nil
		}
		os.Setenv(GRACEFUL_KEY, "")
	}
	if listener == nil {
		listener, err = net.ListenTCP("tcp", _addr)
	}
	if err != nil {
		return
	}
	server := &http.Server{Handler: engine}
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR2)
	go signalHandler()
	go func() {
		<-stopServer
		fmt.Println("Shutting down")
		listener.SetDeadline(time.Now())
		server.SetKeepAlivesEnabled(false)
		time.Sleep(10 * time.Second)
		if err = server.Shutdown(context.Background()); err != nil {
			fmt.Println("Shutting down failed.", err.Error())
		}
	}()
	err = server.Serve(listener)
	fmt.Println("Stopped,", err.Error())
	//err = http.ListenAndServe(address, engine)
	return
}

func signalHandler() {
	for {
		sig := <-signalChannel
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			stopServer <- true
		case syscall.SIGUSR2:
			if err := newProcess(); err != nil {
				fmt.Println("Start newProcess failed, " + err.Error())
			} else {
				fmt.Println("Stop old process")
				fmt.Println("Start newProcess success")
				syscall.Kill(os.Getpid(), syscall.SIGQUIT)
			}
		}
	}
}

func newProcess() (err error) {
	envs := []string{}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, GRACEFUL_KEY) {
			continue
		}
		envs = append(envs, env)
	}
	arg0, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}
	wd, _ := os.Getwd()
	envs = append(envs, fmt.Sprintf("%s=%d", GRACEFUL_KEY, 1))
	lf, err := listener.File()
	if err != nil {
		return err
	}
	_, err = os.StartProcess(arg0, os.Args, &os.ProcAttr{
		Dir:   wd,
		Env:   envs,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr, lf},
	})
	return
}

func isSameAddr(a1, a2 net.Addr) bool {
	if a1.Network() != a2.Network() {
		return false
	}
	a1s := a1.String()
	a2s := a2.String()
	if a1s == a2s {
		return true
	}

	// This allows for ipv6 vs ipv4 local addresses to compare as equal. This
	// scenario is common when listening on localhost.
	const ipv6prefix = "[::]"
	a1s = strings.TrimPrefix(a1s, ipv6prefix)
	a2s = strings.TrimPrefix(a2s, ipv6prefix)
	const ipv4prefix = "0.0.0.0"
	a1s = strings.TrimPrefix(a1s, ipv4prefix)
	a2s = strings.TrimPrefix(a2s, ipv4prefix)
	return a1s == a2s
}

func isGraceful() bool {
	return os.Getenv(GRACEFUL_KEY) != ""
}

func getListenerFd() (listener *net.TCPListener, err error) {
	file := os.NewFile(uintptr(3), "")
	_listener, err := net.FileListener(file)
	if err != nil {
		file.Close()
		fmt.Println("Graceful restart error, File Descripter error")
		return
	}
	if err = file.Close(); err != nil {
		fmt.Println("Graceful restart error, error close fd")
	}
	listener = _listener.(*net.TCPListener)
	return
}

func debugPrint(format string, values ...interface{}) {
	if gin.IsDebugging() {
		fmt.Println("[GIN-debug] "+format)
	}
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			debugPrint("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		debugPrint("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too much parameters")
	}
}
