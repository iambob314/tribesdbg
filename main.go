package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	f, err := os.Open("console.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if n, err := io.Copy(io.Discard, f); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("caught up console.log by %d bytes", n)
	}

	cmd := exec.Command("T1Vista.exe", os.Args...)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	stopCh := make(chan struct{})
	lineCh := make(chan string, 1)
	doneCh := make(chan struct{})
	go func() {
		defer close(lineCh)
		buf := bytes.Buffer{}
		skipNextNL := false

		for {
			if n, err := io.Copy(&buf, f); err != nil {
				log.Fatal(err)
			} else if n != 0 {
				b := buf.Bytes()
				skip := 0

				for len(b) > 0 {
					if skipNextNL && b[0] == '\n' {
						skip++
						b = b[1:]
					}

					idx := bytes.IndexAny(b, "\r\n")
					if idx == -1 {
						break
					}
					lineCh <- string(b[:idx])

					if b[idx] == '\r' {
						skipNextNL = true
					}
					skip += idx + 1
					b = b[idx+1:]
				}

				buf.Next(skip)
				continue
			}

			log.Println("wow")

			select {
			case <-stopCh:
				return
			case <-time.After(time.Second):
				continue
			}
		}
	}()

	go func() {
		defer close(doneCh)
		for line := range lineCh {
			log.Println(line)
		}
	}()

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	close(stopCh)
	<-doneCh
}
