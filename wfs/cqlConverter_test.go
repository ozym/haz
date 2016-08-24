package main

import (
	"bytes"
	"strings"
	"testing"
)

/**
 * test the cqlConverter
 */
func TestGetParenthesisedTokens(t *testing.T) {
	cqlString := `(origin_geom, 174,-41,175, -42)`
	cql := NewCqlConverter(cqlString)
	cql.NextToken()
	tokens := cql.GetParenthesisedTokens()
	//t.Log(tokens)
	if len(tokens) != 6 {
		t.Fail()
	}
}

func TestToBBoxSql(t *testing.T) {
	cqlString := `BBOX(origin_geom, 174,-41,175, -42)`
	expected := `ST_Contains(ST_SetSRID(ST_Envelope($1::geometry),4326),ST_Shift_Longitude(origin_geom))`
	expectedArg := "'LINESTRING(174 -41,175 -42)'"
	var sql bytes.Buffer
	cql := NewCqlConverter(cqlString)
	cql.NextToken()
	args := []interface{}{}
	args, err := cql.ToBBoxSql(&sql, args)

	if err != nil || strings.Trim(sql.String(), " ") != expected {
		t.Error(err)
		t.Error(sql.String(), " not equals expected", expected)
	}
	if len(args) != 1 && args[0] != expectedArg {
		t.Error("args is not expected", args)
	}
}

func TestToWithinSql(t *testing.T) {
	cqlString := `WITHIN(origin_geom,POLYGON((172.951 -41.767,172.001 -42.832,169.564 -44.341,172.312 -45.412,175.748 -42.908,172.951 -41.767)))`
	expected := `ST_Within(origin_geom, ST_GeomFromText($1, 4326))`
	var sql bytes.Buffer
	args := []interface{}{}
	cql := NewCqlConverter(cqlString)
	cql.NextToken()
	args, err := cql.ToWithinSql(&sql, args)

	if err != nil {
		t.Error(err)
	}

	if strings.Trim(sql.String(), " ") != expected {
		t.Error(sql.String(), " not equals expected", expected)
	}
	expectedArg := "POLYGON((172.951 -41.767,172.001 -42.832,169.564 -44.341,172.312 -45.412,175.748 -42.908,172.951 -41.767))"
	if len(args) != 1 && args[0] != expectedArg {
		t.Error("args is not expected", args)
	}
}

func TestToDWithinSql(t *testing.T) {
	cqlString := `DWITHIN(origin_geom,POINT(172.951 -41.767),5000,meters)`
	expected := `ST_DWithin(origin_geom::Geography, ST_GeomFromText($1, 4326)::Geography, $2)`
	var sql bytes.Buffer
	args := []interface{}{}
	cql := NewCqlConverter(cqlString)
	cql.NextToken()
	args, err := cql.ToDWithinSql(&sql, args)

	if err != nil {
		t.Error(err)
	}

	if sql.String() != expected {
		t.Error(sql.String(), " not equals expected", expected)
	}

	if len(args) != 2 {
		t.Error("args is not expected", args)
	}
	if args[0].(string) != "POINT(172.951 -41.767)" {
		t.Error("args is not expected", args[0])
	}
	if args[1].(float64) != 5000 {
		t.Error("args is not expected", args[1])
	}

}

func TestToSql(t *testing.T) {
	cqlString := `(origintime>='2013-06-01' AND origintime<'2016-04-12T22:00:00') or usedphasecount != 60`
	expected := `(origintime >= $1::timestamptz AND origintime < $2::timestamptz) or usedphasecount != $3`
	cql := NewCqlConverter(cqlString)
	sql, args, err := cql.ToSQL()
	//t.Log(sql)
	//t.Log(args)
	if err != nil {
		t.Error(err)
	}
	if strings.Trim(sql, " ") != expected {
		t.Error("sql string not equals expected", sql, expected)
	}

	if len(args) != 3 {
		t.Error("args should have 2 parameters")
	}
	if args[0] != "'2013-06-01'" || args[1] != "'2016-04-12T22:00:00'" || args[2] != "60" {
		t.Error("args should be 2013-06-01, 2016-04-12T22:00:00 ", args)
	}
}
