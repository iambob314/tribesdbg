package main

import (
	"github.com/hpcloud/tail"
	"gopkg.in/tomb.v1"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	stat, err := os.Stat("console.log")
	if err != nil {
		log.Fatal(err)
	}
	off := stat.Size()
	_ = off

	cmd := exec.Command("T1Vista.exe", os.Args[1:]...)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	time.Sleep(3 * time.Second) // allow Tribes to start, or it gets cranky if we have the file open already
	f, err := tail.TailFile("console.log", tail.Config{
		Location: &tail.SeekInfo{Offset: off},
		ReOpen:   true,
		Follow:   true,
		Poll:     true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer f.Cleanup()

	log.Println("starting up...")

	go func() {
		log.Print("Waiting for Tribes...")
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
		log.Print("Tribes has closed, wrapping up...")
		_ = f.StopAtEOF()
	}()

	for line := range f.Lines {
		log.Println("LINE", line.Time.Format("15:04:05"), strings.TrimSpace(line.Text), "ERROR", line.Err)
		if err := f.Err(); err != nil && err != tomb.ErrStillAlive {
			log.Fatal(err)
		}
	}
	if err := f.Wait(); err != nil && err != tomb.ErrStillAlive {
		log.Fatal(err)
	}
}
