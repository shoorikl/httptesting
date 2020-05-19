package httptesting

import (
	"errors"
	"strconv"
	"strings"

	"github.com/elliotchance/orderedmap"
)

type PostgresSqlBuilder struct {
	insertFlag      bool
	updateFlag      bool
	deleteFlag      bool
	selectFlag      bool
	tableName       string
	setParams       *orderedmap.OrderedMap
	whereParams     *orderedmap.OrderedMap
	argumentNames   []string
	argumentValues  []interface{}
	returningParams []string
	err             error
	buffer          StringBuilder
}

// PostgresSqlBuilder ;= PostgresSqlBuilder{}
// PostgresSqlBuilder.Insert("table1").Set(map[string]interface{}{"param1": 1, "param2": true}).String()

func (s *PostgresSqlBuilder) Insert(tableName string) *PostgresSqlBuilder {
	s.tableName = tableName
	s.insertFlag = true
	s.setParams = orderedmap.NewOrderedMap()
	s.whereParams = orderedmap.NewOrderedMap()

	if len(s.tableName) == 0 {
		s.err = errors.New("Table name has to be specified")
	}

	return s
}

func (s *PostgresSqlBuilder) Delete(tableName string) *PostgresSqlBuilder {
	s.tableName = tableName
	s.deleteFlag = true
	s.whereParams = orderedmap.NewOrderedMap()

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

	if len(s.tableName) == 0 {
		s.err = errors.New("Table name has to be specified")
	}

	return s
}

func (s *PostgresSqlBuilder) Select(tableName string) *PostgresSqlBuilder {

	s.tableName = strings.TrimSpace(tableName)
	s.selectFlag = true
	s.whereParams = orderedmap.NewOrderedMap()

	if len(s.tableName) == 0 {
		s.err = errors.New("Table name has to be specified")
	}

	return s
}

func (s *PostgresSqlBuilder) Where(params *orderedmap.OrderedMap) *PostgresSqlBuilder {
	s.whereParams = params
	return s
}

func (s *PostgresSqlBuilder) WhereArg(param string, value interface{}) *PostgresSqlBuilder {
	s.whereParams.Set(param, value)
	return s
}

func (s *PostgresSqlBuilder) Set(params *orderedmap.OrderedMap) *PostgresSqlBuilder {
	s.setParams = params
	return s
}

func (s *PostgresSqlBuilder) SetArg(param string, value interface{}) *PostgresSqlBuilder {
	s.setParams.Set(param, value)
	return s
}

func (s *PostgresSqlBuilder) Returning(params ...string) *PostgresSqlBuilder {
	s.returningParams = params
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
			continue
		}

		if index > 1 {
			sb.Write(" AND ")
		}

		sb.Write(name.(string), "=$", strconv.Itoa(argCount+index))
		s.argumentNames = append(s.argumentNames, name.(string))
		s.argumentValues = append(s.argumentValues, value)
		index++
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

func (s *PostgresSqlBuilder) Build() (string, []interface{}, []string, error) {
	if s.selectFlag {
		s.buffer.Write("SELECT ", buildSelectClause(s), " FROM ", s.tableName, " ")
		s.buffer.Write("WHERE ")
		s.buffer.Write(buildWhereClause(s))
	} else if s.insertFlag {
		s.buffer.Write("INSERT INTO ", s.tableName, " ")
		s.buffer.Write(buildValuesClause(s))
		if len(s.returningParams) > 0 {
			s.buffer.Write(" RETURNING ")
			s.buffer.Write(buildReturnClause(s))
		}
	} else if s.updateFlag {

		s.buffer.Write("UPDATE ", s.tableName, " ")
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
