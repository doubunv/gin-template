package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	_ "project-api/api"
	"project-api/config"
	ml "project-api/middleware"
	"project-api/router"
	"strconv"
	"time"
)

func runApi(r *gin.Engine, addr string) {
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Println("server start port:", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("server stopping")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalln("server stop exception:", err)
	}

	select {
	case <-ctx.Done():
		log.Println("server stop time out")
	}
}

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(ml.ApiCors())
	r.Use(ml.CheckHeader)
	r.Use(ml.ApiParams)
	r.Use(ml.AuthToken)
	r.Use(ml.AuthWhitePath)
	router.InitRouter(r)
	runApi(r, fmt.Sprintf(":%s", strconv.Itoa(*config.Server.Port)))
}
