package application

import (
	"calculator/internal/calculation"
	"calculator/internal/modules"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")
	if config.Addr == "" {
		config.Addr = "8080"
	}
	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (a *Application) Run() error {
	// for {
	// 	log.Println("Input expression:")
	// 	reader := bufio.NewReader(os.Stdin)
	// 	text, err := reader.ReadString('\n')
	// 	text = text[:len(text)-2]
	// 	if err != nil {
	// 		fmt.Println("Failed to read application from console")
	// 	}
	// 	result, err := calculation.Calc(text)
	// 	if text == "exit" {
	// 		log.Println(text, " calculation failed with error: ", err)
	// 	} else {
	// 		log.Println(text, "=", result)
	// 	}
	// }
	text := "(1-(2+3))+2-1"
	result, _ := calculation.Calc(text)
	log.Println(text, "=", result)

	return nil
}

type Request struct {
	Expression string `json:"expression"`
	Id         string `json:"id"`
}

var expressions modules.Expressions

var COMPUTING_POWER int = 3

func Agent(req string) {
	for i := 0; i < COMPUTING_POWER; i++ {
		go Worker(req)
	}
}

func Worker(req string) {
	for {
		ch := make(chan float64)
		result, err := calculation.Calc(req)
		if err != nil {
			return
		}
		ch <- result
	}
}

func (a *Application) RunServer() error {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/calculate", CalcHandler)
	r.HandleFunc("/api/v1/expressions", ExpressionsHandler)
	r.HandleFunc("/api/v1/expressions/{id}", IDHandler).Methods("GET")
	return http.ListenAndServe(":"+a.config.Addr, r)
}
