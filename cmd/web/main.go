package main

import (
	"log"
	"net/http"
	"time"

	"github.com/cetingzhou/bookings/pkg/config"

	"github.com/alexedwards/scs/v2"
	"github.com/cetingzhou/bookings/pkg/handlers"
	"github.com/cetingzhou/bookings/pkg/renders"
)

const portNumber = ":8080"

var App config.AppConfig
var Session *scs.SessionManager

func main() {

	App.InProduction = false

	Session = scs.New()
	Session.Lifetime = 24 * time.Hour
	Session.Cookie.Persist = true
	Session.Cookie.SameSite = http.SameSiteLaxMode
	Session.Cookie.Secure = App.InProduction

	App.Session = Session

	tc, err := renders.CreateTemplateCache()
	if err != nil {
		log.Fatal("Failed to create template cache")
	}
	App.TemplateCache = tc
	App.UseCache = false

	repo := handlers.NewRepo(&App)
	handlers.NewHandlers(repo)

	renders.NewTemplates(&App)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&App),
	}

	log.Printf("Starting application on port%s\n", portNumber)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
