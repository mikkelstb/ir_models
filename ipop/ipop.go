package ipop

import irmodels "github.com/mikkelstb/ir_models"


type IO interface {

	Init(preferences map[string]string)
	GetNextDoc() *irmodels.Article
	HasNext() (bool)
}

