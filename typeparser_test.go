package typeparser_test

import (
	"testing"

	typeparser "github.com/fsamin/go-typeparser"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	typeparser.Verbose = true
	typeparser.Logger = t.Logf

	types, err := typeparser.Parse("fixtures/types.go")
	assert.NoError(t, err, "unable to parse file")

	t1 := types[0]
	assert.NotEmpty(t, t1.Name())
	assert.Equal(t, "TypeA", t1.Name())
	assert.True(t, t1.IsConcrete())

	assert.NotEmpty(t, t1.Docs())
	assert.True(t, t1.Docs().Has("//metadata"))

	assert.NotEmpty(t, t1.FieldNames())
	assert.NotEmpty(t, t1.Fields())

	assert.EqualValues(t, []string{"FieldA", "FieldB"}, t1.FieldNames())

	assert.NotEmpty(t, t1.Field("FieldA").Tags())
	assert.EqualValues(t, []string{"tag:\"value,option\"", "json:\"-\""}, t1.Field("FieldA").Tags())
	assert.EqualValues(t, []string{"value", "option"}, t1.Field("FieldA").TagValue("tag"))

	assert.Len(t, types, 2)

	t2 := types[1]
	assert.NotEmpty(t, t2.Name())
	assert.Equal(t, "InterfaceA", t2.Name())
	assert.True(t, t2.IsInterface())

	assert.EqualValues(t, []string{"foo", "bar"}, t2.MethodNames())

	assert.NotNil(t, t2.Method("foo"))
	assert.NotNil(t, t2.Method("foo").ParamNames())

	assert.NotEmpty(t, t2.Method("foo").Params())
	assert.Len(t, t2.Method("foo").Params(), 1)

	assert.Equal(t, "string", t2.Method("foo").Params()[0].Type())
	assert.EqualValues(t, []string{"string"}, t2.Method("foo").ParamTypes())

	assert.Equal(t, "...int", t2.Method("bar").Params()[0].Type())

	assert.Equal(t, "string", t2.Method("foo").Results()[0].Type())
	assert.Equal(t, "error", t2.Method("foo").Results()[1].Type())
	assert.EqualValues(t, []string{"string", "error"}, t2.Method("foo").ResultTypes())

}
