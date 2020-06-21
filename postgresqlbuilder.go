package httptesting

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/elliotchance/orderedmap"
)

type PostgresSqlBuilder struct {
	insertFlag              bool
	updateFlag              bool
	deleteFlag              bool
	selectFlag              bool
	tableName               string
	setParams               *orderedmap.OrderedMap
	whereParams             *orderedmap.OrderedMap
	whereParamsRelationship *orderedmap.OrderedMap
	orderByParams           *orderedmap.OrderedMap
	argumentNames           []string
	argumentValues          []interface{}
	returningParams         []string
	limit                   int
	err                     error
	buffer                  StringBuilder
}

// PostgresSqlBuilder ;= PostgresSqlBuilder{}
// PostgresSqlBuilder.Insert("table1").Set(map[string]interface{}{"param1": 1, "param2": true}).String()

func (s *PostgresSqlBuilder) Insert(tableName string) *PostgresSqlBuilder {
	s.tableName = tableName
	s.insertFlag = true
	s.setParams = orderedmap.NewOrderedMap()
	s.whereParams = orderedmap.NewOrderedMap()
	s.whereParamsRelationship = orderedmap.NewOrderedMap()

	if len(s.tableName) == 0 {
		s.err = errors.New("Table name has to be specified")
	}

	return s
}

func (s *PostgresSqlBuilder) Delete(tableName string) *PostgresSqlBuilder {
	s.tableName = tableName
	s.deleteFlag = true
	s.whereParams = orderedmap.NewOrderedMap()
	s.whereParamsRelationship = orderedmap.NewOrderedMap()

	if len(s.tableName) == 0 {
		s.err = errors.New("Table name has to be specified")
	}
	return s
}

func (s *PostgresSqlBuilder) Update(tableName string) *PostgresSqlBuilder {
	s.tableName = tableName
	s.updateFlag = true
	s.setParams = orderedmap.NewOrderedMap()
	s.whereParams = orderedmap.NewOrderedMap()
	s.whereParamsRelationship = orderedmap.NewOrderedMap()

	if len(s.tableName) == 0 {
		s.err = errors.New("Table name has to be specified")
	}

	return s
}

func (s *PostgresSqlBuilder) Select(tableName string) *PostgresSqlBuilder {

	s.tableName = strings.TrimSpace(tableName)
	s.selectFlag = true
	s.whereParams = orderedmap.NewOrderedMap()
	s.whereParamsRelationship = orderedmap.NewOrderedMap()
	s.orderByParams = orderedmap.NewOrderedMap()

	if len(s.tableName) == 0 {
		s.err = errors.New("Table name has to be specified")
	}

	return s
}

func (s *PostgresSqlBuilder) Where(params *orderedmap.OrderedMap) *PostgresSqlBuilder {
	s.whereParams = params
	for _, k := range params.Keys() {
		s.whereParamsRelationship.Set(k, "=")
	}
	return s
}

func (s *PostgresSqlBuilder) WhereArg(param string, value interface{}) *PostgresSqlBuilder {
	if s.whereParams == nil {
		s.err = errors.New("In this mode usage of WhereArg is not appropriate")
	} else {
		s.whereParams.Set(param, value)
		s.whereParamsRelationship.Set(param, "=")
	}
	return s
}

func (s *PostgresSqlBuilder) WhereArgRelationship(param string, relationship string, value interface{}) *PostgresSqlBuilder {
	if s.whereParams == nil {
		s.err = errors.New("In this mode usage of WhereArgRelationship is not appropriate")
	} else {
		relationship = strings.TrimSpace(relationship)
		if len(relationship) == 0 {
			s.err = errors.New("Relationship is not defined")
		}
		s.whereParams.Set(param, value)
		s.whereParamsRelationship.Set(param, relationship)
	}
	return s
}

func (s *PostgresSqlBuilder) OrderBy(param string, direction string) *PostgresSqlBuilder {
	if s.whereParams == nil || !s.selectFlag {
		s.err = errors.New("In this mode usage of OrderBy is not appropriate")
	} else {
		param = strings.TrimSpace(param)
		direction = strings.TrimSpace(direction)
		if len(param) == 0 {
			s.err = errors.New("OrderBy Param is not defined")
		}
		if len(direction) == 0 {
			s.err = errors.New("OrderBy direction is not defined")
		}
		s.orderByParams.Set(param, direction)
	}
	return s
}

func (s *PostgresSqlBuilder) Set(params *orderedmap.OrderedMap) *PostgresSqlBuilder {
	s.setParams = params
	return s
}

func (s *PostgresSqlBuilder) SetArg(param string, value interface{}) *PostgresSqlBuilder {
	if s.setParams == nil {
		s.err = errors.New("In this mode usage of SetArg is not appropriate")
	} else {
		s.setParams.Set(param, value)
	}
	return s
}

func (s *PostgresSqlBuilder) Returning(params ...string) *PostgresSqlBuilder {
	s.returningParams = params
	return s
}

func (s *PostgresSqlBuilder) Limit(limit int) *PostgresSqlBuilder {
	s.limit = limit
	return s
}

func buildValuesClause(s *PostgresSqlBuilder) string {
	sb := StringBuilder{}
	sb1 := StringBuilder{}

	if s.setParams.Len() == 0 {
		s.err = errors.New("No insertion parameters passed")
	}

	sb.Write("(")
	sb1.Write("VALUES (")
	index := 1
	argCount := len(s.argumentNames)
	for _, name := range s.setParams.Keys() {
		value, ok := s.setParams.Get(name)
		if !ok {
			continue
		}
		if index > 1 {
			sb.Write(", ")
			sb1.Write(", ")
		}
		sb.Write(name.(string))
		sb1.Write("$", strconv.Itoa(argCount+index))
		s.argumentNames = append(s.argumentNames, name.(string))
		s.argumentValues = append(s.argumentValues, value)
		index++
	}
	sb.Write(") ")
	sb1.Write(")")
	sb.Write(sb1.String())
	return sb.String()
}

