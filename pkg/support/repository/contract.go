package repository

type Entity interface {
	ToMap() map[string]interface{}
}

type Model[E Entity] interface {
	ToEntity() *E
	FromEntity(e E) interface{}
	TableName() string
}

type Condition interface {
	GetQuery() string
	GetValues() []any
}

type Preload interface {
	GetRelation() string
	GetCondition() *Condition
	GetSelectColumns() []string
}
