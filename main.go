package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
	"v-blog/config"
	"v-blog/databases"
	"v-blog/helpers"
	"v-blog/routers"
)

var (
	s          string
	d          bool
	pid        int
	conf       *config.Config
	ln         net.Listener
	httpServer *http.Server
	shutdown   chan bool
	isDaemon   string
)

func init() {
	var err error
	pid = os.Getpid()
	isDaemon = os.Getenv("daemon")

	// 初始化配置
	config.InitConfig(&config.Param{})
	conf, err = config.Get()
	if err != nil {
		log.Fatalf("Config get error: %s", err)
	}

	// 初始化日志
	initLog()

	// cli 参数解析
	flag.StringVar(&s, "s", "", "send signal to a master process: quit, stop, reload, reopen")
	flag.BoolVar(&d, "d", false, "should run in daemon mod")
}

// 初始化日志
// 前台运行默认输出到终端
func initLog() {
	if isDaemon != "true" {
		return
	}
	lg := conf.LogFile
	if lg == "" {
		fmt.Printf("Init Log fail. Log file is not set.")
		os.Exit(1)
	}
	fn := path.Base(lg)
	splitFileName := strings.Split(fn, ".")
	if len(splitFileName) != 2 {
		fmt.Printf("Log file is error. [%s]", fn)
		os.Exit(1)
	}

	lf := fmt.Sprintf("%s/%s.%s.%s", path.Dir(lg), splitFileName[0], "%Y%m%d", splitFileName[1])
	rl, err := rotatelogs.New(lf, rotatelogs.WithRotationTime(time.Hour*24), rotatelogs.WithMaxAge(time.Hour*24*30))
	if err != nil {
		fmt.Printf("New rotatelogs err: %s", err)
		os.Exit(1)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(rl)
	log.SetLevel(log.DebugLevel)
}

func main() {
	flag.Parse()
	// 信号不为空
	if s != "" {
		handleCmdSignals()
		os.Exit(0)
	}

	if d && os.Getenv("daemon") != "true" {
		daemon()
	}
	log.Println("start")
	start()
}

func handleCmdSignals() {
	pid, err := readPid()
	if err != nil {
		log.Fatalf("Pid read fail. err: %s", err)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		log.Fatalf("Find application proc [%d] err: %s", pid, err)
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
		log.Fatalf("Signal send fail. err: %s", err)
	}

	log.Printf("Sent signal: %v", sig)
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

	log.Printf("Run http server. pid: %d", os.Getpid())

	server := &http.Server{
		Handler:      routers.Router,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	httpServer = server

	if os.Getenv("graceful") == "true" {
		var err error
		fd := os.NewFile(3, "")
		ln, err = net.FileListener(fd)
		if err != nil {
			log.Fatalf("Listen file err: %s", err)
		}
	} else {
		var err error
		ln, err = net.Listen("tcp", ":8888")
		if err != nil {
			log.Fatalf("tcp listen err: %s", err)
		}
	}

	go func() {
		err := server.Serve(ln)
		if err != nil {
			log.Printf("http server shutdown: %s", err)
		}
	}()

	log.Println("Server listen in port: 8888")
	err := writePid()
	if err != nil {
		log.Fatalf("Pid write fail. err: %s", err)
	}

	<-shutdown
	log.Println("Proc exited!")
}

func daemon() {
	log.Println("Start in daemon.")
	_ = os.Setenv("daemon", "true")

	stdFile, err := os.OpenFile(conf.DebugLogFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("OpenFile v-blog.debug err: %s", err)
	}

	// fork 子进程
	procAttr := &syscall.ProcAttr{Env: os.Environ(), Files: []uintptr{stdFile.Fd(), stdFile.Fd(), stdFile.Fd()}}
	pid, err := syscall.ForkExec(os.Args[0], os.Args, procAttr)
	if err != nil {
		log.Fatalf("Start daemon err: %s", err)
	}

	log.Printf("Daemon proc started. pid: %d", pid)
	// 主进程退出
	os.Exit(0)
}

// 安装信号器
// 前台运行忽忽略信号器
func installSignals() {
	if isDaemon != "true" {
		return
	}
	s := make(chan os.Signal)
	// SIGINT - 中断； SIGTERM - 结束；SIGQUIT - 退出
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	for sig := range s {
		log.Printf("Receive signal: %v ", sig)
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

	log.Println("Signals intalled.")
}

// 立即退出进程
func quit() {
	if err := clearPid(); err != nil {
		log.Fatalf("Pid file clear fail. err: %s. pid: %d", err, os.Getpid())
	}
	log.Println("Proc quit!")
	os.Exit(0)
}

// 停止程序
func stop() {
	err := stopServe()
	if err != nil {
		log.Fatalf("Proc stop fail. err: %s", err)
	}
	if err := clearPid(); err != nil {
		log.Fatalf("Pid file clear fail. err: %s. pid: %d", err, os.Getpid())
	}
	shutdown <- true
	log.Info("Proc stop. pid: %d", os.Getpid())
}

// 重启进程
func reopen() {
	if err := os.Setenv("graceful", "false"); err != nil {
		log.Fatalf("Set env fail. err: %s", err)
	}

	if err := stopServe(); err != nil {
		log.Fatalf("Stop http server error: %s", err)
	}

	cmd := exec.Cmd{
		Path: os.Args[0],
		Args: os.Args,
		Env:  os.Environ(),
		// Stdin: os.Stdin,
		// Stdout: os.Stdout,
		// Stderr: os.Stderr,
	}

	if err := clearPid(); err != nil {
		log.Fatalf("Pid file clear fail. err: %s", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Reopen proc fail. err: %s", err)
	}
	os.Exit(0)
}

// 重载进程 - 热重启
func reload() {
	if err := os.Setenv("graceful", "true"); err != nil {
		log.Fatalf("Set env err: %s", err)
	}

	tl, ok := ln.(*net.TCPListener)
	if !ok {
		log.Fatalln("listener is not tcp listener")
	}

	tld, err := tl.File()
	if err != nil {
		log.Fatalf("Get fd err: %s", err)
	}

	if err := clearPid(); err != nil {
		log.Fatalf("Pid file clear fail. err: %s", err)
	}

	cmd := exec.Cmd{
		Path: os.Args[0],
		Args: os.Args,
		Env:  os.Environ(),
		//Stdin: os.Stdin,
		//Stdout: os.Stdout,
		//Stderr: os.Stderr,
		ExtraFiles: []*os.File{tld},
	}

	if err := cmd.Start(); err != nil {
		log.Fatalf("Reload err: %s", err)
	}

	_ = stopServe()
	shutdown <- true
}

// 写 pid 文件
func writePid() error {
	if isDaemon != "true" {
		return nil
	}

	pidPath := conf.PidFile
	exist, err := helpers.PathExists(pidPath)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("pid file already exist")
	}
	pidFile, err := os.OpenFile(pidPath, os.O_RDWR|os.O_CREATE, 0766)
	defer func() {
		_ = pidFile.Close()
	}()

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

func stopServe() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	return httpServer.Shutdown(ctx)
}
