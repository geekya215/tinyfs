package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	port     int
	serveDir string
)

func CheckDir(dir string) error {
  f, err := os.Stat(dir)

  if os.IsNotExist(err) {
    return fmt.Errorf("%s is not exist", dir)
  }

  if !f.IsDir() {
    return fmt.Errorf("%s is not a directory", dir)
  }

  return nil
}

func main() {
	flag.IntVar(&port, "port", 9876, "set tinyfs port")
	flag.StringVar(&serveDir, "dir", "./", "set tinyfs serve directory")
	flag.Parse()

  e := CheckDir(serveDir)
  if e != nil {
    log.Fatal(e)
    os.Exit(1)
  }

	http.Handle("/", http.FileServer(http.Dir(serveDir)))

	go func() {
	  e := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if e != nil {
      log.Fatal(e)
			os.Exit(1)
		}
	}()
	osCh := make(chan os.Signal, 1)
	fmt.Printf("tinyfs start at port:%d, serve directory:%s\n", port, serveDir)
	signal.Notify(osCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT) // , syscall.SIGSTOP) cannot compile on windows
	fmt.Printf("\rGot a signal [%s]\n", <-osCh)
}
