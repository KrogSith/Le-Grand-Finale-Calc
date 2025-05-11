package application

import (
	"calculator/internal/calculation"
	"encoding/json"
	"fmt"
	"net/http"

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

	if !UserIsAuthorised() {
		w.WriteHeader(http.StatusUnauthorized)
		response := struct {
			Error string `json:"error"`
		}{
			Error: "Please, sign in first",
		}

		jsonData, _ := json.MarshalIndent(response, "", "  ")

		w.Write(jsonData)
		return
	}

	ctx, db, err := InitiateDatabase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
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

		expression := NewExpression(request.Expression, current_user, err.Error(), result)
		err := insertExpression(ctx, db, expression)
		if err != nil {
			panic(err)
		}

		response := struct {
			Error string `json:"error"`
		}{
			Error: expression.Status,
		}

		jsonData, _ := json.MarshalIndent(response, "", "  ")

		w.Write(jsonData)
		return
	}
	w.WriteHeader(http.StatusCreated)

	expression := NewExpression(request.Expression, current_user, "OK", result)
	err = insertExpression(ctx, db, expression)
	if err != nil {
		panic(err)
	}

	response := struct {
		Id string `json:"id"`
	}{
		Id: expression.ID,
	}

	jsonData, _ := json.MarshalIndent(response, "", "  ")

	w.Write(jsonData)
}

func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if !UserIsAuthorised() {
		w.WriteHeader(http.StatusUnauthorized)
		response := struct {
			Error string `json:"error"`
		}{
			Error: "Please, sign in first",
		}

		jsonData, _ := json.MarshalIndent(response, "", "  ")

		w.Write(jsonData)
		return
	}

	ctx, db, err := InitiateDatabase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	expressions, err := selectExpressions(ctx, db)
	if err != nil {
		panic(err)
	}

	response := map[string]interface{}{
		"expressions": expressions,
	}

	jsonData, _ := json.MarshalIndent(response, "", "  ")

	w.Write(jsonData)
}

func IDHandler(w http.ResponseWriter, r *http.Request) {
	if !UserIsAuthorised() {
		w.WriteHeader(http.StatusUnauthorized)
		response := struct {
			Error string `json:"error"`
		}{
			Error: "Please, sign in first",
		}

		jsonData, _ := json.MarshalIndent(response, "", "  ")

		w.Write(jsonData)
		return
	}

	ctx, db, err := InitiateDatabase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	expressions, err := selectExpressions(ctx, db)
	if err != nil {
		panic(err)
	}

	vars := mux.Vars(r)
	id := vars["id"]

	for _, expression := range expressions {
		if expression.ID == id {
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

type Register_Request struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func RegistrationHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Register_Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, db, err := InitiateDatabase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	users, err := selectUsers(ctx, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	for _, user := range users {
		if user.Name == request.Name {
			w.WriteHeader(http.StatusNotAcceptable)
			response := struct {
				Error string `json:"error"`
			}{
				Error: "User with this name already exists",
			}

			jsonData, _ := json.MarshalIndent(response, "", "  ")

			w.Write(jsonData)

			return
		}
	}

	user := NewUser(request.Name, request.Password)
	err = insertUser(ctx, db, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	response := struct {
		Message string `json:"message"`
	}{
		Message: fmt.Sprintf("Registration succesful. Welcome, %v!", user.Name),
	}

	jsonData, _ := json.MarshalIndent(response, "", "  ")

	w.Write(jsonData)
	current_user = user.ID
}

func LogInHandler(w http.ResponseWriter, r *http.Request) {
	request := new(Register_Request)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, db, err := InitiateDatabase()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	users, err := selectUsers(ctx, db)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}

	for _, user := range users {
		if user.Name == request.Name {
			if user.Password != request.Password {
				w.WriteHeader(http.StatusNotAcceptable)
				response := struct {
					Error string `json:"error"`
				}{
					Error: "Wrong password",
				}

				jsonData, _ := json.MarshalIndent(response, "", "  ")

				w.Write(jsonData)
				return
			}

			user, err := selectUserByID(ctx, db, user.ID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				panic(err)
			}

			w.WriteHeader(http.StatusOK)
			response := struct {
				Message string `json:"message"`
			}{
				Message: fmt.Sprintf("Welcome, %v!", user.Name),
			}

			jsonData, _ := json.MarshalIndent(response, "", "  ")

			w.Write(jsonData)
			current_user = user.ID
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	response := struct {
		Error string `json:"error"`
	}{
		Error: "User not found",
	}

	jsonData, _ := json.MarshalIndent(response, "", "  ")

	w.Write(jsonData)
}
