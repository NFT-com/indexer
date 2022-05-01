package filters

import (
	"fmt"
	"strings"
)

type Where func() string

func Eq(field string, value string) Where {
	return func() string {
		return fmt.Sprintf("%s = %s", field, value)
	}
}

func In(field string, values ...string) Where {
	return func() string {
		value := strings.Join(values, ",")
		value = "(" + value + ")"
		return fmt.Sprintf("%s IN %s", field, value)
	}
}
