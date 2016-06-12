package main

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/**
 * A cql converter that converts CQL query to SQL query
 * Baishan 15/4/2016
 */
const (
	DATE_PATTERN      = "^(\\d{4})-(\\d{1,2})-(\\d{1,2})$"
	DATE_TIME_PATTERN = "^(\\d{4})-(\\d{1,2})-(\\d{1,2})T(\\d{1,2}):(\\d{1,2}):(\\d{1,2})$"
	NUMBER_PATTERN    = "^-?[.\\d]+$"
	//token types
	TT_EOF          = -1
	TT_WORD         = -3
	TT_STRING       = 999  // quoted string
	TT_SORTBY       = 1008 // The "sortby" operator
	TT_PARENTHESIS  = 1100 //()
	TT_LOGIC_OPS    = 1101 //and or
	TT_RELATION_OPS = 1102 //>, <, =, >=, <=, !=, ==
	TT_TIMESTRING   = 1103 //2016-4-12T22:00:00
	TT_BBOX         = 1104 //BBOX(origin_geom,174,-41,175,-42)
	TT_WITHIN       = 1105 //WITHIN(origin_geom,POLYGON((172.951 -41.767, 172.001 -42.832, 169.564 -44.341,+172.312 -45.412, 175.748 -42.908, 172.951 -41.767)))
	TT_DWITHIN      = 1106 //DWITHIN(origin_geom,Point+(175+-41),0.5,feet)

	//database field type
	WFS_DB_FIELD_TYPE_STRING  = 2001
	WFS_DB_FIELD_TYPE_TIME    = 2002
	WFS_DB_FIELD_TYPE_NUMBER  = 2003
	WFS_DB_FIELD_TYPE_UNKNOWN = 2004
)

type CqlConverter struct {
	cqlCharArray     []rune //char array of the CQL string
	cqlLen           int    //length of the CQL string
	currentIndex     int    //current index in the cqlCharArray
	currentTokenType int    //current token type in the CQL string
	currentToken     string //current token in the CQL
	//immediate previous token, if previous token is TT_RELATION_OPS, current one can be a place holder
	//if current type is TT_RELATION_OPS, previous token need be validated against available fields
	previousTokenType1 int    //immediate previous token type in the CQL string
	previousToken1     string //immediate previous token in the CQL
	//previous previous token
	previousTokenType2 int    //previous previous token type in the CQL string
	previousToken2     string //previous previous token in the CQL
	BBOX               string
}

//search string contains another string
func SearchStringInString(s string, ch string) bool {
	return (strings.Index(s, ch)) >= 0
}

/**
 * convenient pattern match (treat error as false)
 */
func PatternMatch(pattern string, str string) bool {
	match, err := regexp.MatchString(pattern, str)
	if err != nil {
		return false
	}
	return match
}

// NewCqlConverterL returns a pointer to a CqlConverter struct.
func NewCqlConverter(cql string) *CqlConverter {
	cqlConverter := CqlConverter{}
	cqlConverter.cqlCharArray = []rune(cql)
	cqlConverter.cqlLen = len(cqlConverter.cqlCharArray)
	return &cqlConverter
}

func (cql *CqlConverter) Render() string {
	var ret string
	switch cql.currentTokenType {
	case TT_BBOX:
		ret = "bboc"
	case TT_WITHIN:
		ret = "within"
	case TT_LOGIC_OPS:
		ret = "logic"
	case TT_PARENTHESIS:
		ret = "paranthesis"
	case TT_RELATION_OPS:
		ret = "relation"
	case TT_DWITHIN:
		ret = "dwithin"
	case TT_TIMESTRING:
		ret = "time"
	case TT_WORD:
		ret = "word"
	case TT_STRING:
		ret = "string"
	case TT_EOF:
		ret = "END"
	default:
		ret = ""
	}
	ret += ">>" + cql.currentToken
	return ret
}

/**
 * move to next token, update token and type
 * token separators: "()/<>= \t\r\n"
 */
