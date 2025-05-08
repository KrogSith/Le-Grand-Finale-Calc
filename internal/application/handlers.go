package application

import (
	"calculator/internal/calculation"
	"calculator/internal/modules"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	res := make(chan float64)
	er := make(chan error)
	request := new(Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Agent(request.Expression)

	for i := 0; i < COMPUTING_POWER; i++ {
		go func() {
			for {
				result, err := calculation.Calc(request.Expression)
				if err != nil {
					w.WriteHeader(http.StatusUnprocessableEntity)
				}
				res <- result
				er <- err
			}
		}()
	}
	result := <-res
	err = <-er

	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)

		id := uuid.NewString()
		expression := modules.NewExpression(id, request.Expression, err.Error(), result)
		expressions = append(expressions, expression)

		response := struct {
			Error string `json:"error"`
		}{
			Error: expression.Status,
		}

		jsonData, _ := json.MarshalIndent(response, "", "  ")

		w.Write(jsonData)
	} else {
		w.WriteHeader(http.StatusCreated)

		id := uuid.NewString()
		expression := modules.NewExpression(id, request.Expression, "OK", result)
		expressions = append(expressions, expression)

		response := struct {
			Id string `json:"id"`
		}{
			Id: expression.Id,
		}

		jsonData, _ := json.MarshalIndent(response, "", "  ")

		w.Write(jsonData)
	}
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"expressions": expressions,
	}

	jsonData, _ := json.MarshalIndent(response, "", "  ")

	w.Write(jsonData)
}

func IDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	for _, expression := range expressions {
		if expression.Id == id {
			jsonData, err := json.MarshalIndent(expression, "", "  ")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(jsonData)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
