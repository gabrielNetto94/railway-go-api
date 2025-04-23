package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Application struct {
	Config Config
	Logger *slog.Logger
}

type Config struct {
	Port int
}

const version = "0.0.1"

func (app *Application) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {

	app.Logger.Info("healthcheck endpoint hit")
	// Set the content type to application/json
	data := map[string]string{
		"status":  "ok3",
		"version": version,
	}

	// No need to add acess control origin headers. On other routes, that may be necessary
	if err := app.WriteJSON(w, http.StatusOK, data, nil); err != nil {
		app.Logger.Error(err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func (app *Application) Ping(w http.ResponseWriter, r *http.Request) {

	// Set the content type to application/json
	data := map[string]string{
		"message": "pong",
	}

	// No need to add acess control origin headers. On other routes, that may be necessary
	if err := app.WriteJSON(w, http.StatusOK, data, nil); err != nil {
		app.Logger.Error(err.Error())
		http.Error(w, "server error", http.StatusInternalServerError)
	}
}

func (app *Application) Routes() *httprouter.Router {
	router := httprouter.New()

	// Define the available routes
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.HealthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/ping", app.Ping)

	return router
}

// The writeJSON() method is a generic helper for writing JSON to a response
func (app *Application) WriteJSON(w http.ResponseWriter, sCode int, data any, headers http.Header) error {
	marshalledJson, err := json.Marshal(data)

	if err != nil {
		return err
	}
	marshalledJson = append(marshalledJson, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sCode)
	w.Write(marshalledJson)

	return nil
}
