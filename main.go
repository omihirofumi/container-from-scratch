package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bad command")
	}
}

func run() {
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

	// /proc/self/exe child /bin/bash
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	
	cmd.Run()
}

// container
func child() {
	fmt.Printf("Running %v as PID %d\n", os.Args[2:], os.Getpid())

	syscall.Sethostname([]byte("container"))
	// change root
	must(syscall.Chroot("/my-root-fs"))
	must(syscall.Chdir("/"))
	// mount proc
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	

	cmd.Run()

	syscall.Unmount("proc", 0)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
