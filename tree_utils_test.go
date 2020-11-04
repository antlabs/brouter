package brouter

import (
	"fmt"
	"net/http"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	insertPath string
	lookupPath string
	paramKey   []string
	paramValue []string
}

type testCases []testCase

func (tcs *testCases) run(t *testing.T) {
	d := newTree()
	done := 0

	for _, tc := range *tcs {
		d.insert(tc.insertPath, func(w http.ResponseWriter, r *http.Request, p Params) {
			done++
		})
		//d.debug(80, tc.insertPath, 0, 0, 0)
	}

	for k, tc := range *tcs {

		p2 := make(Params, 0, d.maxParam)
		h := d.lookup(tc.lookupPath, &p2)

		p := p2
		cb := func() {
			handleToUint := *(*uintptr)(unsafe.Pointer(&h))
			assert.NotEqual(t, handleToUint, uintptr(0), fmt.Sprintf("lookup word(%s)", tc.lookupPath))

			fmt.Printf("testCases.run return h address:%x\n", handleToUint)
			h(nil, nil, nil)
			b := assert.Equal(t, done, k+1)
			if !b {
				panic(fmt.Sprintf("done(%d) != %d", done, k+1))
			}

			for index, needKey := range tc.paramKey {
				if len(needKey) == 0 {
					fmt.Printf("index = %d, needKey = 0\n", k)
					continue
				}

				needVal := tc.paramValue
				b := assert.Equal(t, p[index].Key, needKey, fmt.Sprintf("lookup key(%s)", needKey))
				if !b {
					return
				}

				b = assert.Equal(t, p[index].Value, needVal[index], fmt.Sprintf("lookup key(%s)", needKey))
				if !b {
					return
				}
			}
		}

		b := assert.NotPanics(t, cb, fmt.Sprintf("lookup path is(%s)", tc.lookupPath))
		if !b {
			break
		}
	}
}
