package tasks

import "encoding/json"

const DLPGraphQL = "https://api.disneylandparis.com/query"

type GraphQlQuery struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func (q GraphQlQuery) ToJSON() ([]byte, error) {
	return json.Marshal(q)
}
