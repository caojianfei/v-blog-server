package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {

	fmt.Printf("Proc pid: %d, os.Args: %v\n", os.Getpid(), os.Args)

	path := os.Args[0]
	fmt.Println("path", path)

	pid, err := syscall.ForkExec(path, os.Args, &syscall.ProcAttr{
		Env: os.Environ(),
		Files: []uintptr{
			os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(),
		}})
	if err != nil {
		fmt.Println("fork error: ", err)
		os.Exit(0)
	}
	fmt.Println("fork proc pid: ", pid)
}
