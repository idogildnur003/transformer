package models

import "sync"

type UserModel struct {
	Users      []map[string]interface{}
	Activities []map[string]interface{}
	UserMutex  sync.Mutex
}

type DefaultStructure struct {
	Id             string          `json:"id"`
	ExternalID     string          `json:"external_id"`
	Mail           string          `json:"mail"`
	Type           string          `json:"type"`
	Location       string          `json:"location"`
	IsEnabled      bool            `json:"is_enabled"`
	FirstName      string          `json:"first_name"`
	LastName       string          `json:"last_name"`
	SignInActivity *SignInActivity `json:"sign_in_activity,omitempty"`
}

type SignInActivity struct {
	LastSignInDateTime                string `json:"lastSignInDateTime"`
	LastSignInRequestId               string `json:"lastSignInRequestId"`
	LastNonInteractiveSignInDateTime  string `json:"lastNonInteractiveSignInDateTime"`
	LastNonInteractiveSignInRequestId string `json:"lastNonInteractiveSignInRequestId"`
	LastSuccessfulSignInDateTime      string `json:"lastSuccessfulSignInDateTime"`
	LastSuccessfulSignInRequestId     string `json:"lastSuccessfulSignInRequestId"`
}

const DEFAULT_RULES_MODEL = `
{
	"id": "id",
	"external_id": "userPrincipalName",
	"mail": "mail",
	"type": "userType",
	"location": "usageLocation",
	"is_enabled": "accountEnabled",
	"first_name": "givenName",
	"last_name": "surname",
	"sign_in_activity": {
		"lastSignInDateTime": "signInActivity.lastSignInDateTime",
		"lastSignInRequestId": "signInActivity.lastSignInRequestId",
		"lastNonInteractiveSignInDateTime": "signInActivity.lastNonInteractiveSignInDateTime",
		"lastNonInteractiveSignInRequestId": "signInActivity.lastNonInteractiveSignInRequestId",
		"lastSuccessfulSignInDateTime": "signInActivity.lastSuccessfulSignInDateTime",
		   "lastSuccessfulSignInRequestId": "signInActivity.lastSuccessfulSignInRequestId"
	}
}`