func (cql *CqlConverter) NextToken() {
	//eat whitespace
	for cql.currentIndex < cql.cqlLen && SearchStringInString(" \t\r\n", string(cql.cqlCharArray[cql.currentIndex])) {
		cql.currentIndex++
	}
	//store previous type and tokens
	cql.previousTokenType2 = cql.previousTokenType1
	cql.previousToken2 = cql.previousToken1
	cql.previousTokenType1 = cql.currentTokenType
	cql.previousToken1 = cql.currentToken

	//eof
	if cql.currentIndex == cql.cqlLen {
		cql.currentTokenType = TT_EOF
		cql.currentToken = ""
		return
	}
	//get current char
	c := string(cql.cqlCharArray[cql.currentIndex])
	var buf bytes.Buffer
	var lval string
	//1. separators
	if SearchStringInString("()", c) {
		cql.currentTokenType = TT_PARENTHESIS
		cql.currentToken = c
		cql.currentIndex++
		//2. comparitor
	} else if SearchStringInString("<>=!", c) {
		cql.currentIndex++
		//two-char comparitor
		twoCharOps := false

		//log.Println("2. relation c " + c)

		if cql.currentIndex < cql.cqlLen {
			d := string(cql.cqlCharArray[cql.currentIndex])
			//log.Println("2. relation d " + d)
			if SearchStringInString("<>=", d) {
				cql.currentToken = c + d
				cql.currentTokenType = TT_RELATION_OPS
				twoCharOps = true
				cql.currentIndex++
			}
		}
		if !twoCharOps { //single char ops
			if SearchStringInString("<>=", c) {
				cql.currentToken = c
				cql.currentTokenType = TT_RELATION_OPS
			} else { //only the ! char
				cql.currentIndex--
			}
		}
		//3. quoted string
	} else if SearchStringInString("\"", c) { //no single-quotes
		cql.currentTokenType = TT_STRING
		//remember quote char
		mark := c
		cql.currentIndex++
		escaped := false

		buf.Reset() //reset buffer
		for cql.currentIndex < cql.cqlLen {
			if !escaped && (string(cql.cqlCharArray[cql.currentIndex]) == mark) { //terminator
				break
			}
			if escaped && SearchStringInString("*?^\\", string(cql.cqlCharArray[cql.currentIndex])) { //no escaping for d-quote
				buf.WriteString("\\")
			}
			if !escaped && string(cql.cqlCharArray[cql.currentIndex]) == "\\" { //escape-char
				escaped = true
				cql.currentIndex++
				continue
			}
			escaped = false //reset escape
			buf.WriteString(string(cql.cqlCharArray[cql.currentIndex]))
			cql.currentIndex++
		}
		cql.currentToken = buf.String()
		lval = strings.ToLower(cql.currentToken)
		if cql.currentIndex < cql.cqlLen {
			cql.currentIndex++
		} else { //unterminated
			cql.currentTokenType = TT_EOF //notify error
		}
	} else { //4. WORD, including: logic ops, bbox, within, dwithin, timestring
		cql.currentTokenType = TT_WORD
		buf.Reset() //reset buffer
		for cql.currentIndex < cql.cqlLen && !SearchStringInString("()/<>= \t\r\n", string(cql.cqlCharArray[cql.currentIndex])) {
			buf.WriteString(string(cql.cqlCharArray[cql.currentIndex]))
			cql.currentIndex++
		}
		cql.currentToken = buf.String()
		lval = strings.ToLower(cql.currentToken)
		//strip leading and tailing single quotes
		stripVal := strings.Trim(cql.currentToken, "'")

		if lval == "or" || lval == "and" || lval == "not" {
			cql.currentTokenType = TT_LOGIC_OPS
		} else if lval == "bbox" {
			cql.currentTokenType = TT_BBOX
		} else if lval == "within" {
			cql.currentTokenType = TT_WITHIN
		} else if lval == "dwithin" {
			cql.currentTokenType = TT_DWITHIN
		} else if lval == "sortby" {
			cql.currentTokenType = TT_SORTBY
		} else if PatternMatch(DATE_PATTERN, stripVal) || PatternMatch(DATE_TIME_PATTERN, stripVal) {
			cql.currentTokenType = TT_TIMESTRING
		}
	}
}

