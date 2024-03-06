package repository

type preload struct {
	relation  string
	condition *Condition
	columns   []string
}

func (preload preload) GetRelation() string {
	return preload.relation
}

func (preload preload) GetCondition() *Condition {
	return preload.condition
}

func (preload preload) GetSelectColumns() []string {
	return preload.columns
}

func With(relation string, vars ...interface{}) Preload {
	var condition Condition = nil
	columns := []string{"*"}
	if len(vars) > 0 && vars[0] != nil {
		condition = vars[0].(Condition)
	}
	if len(vars) > 1 && vars[1] != nil {
		columns = vars[1].([]string)
	}
	return preload{
		relation:  relation,
		condition: &condition,
		columns:   columns,
	}
}
