package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	cache *Cache
}

func StartServer(cache *Cache) {
	h := Handler{
		cache: cache,
	}
	http.HandleFunc("/", h.getInfo)

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
	} else {
		fmt.Println("Server start")
	}
}

func (h *Handler) getInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form.Get("id")
	if id != "" {
		orderF, ok := h.cache.Get(id)
		if !ok {
			fmt.Println(ok)
			w.WriteHeader(400)
			return
		}
		order, err := json.MarshalIndent(orderF, "", "    ")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(order)
	} else {
		http.ServeFile(w, r, "templates/index.html")
	}
}
