package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/model"
	"github.com/margulan-kalykul/JustQuiz/pkg/quiz/validator"
)

func (app *application) createQuizHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Category  string   `json:"category"`
		Reward    int      `json:"reward"`
		Questions []string `json:"questions"`
		Answers   []string `json:"answers"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		log.Println(err)
		app.errorResponse(w, r, http.StatusBadRequest, "Invalid request payload")
		return
	}

	quiz := &model.Quiz{
		Category:  input.Category,
		Reward:    input.Reward,
		Questions: input.Questions,
		Answers:   input.Answers,
	}

	err = app.models.Quizes.Insert(quiz)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"quiz": quiz}, nil)
}

func (app *application) getQuizesList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Category   string
		RewardFrom int
		RewrdTo    int
		model.Filters
	}
	v := validator.New()
	qs := r.URL.Query()

	// Use our helpers to extract the name and score value range query string values, falling back to the
	// defaults of an empty string and an empty slice, respectively, if they are not provided
	// by the client.
	input.Category = app.readStrings(qs, "category", "")
	input.RewardFrom = app.readInt(qs, "rewardFrom", 0, v)
	input.RewrdTo = app.readInt(qs, "rewardTo", 0, v)

	// Ge the page and page_size query string value as integers. Notice that we set the default
	// page value to 1 and default page_size to 20, and that we pass the validator instance
	// as the final argument.
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply an ascending sort on quiz ID).
	input.Filters.Sort = app.readStrings(qs, "sort", "id")

	// Add the supported sort value for this endpoint to the sort safelist.
	// name of the column in the database.
	input.Filters.SortSafeList = []string{
		// ascending sort values
		"id", "category", "reward",
		// descending sort values
		"-id", "-category", "-reward",
	}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	quizes, metadata, err := app.models.Quizes.GetAll(input.Category, input.RewardFrom, input.RewrdTo, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"quizes": quizes, "metadata": metadata}, nil)
}

func (app *application) getQuizHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	quiz, err := app.models.Quizes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"quiz": quiz}, nil)
}

func (app *application) updateQuizHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	quiz, err := app.models.Quizes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Category  *string   `json:"category"`
		Reward    *int      `json:"reward"`
		Questions *[]string `json:"questions"`
		Answers   *[]string `json:"answers"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check fileds
	if input.Category != nil {
		quiz.Category = *input.Category
	}
	if input.Reward != nil {
		quiz.Reward = *input.Reward
	}

	v := validator.New()

	if model.ValidateQuiz(v, quiz); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Quizes.Update(quiz)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"quiz": quiz}, nil)
}

func (app *application) deleteQuizHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Quizes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.writeJSON(w, http.StatusOK, envelope{"message": "success"}, nil)
}
