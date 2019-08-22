package main

import (
	"flag"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	f := flag.String("from", "", "from path")
	t := flag.String("to", "", "to path")
	flag.Parse()

	workers := runtime.GOMAXPROCS(0)
	ch := make(chan string, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range ch {
				exec.Command("cp", task, *t).Run()
			}
		}()
	}

	paths := dirwalk(*f)
	for i := 0; i < len(paths); i++ {
		ch <- paths[i]
	}
	close(ch)
	wg.Wait()
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}
