package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/pmlogist/1am-audio-srv/internal/data"
	"github.com/pmlogist/1am-audio-srv/internal/validator"
)

func (app *application) showTrackHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	track, err := app.models.Track.Find(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, payload{"track": track}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}

func (app *application) createTrackHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string        `json:"title"`
		Duration time.Duration `json:"duration"`
		Genres   []string      `json:"genres"`
		Albums   []string      `json:"albums"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	track := &data.Track{
		Title:    input.Title,
		Duration: input.Duration,
		Genres:   input.Genres,
		Albums:   input.Albums,
	}

	v := validator.New()

	if data.ValidateTrack(v, track); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Track.Create(track)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/tracks/%d", track.ID))

	err = app.writeJSON(w, http.StatusCreated, payload{"track": track}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTrackHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	track, err := app.models.Track.Find(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflitResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// if r.Header.Get("X-Expected-Version") != "" {
	// 	if track.Version.String() != r.Header.Get("X-Expected-Version") {
	// 		app.editConflitResponse(w, r)
	// 		return
	// 	}
	// }

	var input struct {
		Title    *string        `json:"title"`
		Duration *time.Duration `json:"duration"`
		Genres   []string       `json:"genres"`
		Albums   []string       `json:"albums"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		track.Title = *input.Title
	}
	if input.Duration != nil {
		track.Duration = *input.Duration
	}
	if input.Genres != nil {
		track.Genres = input.Genres
	}
	if input.Albums != nil {
		track.Albums = input.Albums
	}

	v := validator.New()

	if data.ValidateTrack(v, track); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Track.Update(track)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflitResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, payload{"track": track}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTrackHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Track.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, payload{"message": "track successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listTracksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		Albums string
		data.Filter
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Albums = app.readString(qs, "albums", "")
	input.Genres = app.readCSV(qs, "genres", []string{})

	input.Filter.Page = app.readInt(qs, "page", 1, v)
	input.Filter.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filter.Sort = app.readString(qs, "sort", "id")
	input.Filter.SortSafelist = []string{"id", "title", "genres", "albums", "-id", "-title", "-genres", "-albums"}

	if data.ValidateFilters(v, input.Filter); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	tracks, metadata, err := app.models.Track.FindAll(input.Title, input.Genres, input.Filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, payload{"tracks": tracks, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
