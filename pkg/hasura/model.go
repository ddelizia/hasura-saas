package hasura

type BasePayload struct {
	SessionVariables map[string]interface{} `json:"session_variables"`
}
