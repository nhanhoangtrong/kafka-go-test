package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

const maxWorkers = 5

type Result struct {
	path string
	sum  [md5.Size]byte
	err  error
}

func measureTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func main() {
	defer measureTime(time.Now(), "Bounded Parallel")
	m, err := MD5All(os.Args[1])
	if err != nil {
		panic(err)
	}
	var paths []string
	for path := range m {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x %s\n", m[path], path)
	}
}

func MD5All(root string) (map[string][md5.Size]byte, error) {
	done := make(chan struct{})
	defer close(done)

	res := make(chan Result)
	paths, errc := walkFiles(done, root)
	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			digester(done, paths, res)
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	m := make(map[string][md5.Size]byte)
	for r := range res {
		if r.err != nil {
			return nil, r.err
		}
		m[r.path] = r.sum
	}
	if err := <-errc; err != nil {
		return nil, err
	}
	return m, nil
}

func walkFiles(done <-chan struct{}, root string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(paths)
		errc <- filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			select {
			case paths <- path:
			case <-done:
				return errors.New("walk cancelled")
			}
			return nil
		})
	}()
	return paths, errc
}

func digester(done <-chan struct{}, paths <-chan string, res chan<- Result) {
	for path := range paths {
		data, err := ioutil.ReadFile(path)
		select {
		case res <- Result{path, md5.Sum(data), err}:
		case <-done:
			return
		}
	}
}
