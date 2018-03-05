package messages

import (
	"bytes"
	"context"
	"database/sql"
	"strings"
)

// Execer lets functions accept a DB or a Tx without knowing the difference
type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// A Queue represents a Vitess message queue
type Queue struct {
	name       string
	fieldNames []string

	// predefine these sql strings
	insertSQL          string
	insertScheduledSQL string

	s *subscription
}

// NewQueue returns a queue definition
func NewQueue(ctx context.Context, name string, fieldNames []string) (*Queue, error) {
	q := &Queue{
		name:       name,
		fieldNames: fieldNames,
	}

	// only do this string manipulation once
	q.insertSQL = q.generateInsertSQL()
	q.insertScheduledSQL = q.generateInsertScheduledSQL()

	return q, nil
}

// MustNewQueue returns a queue definition and panics if there is an error
func MustNewQueue(ctx context.Context, name string, fieldNames []string) *Queue {
	q, err := NewQueue(ctx, name, fieldNames)
	if err != nil {
		panic(err)
	}
	return q
}

// generateInsertSQL does the string manipulation to generate the insert statement
func (q *Queue) generateInsertSQL() string {
	buf := bytes.Buffer{}

	// generate default insert into queue with required fields
	buf.WriteString("INSERT INTO `")
	buf.WriteString(q.name)
	buf.WriteString("` (id")

	// add quoted user fields to the insert statement
	for _, f := range q.fieldNames {
		buf.WriteString(", `")
		buf.WriteString(f)
		buf.WriteString("`")
	}
	buf.WriteString(") VALUES (?")

	// add params representing user data
	buf.WriteString(strings.Repeat(",?", len(q.fieldNames)))

	// close VALUES
	buf.WriteString(")")

	return buf.String()
}

// generateInsertScheduledSQL does the string manipulation to generate the insertFuture statement
func (q *Queue) generateInsertScheduledSQL() string {
	buf := bytes.Buffer{}

	// generate default insert into queue with required fields
	buf.WriteString("INSERT INTO `")
	buf.WriteString(q.name)
	buf.WriteString("` (time_scheduled, id")

	// add quoted user fields to the insert statement
	for _, f := range q.fieldNames {
		buf.WriteString(", `")
		buf.WriteString(f)
		buf.WriteString("`")
	}
	buf.WriteString(") VALUES (?,?")

	// add params representing user data
	buf.WriteString(strings.Repeat(",?", len(q.fieldNames)))

	// close VALUES
	buf.WriteString(")")

	return buf.String()
}
