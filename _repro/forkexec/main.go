// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {
	if os.Args[1] != "1" {
		return
	}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// does not matter what we exec, just exec itself
			cmd := exec.Command("./forkexec.exe", "0")
			cmd.Run()
		}()
	}

	var wg2 sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			// does not matter what we exec, just exec itself
			client := http.DefaultClient
			client.Timeout = 10 * time.Second
			resp, err := client.Get("https://google.com")
			if err != nil {
				// log.Warnf("HTTP error: %s", err)
			}

			resp.Body.Close()
		}()
	}

	wg.Wait()
	wg2.Wait()
}
