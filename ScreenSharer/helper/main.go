package main

import (
	"context"
	"fmt"
	helper "helper/pkgs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	if err := helper.InitDB(); err != nil {
		log.Fatalf("数据库错误: %v", err)
		return
	}

	server := helper.Server()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8080),
		Handler: server,
	}

	go func() {
		log.Printf("url: http://%s:%s/home", helper.Host, helper.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown, err: %v", err)
	}
	log.Println("Server exiting")

}
