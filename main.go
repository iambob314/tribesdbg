package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
)

var syntaxErrPat = regexp.MustCompile(`^(.*) Line: ([0-9]+) - Syntax error\.$`)
var missingFuncPat = regexp.MustCompile(`^(.*): Unknown command\.$`)

func handleLine(l string) {
	switch {
	case syntaxErrPat.MatchString(l), missingFuncPat.MatchString(l):
		log.Println(l)
	}
}

func main() {
	if len(os.Args) == 1 {
		log.Fatalf("usage: tribesdbg.exe tribes-exe [args..]")
	}

	stat, err := os.Stat("console.log")
	if err != nil {
		log.Fatal(err)
	}
	off := stat.Size()
	_ = off

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	log.Print("Waiting for Tribes...")
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	log.Print("Tribes has closed, wrapping up...")

	f, err := os.Open("console.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := f.Seek(off, io.SeekStart); err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(f)
	for {
		switch l, _, err := r.ReadLine(); err {
		case nil:
			handleLine(string(l))
		case io.EOF:
			return
		default:
			log.Fatal(err)
		}
	}
}
