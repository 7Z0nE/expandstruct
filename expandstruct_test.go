package expandstruct

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TT struct {
	A string
	B int
	C float32
}
type T struct {
	I int
	J float32
	T TT
}

func Test_fieldByPath(t *testing.T) {
	s := T{0, 0.0, TT{"hallo", 0, 0.0}}
	sVal := reflect.Indirect(reflect.ValueOf(&s))
	fVal, err := fieldByPath(sVal, "T.A")

	assert.Nil(t, err)
	assert.Equal(t, s.T.A, fVal.Interface())
}

func Test_fieldByPathLowercase(t *testing.T) {
	s := T{0, 0.0, TT{"hallo", 0, 0.0}}
	sVal := reflect.Indirect(reflect.ValueOf(&s))
	fVal, err := fieldByPath(sVal, "t.a")

	assert.Nil(t, err)
	assert.Equal(t, s.T.A, fVal.Interface())
}

func Test_ExpandToStruct(t *testing.T) {
	var m T
	mm := map[string]interface{}{"I": 1, "J": 2.0, "T.A": "42", "T.B": 2, "T.C": 0.5}

	err := ExpandToStruct(mm, &m)

	assert.Nil(t, err)
	assert.Equal(t, m, T{1, 2.0, TT{"42", 2, 0.5}})
}
