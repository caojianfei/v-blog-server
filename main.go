package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"v-blog/config"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/routers"
)

var (
	s string
	d bool
	pid int
	conf *config.Config
	ln net.Listener
	httpServer *http.Server
	shutdown chan bool
)

func init() {
	var err error
	pid = os.Getpid()
	// 初始化配置
	config.InitConfig(&config.Param{})
	conf, err = config.Get()
	if err != nil {
		log.Fatalf("Config get error: %s\n", err)
	}
	// 参数设置
	flag.StringVar(&s, "s", "", "send signal to a master process: quit, stop, reload, reopen")
	flag.BoolVar(&d, "d", false, "should run in daemon mod")
}

func main() {
	fmt.Println("=============================")
	fmt.Println("os.Args", os.Args)
	fmt.Println("os.Environ", os.Environ())
	flag.Parse()
	// 信号不为空
	if s != "" {
		handleCmdSignals()
		os.Exit(0)
	}
	if d && os.Getenv("daemon") != "true" {
		daemon()
	}

	start()
}

func start() {

	shutdown = make(chan bool)

	// 安装信号监听期
	go installSignals()
	// 初始化数据库
	databases.InitDatabase()
	// 初始化验证器
	helpers.InitValidator()
	// 初始化路由
	routers.InitRouter()

	log.Printf("Run http server. pid: %d\n", os.Getpid())

	server := &http.Server{
		Handler:     routers.Router,
		ReadTimeout: time.Second,
	}

	httpServer = server

	if os.Getenv("graceful") == "true" {
		var err error
		fd := os.NewFile(3, "")
		ln, err = net.FileListener(fd)
		if err != nil {
			log.Fatalf("Listen file err: %s\n", err)
		}
	} else {
		var err error
		ln, err = net.Listen("tcp", ":8888")
		if err != nil {
			log.Fatalf("tcp listen err: %s\n", err)
		}
	}

	go func() {
		err := server.Serve(ln)
		if err != nil {
			log.Printf("http server shutdown: %s\n", err)
		}
	}()

	log.Println("Server listen in port: 8888")
	err := writePid()
	if err != nil {
		log.Fatalf("Pid write fail. err: %s\n", err)
	}

	<-shutdown
	log.Println("Proc exited!")
}

func daemon() {
	log.Println("Start in daemon.")
	_ = os.Setenv("daemon", "true")

	// fork 子进程
	procAttr := &syscall.ProcAttr{Env: os.Environ(), Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()}}
	// procAttr := &syscall.ProcAttr{Env: os.Environ(), Files: []uintptr{}}
	pid, err := syscall.ForkExec(os.Args[0], os.Args, procAttr)
	if err != nil {
		log.Fatalf("Start daemon err: %s\n", err)
	}

	log.Printf("Daemon proc started. pid: %d\n", pid)
	// 主进程退出
	os.Exit(0)
}

func installSignals() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Recover from installSignals func: ", err)
		}
	}()

	s := make(chan os.Signal)
	// SIGINT - 中断； SIGTERM - 结束；SIGQUIT - 退出
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for sig := range s {
		log.Printf("Receive signal: %v \n", sig)
		switch sig {
		case syscall.SIGINT:
			quit()
		case syscall.SIGTERM, syscall.SIGQUIT:
			stop()
		case syscall.SIGUSR1: // 重启
			reopen()
		case syscall.SIGUSR2: // 热重启
			reload()
		default:
			log.Println("Received signal ignored: ", sig)
		}
	}
}

func handleCmdSignals() {
	pid, err := readPid()
	if err != nil {
		log.Fatalf("Pid read fail. err: %s\n", err)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Fatalf("Find application proc [%d] err: %s\n", pid, err)
	}

	var sig os.Signal

	switch s {
	case "quit":
		sig = syscall.SIGINT
	case "stop":
		sig = syscall.SIGTERM
	case "reopen":
		sig = syscall.SIGUSR1
	case "reload":
		sig = syscall.SIGUSR2
	}

	if sig == nil {
		log.Fatalf("Unknow signal: %s", s)
	}

	err = proc.Signal(sig)
	if err != nil {
		log.Fatalf("Signal send fail. err: %s\n", err)
	}

	log.Printf("Sended signal: %v\n", sig)
}

// 立即退出进程
func quit() {
	if err := clearPid(); err != nil {
		log.Fatalf("Pid file clear fail. err: %s\n", err)
	}
	log.Println("Proc quit!")
	os.Exit(0)
}

// 停止程序
func stop() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second * 5)
	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Proc stop fail. err: %s\n", err)
	}
	_ = clearPid()
	shutdown <- true
}

// 重启进程
func reopen() {
	if err := os.Setenv("graceful", "false"); err != nil {
		log.Fatalf("Set env fail. err: %s", err)
	}
	cmd := exec.Cmd{
		Path: os.Args[0],
		Args: os.Args,
		Env: os.Environ(),
		Stdin: os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := clearPid(); err != nil {
		log.Fatalf("Pid file clear fail. err: %s\n", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Reopen proc fail. err: %s\n", err)
	}
	os.Exit(0)
}

// 重载进程 - 热重启
func reload() {
	if err := os.Setenv("graceful", "true"); err != nil {
		log.Fatalf("Set env err: %s\n", err)
	}

	tl, ok := ln.(*net.TCPListener)
	if !ok {
		log.Fatalln("listener is not tcp listener")
	}

	tld, err := tl.File()
	if err != nil {
		log.Fatalf("Get fd err: %s\n", err)
	}

	cmd := exec.Cmd{
		Path: os.Args[0],
		Args: os.Args,
		Env: os.Environ(),
		Stdin: os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		ExtraFiles: []*os.File{tld},
	}

	if err := clearPid(); err != nil {
		log.Fatalf("Pid file clear fail. err: %s\n", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Reload err: %s\n", err)
	}
}

// 写 pid 文件
func writePid() error {
	var pidFile *os.File
	defer func() {
		if pidFile != nil {
			_ = pidFile.Close()
		}
	}()

	pidPath := conf.PidFile
	exist, err := helpers.PathExists(pidPath)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("pid file already exist")
	}
	pidFile, err = os.OpenFile(pidPath, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return errors.New(fmt.Sprintf("pid file create err: %s", err))
	}

	log.Println("Writing pid: ", pid)
	_, err = pidFile.Write([]byte(strconv.Itoa(pid)))
	if err != nil {
		return errors.New(fmt.Sprintf("pid file write err: %s", err))
	}

	return nil
}

// 删除pid 文件
func clearPid() error {
	pidPath := conf.PidFile
	exist, err := helpers.PathExists(pidPath)
	if err != nil {
		return err
	}

	if exist == false {
		return errors.New("pid file is not exist")
	}

	return os.Remove(pidPath)
}

// 读取 pid
func readPid() (int, error) {
	pidPath := conf.PidFile
	exist, err := helpers.PathExists(pidPath)
	if err != nil {
		return 0, err
	}

	if exist == false {
		 return 0, errors.New("pid file is not exist")
	}

	pidByte, err := ioutil.ReadFile(pidPath)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(pidByte))
}






