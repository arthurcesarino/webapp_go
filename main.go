package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/arthurcesarino/webapp_go/handlers"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {

	l := log.New(os.Stdout, "product-api", log.LstdFlags)
	hh := handlers.NewProducts(l)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/products", hh.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", hh.ListSingle)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/products", hh.Update)
	putRouter.Use(hh.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/products", hh.Create)
	postRouter.Use(hh.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", hh.Delete)
	//sm.Handle("/products", hh)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	server := &http.Server{
		Addr:         ":9090",
		Handler:      sm,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Println("Service starting at port:", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	server.Shutdown(tc)
}