/**
  *sql:
  *  BBOX(origin_geom,174,-41,175,-42)
  sql:
  ST_Contains(ST_SetSRID(ST_Envelope('LINESTRING(174 -41,175 -42)'::geometry),4326), origin_geom)
  * @param sql
*/
func (cql *CqlConverter) ToBBoxSql(sql *bytes.Buffer, args []interface{}) ([]interface{}, error) {
	token := cql.currentToken //bbox
	//log.Println("token1 " + token)
	cql.NextToken()

	if cql.currentTokenType != TT_PARENTHESIS {
		//log.Fatalln("Error, not valid bbox!!")
		return args, errors.New("Invalid bbox!!")
	}
	tokens := make([]string, 0)

	for ")" != cql.currentToken { //move until the close parenthesis
		cql.NextToken()
		token = cql.currentToken // origin_geom,174,-41,175,-42
		//log.Println("token2 " + token)
		if ")" != token {
			//String[] bboxContents = token3.split("[, ]+");
			for _, value := range regexp.MustCompile("[, ]").Split(token, -1) {
				if value != "" {
					tokens = append(tokens, value)
				}
			}
		}
	}

	if cql.currentTokenType != TT_PARENTHESIS {
		//log.Fatal("Error, not valid bbox!!")
		return args, errors.New("Invalid bbox!!")
	}

	if len(tokens) == 5 {
		for i, v := range tokens {
			if i > 0 { //check numbers, first one is not a number
				if _, e := strconv.ParseFloat(v, 64); e != nil {
					return args, errors.New("Invalid bbox " + token + e.Error())
				}
			}
		}
		bbox := fmt.Sprintf("%s %s,%s %s", tokens[1], tokens[2], tokens[3], tokens[4])
		args = append(args, fmt.Sprintf("LINESTRING(%s)", bbox))
		sql.WriteString(fmt.Sprintf(" ST_Contains(ST_SetSRID(ST_Envelope($%d::geometry),4326),%s)", len(args), tokens[0]))
		//keep the bbox string for later user
		cql.BBOX = tokens[1] + "," + tokens[2] + "," + tokens[3] + "," + tokens[4]
	} else {
		return args, errors.New("Invalid bbox " + token)
	}
	return args, nil
}

/**
* get comma (and white space) separated tokens between parenthesis
* (origin_geom,POLYGON((172.951 -41.767,172.001 -42.832,169.564 -44.341,172.312 -45.412,175.748 -42.908,172.951 -41.767)))
* (origin_geom,Point(175 -41),500,meters)
 */
func (cql *CqlConverter) GetParenthesisedTokens() []string {
	//get all tokens for the WITHIN claus, parse all comma separated tokens
	var token string
	openParenthesis := 1 //already at open parenthesis
	closeParenthesis := 0
	tokens := make([]string, 0)
	for true {
		cql.NextToken()          //this get white space separated
		token = cql.currentToken // origin_geom,174,-41,175,-42
		//log.Println("token2 " + token)

		for _, value := range strings.Split(token, ",") {
			if value != "" {
				tokens = append(tokens, value)
			}
		}

		if "(" == token {
			openParenthesis++
		}
		if ")" == token {
			closeParenthesis++
		}

		//fmt.Printf("openParenthesis %d closeParenthesis %d", openParenthesis, closeParenthesis)
		if closeParenthesis > 0 && openParenthesis == closeParenthesis {
			break
		}
	}
	return tokens
}

/**
  * cql:
  WITHIN(origin_geom,POLYGON((172.951 -41.767,172.001 -42.832,169.564 -44.341,172.312 -45.412,175.748 -42.908,172.951 -41.767)))
  sql:
  ST_Within(origin_geom, ST_GeomFromText('POLYGON((172.951 -41.767,172.001 -42.832,169.564 -44.341,172.312 -45.412,175.748 -42.908,172.951 -41.767))',4326))

  * @param sql
*/

