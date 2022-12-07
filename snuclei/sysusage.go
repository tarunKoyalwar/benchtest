package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func main() {

	arr := os.Args[1:]
	givencmd := strings.Join(arr, " ")
	fmt.Printf("[*] %v\n", givencmd)

	cmd := exec.Command("/bin/sh", "-c", givencmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	if cmd.ProcessState != nil {
		usage := cmd.ProcessState.SysUsage().(*syscall.Rusage)
		fmt.Printf("\nMax RSS: %v MB ,Voluntary Context Switch (nvcsw): %v\n", usage.Maxrss/(1024*1024), usage.Nvcsw)
		fmt.Printf("Involuntary Context Switch(nivcsw): %v,nswap: %v\n", usage.Nivcsw, usage.Nswap)
		fmt.Printf("Sys Time: %vs , User Time: %vs\n", usage.Stime.Sec, usage.Utime.Sec)
	}

}
