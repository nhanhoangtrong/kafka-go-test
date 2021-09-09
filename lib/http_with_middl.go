package main

import (
	"context"
	"fmt"
	"net/http"
)

func main() {
	finalHandler := http.HandlerFunc(testHandler)
	http.Handle("/test", addContextWithValue(finalHandler))

	http.ListenAndServe("localhost:8080", nil)
}

func addContextWithValue(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "Home", "Hello world!")
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func testHandler(w http.ResponseWriter, req *http.Request) {
	val := req.Context().Value("Home").(string)
	fmt.Println(val)
	w.Write([]byte("Bye"))
}
