package repository

type PrimaryKey = interface {
	int64 | string
}

type Entity[P PrimaryKey] interface {
	ToMap() map[string]interface{}
}

type Model[P PrimaryKey, E Entity[P]] interface {
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
