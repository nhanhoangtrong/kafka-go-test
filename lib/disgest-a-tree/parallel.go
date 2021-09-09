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
	defer measureTime(time.Now(), "Parallel")
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

	res, errc := sumFiles(done, root)
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

func sumFiles(done <-chan struct{}, root string) (<-chan Result, <-chan error) {
	// For each regular file, start a goroutine that sums the file and sends
	// the result on c, send the result of the walk on errc
	c := make(chan Result)
	errc := make(chan error)

	go func() {
		var wg sync.WaitGroup
		err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			wg.Add(1)
			go func() {
				data, err := ioutil.ReadFile(path)
				select {
				case c <- Result{path, md5.Sum(data), err}:
				case <-done:
				}
				wg.Done()
			}()
			select {
			case <-done:
				return errors.New("walk cancelled!")
			default:
				return nil
			}
		})

		go func() {
			wg.Wait()
			close(c)
		}()
		errc <- err
	}()

	return c, errc
}
