package httptesting

import (
	"database/sql"
)

func CountRows(db *sql.DB, tableName string) (int, error) {
	builder := PostgresSqlBuilder{}
	query, inArgs, _, err := builder.Select(tableName).Returning("count(id)").
		WhereArgRelationship("id", ">=", 0).Build()
	if err != nil {
		return -1, err
	}

	var count int
	err = db.QueryRow(query, inArgs...).Scan(&count)
	if err != nil {
		return -1, err
	}

	return count, nil
}
