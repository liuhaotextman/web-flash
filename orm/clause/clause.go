package clause

type Type int

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type]interface{}
}

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)
