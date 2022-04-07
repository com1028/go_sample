package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"
	// "github.com/com1028/go_sample/database"
	// "github.com/com1028/go_sample/jobs"
	// "github.com/com1028/go_sample/models"
	// "github.com/com1028/go_sample/router"
)

func main() {
	database.ConnectDB(false)
	defer database.CloseDB()

	models.Migration(database.GetDB())

	// get entries in 5 minutes cycle
	go func() {
		jobs.GetContentfulEntries("get entries")
	}()

	e := router.NewRouter()

	go func() {
		if err := e.Start(":3030"); err != nil {
			fmt.Printf("[Error]: %s", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
