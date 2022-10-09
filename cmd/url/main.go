package main

import (
	"fmt"
	"time"

	app "github.com/kormiltsev/url-testtask/internal/app"
	serv "github.com/kormiltsev/url-testtask/internal/server"
	st "github.com/kormiltsev/url-testtask/internal/storage"
)

func main() {
	// upload cofig and db
	app.Initial()

	if app.GetDBtype() == "local" {
		//autosave every n seconds
		go st.AutosaverDB(app.GetCatalog(), time.Millisecond*time.Duration(app.GetSaveTimer()))
	} else {
		fmt.Println("Postgres status: ", st.Pingdb())
	}
	// start server
	serv.StartServe()
}
