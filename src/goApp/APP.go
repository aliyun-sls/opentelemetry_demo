package main

import (
	"net/http"
)

type APP struct {
}

func (a *APP) GetApp01(w http.ResponseWriter, r *http.Request) {
	// 处理请求的逻辑
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from /getapp01"))
}

func main() {
	app := &APP{}
	http.HandleFunc("/getapp01", app.GetApp01)
	http.ListenAndServe(":8080", nil)
}
