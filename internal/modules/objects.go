package modules

type Expression struct {
	Id     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

func NewExpression(id string, status string, result float64) Expression {
	return Expression{
		Id:     id,
		Status: status,
		Result: result,
	}
}

type Expressions []Expression

type Task struct {
	Operation      string `json:"operation"`
	Operation_time int    `json:"operation_time"`
}

type Tasks []Task
