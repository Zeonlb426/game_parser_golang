package types

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"regexp"
	"strconv"
	"strings"
)

var (
	expression = regexp.MustCompile(`-?[\d.]+`)

	errInvalidPattern = errors.New("invalid pattern")
)

type Point struct {
	Latitude  float64
	Longitude float64
}

func (p Point) Value() (driver.Value, error) {
	if 0 == p.Latitude && 0 == p.Longitude {
		return nil, nil
	}

	return fmt.Sprintf("(%f, %f)", p.Latitude, p.Longitude), nil
}

func (p *Point) Scan(value interface{}) error {
	split := strings.Split(value.(string), ",")

	if 2 != len(split) {
		return nil
	}

	var err error

	if p.Latitude, err = p.convertStringToFloat(split[0]); nil != err {
		return err
	}

	if p.Longitude, err = p.convertStringToFloat(split[1]); nil != err {
		return err
	}

	return nil
}

func (p Point) convertStringToFloat(value string) (float64, error) {
	matches := expression.FindStringSubmatch(value)

	if 1 != len(matches) {
		return 0, errInvalidPattern
	}

	return strconv.ParseFloat(matches[0], 64)
}

func (p Point) GormDataType() string {
	return "point"
}

func (p Point) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "postgres":
		fallthrough
	case "mysql":
		return "point"
	case "sqlite":
		return "string"
	}

	return ""
}
