package models

import "github.com/oxyno-zeta/golang-graphql-example/pkg/golang-graphql-example/database/common"

type SortOrder struct {
	CreatedAt *common.SortOrderEnum `dbfield:"created_at"`
	UpdatedAt *common.SortOrderEnum `dbfield:"updated_at"`
	Text      *common.SortOrderEnum `dbfield:"text"`
	Done      *common.SortOrderEnum `dbfield:"done"`
}

type Filter struct {
	CreatedAt *common.DateFilter    `dbfield:"created_at"`
	UpdatedAt *common.DateFilter    `dbfield:"updated_at"`
	Text      *common.GenericFilter `dbfield:"text"`
	Done      *common.GenericFilter `dbfield:"done"`
}