func buildSetClause(s *PostgresSqlBuilder) string {
	sb := StringBuilder{}

	if s.setParams.Len() == 0 {
		s.err = errors.New("No update parameters passed")
	}

	index := 1
	argCount := len(s.argumentNames)
	for _, name := range s.setParams.Keys() {
		value, ok := s.setParams.Get(name)
		if !ok {
			continue
		}

		if index > 1 {
			sb.Write(",")
		}

		sb.Write(name.(string), "=$", strconv.Itoa(argCount+index))
		s.argumentNames = append(s.argumentNames, name.(string))
		s.argumentValues = append(s.argumentValues, value)
		index++
	}
	return sb.String()
}

func buildWhereClause(s *PostgresSqlBuilder) string {
	sb := StringBuilder{}

	if s.whereParams.Len() == 0 {
		s.err = errors.New("Where clause has to be specified")
	}

	index := 1
	argCount := len(s.argumentNames)
	for _, name := range s.whereParams.Keys() {
		value, ok := s.whereParams.Get(name)
		if !ok {
			s.err = errors.New("Incomplete where arguments")
			continue
		}

		if index > 1 {
			sb.Write(" AND ")
		}

		relationship, ok := s.whereParamsRelationship.Get(name.(string))
		if !ok {
			s.err = errors.New("Incomplete where relationships")
			continue
		}

		sb.Write(name.(string), relationship.(string), "$", strconv.Itoa(argCount+index))
		s.argumentNames = append(s.argumentNames, name.(string))
		s.argumentValues = append(s.argumentValues, value)
		index++
	}

	return sb.String()
}

func buildOrderByClause(s *PostgresSqlBuilder) string {
	sb := StringBuilder{}
	if s.orderByParams.Len() > 0 {
		sb.Write(" ORDER BY ")
		for index, name := range s.orderByParams.Keys() {
			value, ok := s.orderByParams.Get(name)
			if !ok {
				s.err = errors.New("Incomplete order by arguments")
				continue
			}
			if index > 0 {
				sb.Write(", ")
			}
			sb.Write(name.(string), " ", value.(string))
		}

	}
	return sb.String()
}

func buildReturnClause(s *PostgresSqlBuilder) string {
	sb := StringBuilder{}

	if len(s.returningParams) > 0 {
		for i, name := range s.returningParams {
			if i > 0 {
				sb.Write(", ")
			}
			sb.Write(name)
		}
	}
	return sb.String()
}

func buildSelectClause(s *PostgresSqlBuilder) string {
	sb := StringBuilder{}

	if len(s.returningParams) == 0 {
		s.err = errors.New("Returning params have to be specified")
	}

	if len(s.returningParams) > 0 {
		for i, name := range s.returningParams {
			if i > 0 {
				sb.Write(", ")
			}
			sb.Write(name)
		}
	}
	return sb.String()
}

func buildLimitClause(s *PostgresSqlBuilder) string {
	sb := StringBuilder{}

	if !s.selectFlag {
		s.err = errors.New("Limit clause only supported for select")
	}

	if s.limit > 0 {
		sb.Write(" LIMIT ").Write(strconv.Itoa(s.limit))
	}
	return sb.String()
}

func (s *PostgresSqlBuilder) Build() (string, []interface{}, []string, error) {
	if s.selectFlag {
		s.buffer.Write("SELECT ", buildSelectClause(s), " FROM ", s.tableName, " ")
		s.buffer.Write("WHERE ")
		s.buffer.Write(buildWhereClause(s))
		s.buffer.Write(buildOrderByClause(s))
		s.buffer.Write(buildLimitClause(s))
	} else if s.insertFlag {
		s.buffer.Write("INSERT INTO ", s.tableName, " ")
		s.buffer.Write(buildValuesClause(s))
		if len(s.returningParams) > 0 {
			s.buffer.Write(" RETURNING ")
			s.buffer.Write(buildReturnClause(s))
		}
	} else if s.updateFlag {

		s.buffer.Write("UPDATE ", s.tableName, " SET ")
		s.buffer.Write(buildSetClause(s))
		s.buffer.Write(" WHERE ")
		s.buffer.Write(buildWhereClause(s))
		if len(s.returningParams) > 0 {
			s.buffer.Write(" RETURNING ")
			s.buffer.Write(buildReturnClause(s))
		}
	} else if s.deleteFlag {
		s.buffer.Write("DELETE FROM ", s.tableName, " ")
		s.buffer.Write("WHERE ")
		s.buffer.Write(buildWhereClause(s))
		if len(s.returningParams) > 0 {
			s.buffer.Write(" RETURNING ")
			s.buffer.Write(buildReturnClause(s))
		}
	}

	if s.err != nil {
		return "", nil, nil, s.err
	}
	return s.buffer.String(), s.argumentValues, s.returningParams, s.err
}

func ScanOneToMap(rows *sql.Rows) (map[string]interface{}, error) {
	res, err := ScanToMap(rows, 1)
	if err != nil {
		return nil, err
	}

	if len(res) < 1 {
		return nil, errors.New("No matches found")
	}
	return res[0], nil
}

func ScanToMap(rows *sql.Rows, limit int) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	res := make([]map[string]interface{}, 0)
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		res = append(res, m)
		if limit != 0 && len(res) >= limit {
			break
		}
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return res, nil
}
