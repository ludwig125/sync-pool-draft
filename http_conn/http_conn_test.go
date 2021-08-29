package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func getName(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println("id:", id)
	if _, err := w.Write([]byte("Hello, Gophers!")); err != nil {
		log.Printf("failed to Write: %v", err)
	}
}

func startServer() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	mux := http.NewServeMux()
	mux.HandleFunc("/", getName)

	wg.Done()
	go func() {
		log.Println("start server")
		if err := http.ListenAndServe(":3000", mux); err != http.ErrServerClosed {
			log.Printf("failed to ListenAndServe: %v", err)
		}
		log.Println("server shutdown")
	}()
	return &wg
}

func requestClient() string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:3000", nil)
	if err != nil {
		log.Printf("NewRequest failed: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to Do request: %v", err)
	}
	defer resp.Body.Close()

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(res)
}

func init() {
	wg := startServer()
	wg.Wait()
	fmt.Println("done wg.Wait")
}

func TestRequest(t *testing.T) {
	got := requestClient()
	t.Log(got)

	time.Sleep(20 * time.Second)
}

func BenchmarkRequest(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		requestClient()
	}
}

func BenchmarkRequest2(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		requestClient()
	}
}