func (cql *CqlConverter) ToWithinSql(sql *bytes.Buffer, args []interface{}) ([]interface{}, error) {
	token := cql.currentToken //within
	//log.Println("token1 " + token)
	cql.NextToken()

	if "(" != cql.currentToken {
		return args, errors.New("Invalid WITHIN format")
	}

	//get all tokens for the WITHIN claus, parse all comma and space separated tokens
	tokens := cql.GetParenthesisedTokens()
	geomFieldName := tokens[0]
	if WFS_DB_GEOM_FIELD != strings.ToLower(geomFieldName) {
		return args, errors.New("Invalid field name " + geomFieldName)
	}

	openParenthesis := 0
	closeParenthesis := 0
	i := 1 //start from 1
	numeric := 0
	coordinatePos := 0
	var coordinates bytes.Buffer

	//get the polygon contents
	for i < len(tokens) {
		token = tokens[i]
		i++
		if "(" == token {
			openParenthesis++
			if openParenthesis == 2 { //corrdinates starts here
				coordinatePos = i + 1
			}

		}
		if ")" == token {
			closeParenthesis++
		}

		if "POLYGON" == strings.ToUpper(token) {
			numeric = 0
		}
		//get the coordinates
		if openParenthesis == 2 && closeParenthesis == 0 {
			//if PatternMatch("-?[.\\d]+", token) { //numeric
			if i >= coordinatePos {
				if PatternMatch(NUMBER_PATTERN, token) { //check numeric
					numeric++
					if numeric > 1 {
						if numeric%2 == 0 {
							coordinates.WriteString(" ")
						} else {
							coordinates.WriteString(",")
						}
					}
					coordinates.WriteString(token)
				} else { //not valid
					return args, errors.New("Invalid number format " + token)
				}
			}
		}
		//fmt.Printf("openParenthesis %d closeParenthesis %d", openParenthesis, closeParenthesis)
		if openParenthesis == closeParenthesis && closeParenthesis > 0 {
			break
		}
	}

	if numeric%2 != 0 {
		return args, errors.New("Invalid WITHIN format, coordinate must be even numbers")
	}

	args = append(args, fmt.Sprintf("POLYGON((%s))", coordinates.String()))
	//use placeholder and args
	sql.WriteString(fmt.Sprintf("ST_Within(%s, ST_GeomFromText($%d, 4326))", geomFieldName, len(args)))

	return args, nil
}

/**
 * cql:
 * DWITHIN(origin_geom,Point(175 -41),500,meters)
 * sql:
 * ST_DWithin(origin_geom::Geography, ST_GeomFromText('POINT(172.951 -41.767)', 4326)::Geography, 5000)
 * @param sql
 */
func (cql *CqlConverter) ToDWithinSql(sql *bytes.Buffer, args []interface{}) ([]interface{}, error) {
	token := cql.currentToken //within
	//log.Println("token1 " + token)
	cql.NextToken()

	if "(" != cql.currentToken {
		//log.Fatalln("Error, not valid DWITHIN!!")
		return args, errors.New("not valid DWITHIN!!")
	}

	//get all tokens for the WITHIN claus, parse all comma and space separated tokens
	tokens := cql.GetParenthesisedTokens()
	//1. get and check geomFieldName
	geomFieldName := tokens[0]
	if WFS_DB_GEOM_FIELD != strings.ToLower(geomFieldName) {
		return args, errors.New("Invalid field name " + geomFieldName)
	}

	openParenthesis := 0
	closeParenthesis := 0
	i := 1
	numeric := 0
	numParenthesis := 0
	geom := ""
	coordinatePos := 0
	var geomString bytes.Buffer //'POINT(172.951 -41.767)
	//2. get the geometry type
	for i < len(tokens) {
		token = tokens[i]
		i++
		if "(" == token {
			openParenthesis++
		}
		if ")" == token {
			closeParenthesis++
		}
		//log.Println("token3 " + token)
		if "POLYGON" == strings.ToUpper(token) || "POLYLINE" == strings.ToUpper(token) || "POINT" == strings.ToUpper(token) {
			geom = strings.ToUpper(token)
			if "POLYGON" == geom || "POLYLINE" == geom {
				geomString.WriteString(geom)
				geomString.WriteString("((")
				numParenthesis = 2
			} else if "POINT" == geom {
				geomString.WriteString("POINT(")
				numParenthesis = 1
			} else {
				return args, errors.New("No valid geometry for function DWITHIN " + geom)
			}
			numeric = 0
			coordinatePos = i + numParenthesis + 1
		}

		//3. get the coordinates
		if openParenthesis == numParenthesis && closeParenthesis == 0 {
			//if PatternMatch(NUMBER_PATTERN, token) { //numeric
			if i >= coordinatePos {
				if PatternMatch(NUMBER_PATTERN, token) { //check numeric
					numeric++
					if numeric > 1 {
						if numeric%2 == 0 {
							geomString.WriteString(" ")
						} else {
							geomString.WriteString(",")
						}
					}
					geomString.WriteString(token)
				} else {
					return args, errors.New("Invalid number format " + token)
				}
			}
		}
		//fmt.Printf("openParenthesis %d closeParenthesis %d", openParenthesis, closeParenthesis)
		if openParenthesis == closeParenthesis && closeParenthesis > 0 { //end of geometry
			break
		}
	}
	if "POLYGON" == geom || "POLYLINE" == geom {
		geomString.WriteString("))")
	} else if "POINT" == geom {
		geomString.WriteString(")")
	}

	//4. get the distance
	if i < len(tokens) {
		token = tokens[i]
		//log.Println("## distance " + token)
		if PatternMatch("^[.\\d]+$", token) {
			distance, err := strconv.ParseFloat(token, 64)
			if err != nil {
				distance = -1
				return args, errors.New("error parsing distance!!")
			}
			i++
			if i < len(tokens) { //get unit
				token = strings.ToLower(tokens[i])
				//fmt.Println("## unit ", token)
				if "meters" == token {
				} else if "feet" == token {
					distance = 0.3048 * distance
				} else {
					return args, errors.New("wrong unit for distance, should be in [meters, feet]!!")
				}
			}
			//geomString: "POINT(172.951 -41.767"
			args = append(args, geomString.String(), distance)
			//ST_DWithin(origin_geom::Geography, ST_GeomFromText('POINT(172.951 -41.767)', 4326)::Geography, 5000)
			sql.WriteString(fmt.Sprintf("ST_DWithin(%s::Geography, ST_GeomFromText($%d, 4326)::Geography, $%d)", geomFieldName, (len(args) - 1), len(args)))
		} else {
			return args, errors.New("Invalid distance " + token)
		}
	} else { //no distance
		return args, errors.New("No distance specified for  DWITHIN")
	}

	return args, nil
}

