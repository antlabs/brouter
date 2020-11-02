package brouter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testMethodCase struct {
	need   int
	method string
}

func Test_Method(t *testing.T) {
	for _, v := range []testMethodCase{
		{need: 0, method: "GET"},
		{need: 1, method: "POST"},
		{need: 2, method: "PUT"},
		{need: 3, method: "PATCH"},
		{need: 4, method: "DELETE"},
		{need: 5, method: "HEAD"},
	} {
		gotIndex, _ := methodIndex(v.method)
		assert.Equal(t, gotIndex, v.need)
	}
}