/**
*   convert cql to sql
 */
func (cql *CqlConverter) ToSQL() (string, []interface{}, error) {
	var (
		sql bytes.Buffer
		err error
	)
	err = nil
	args := []interface{}{}
	//go through the CQL string and convert
	for cql.currentTokenType != TT_EOF {
		cql.NextToken()

		//validate data format
		if cql.previousTokenType1 == TT_RELATION_OPS && (cql.currentTokenType == TT_WORD || cql.currentTokenType == TT_TIMESTRING) {
			//get the field type for the token before the operator e.g. magnitude>=3
			fType := checkDatabaseFieldType(strings.ToLower(cql.previousToken2))
			if fType == WFS_DB_FIELD_TYPE_UNKNOWN {
				return sql.String(), args, errors.New("Unkown parameter name " + cql.previousToken2)
			} else if fType == WFS_DB_FIELD_TYPE_TIME { //check time format
				if cql.currentTokenType != TT_TIMESTRING {
					return sql.String(), args, errors.New("Invalid time format " + cql.currentToken)
				}
			} else if fType == WFS_DB_FIELD_TYPE_NUMBER { //check number format
				if !PatternMatch(NUMBER_PATTERN, cql.currentToken) {
					return sql.String(), args, errors.New("Invalid number format " + cql.currentToken)
				}
			}
		}
		//1. leading space
		if cql.currentTokenType == TT_RELATION_OPS || cql.currentTokenType == TT_LOGIC_OPS {
			sql.WriteString(" ")
		}

		//2. do content
		if cql.currentTokenType == TT_TIMESTRING {
			args = append(args, cql.currentToken)
			sql.WriteString(fmt.Sprintf("$%d::timestamptz", len(args)))
		} else if cql.currentTokenType == TT_RELATION_OPS && cql.currentToken == "==" { //convert to =
			sql.WriteString("=")
		} else if cql.currentTokenType == TT_BBOX {
			args, err = cql.ToBBoxSql(&sql, args)
		} else if cql.currentTokenType == TT_WITHIN {
			args, err = cql.ToWithinSql(&sql, args)
		} else if cql.currentTokenType == TT_DWITHIN {
			args, err = cql.ToDWithinSql(&sql, args)
		} else {
			//placeholder for parameter
			if cql.previousTokenType1 == TT_RELATION_OPS && cql.currentTokenType == TT_WORD {
				args = append(args, cql.currentToken)
				sql.WriteString(fmt.Sprintf("$%d", len(args)))
			} else { //TT_RELATION_OPS, TT_LOGIC_OPS
				sql.WriteString(cql.currentToken)
			}

		}
		//3. end space
		if cql.currentTokenType == TT_RELATION_OPS || cql.currentTokenType == TT_LOGIC_OPS {
			sql.WriteString(" ")
		}

	}

	return sql.String(), args, err

}
