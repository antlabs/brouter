package brouter

import "testing"

func Test_demo_1(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/",
			lookupPath: "/",
			paramKey:   []string{""},
			paramValue: []string{""},
		},
		{
			insertPath: "/hello/:name",
			lookupPath: "/hello/guonaihong",
			paramKey:   []string{"name"},
			paramValue: []string{"guonaihong"},
		},
	}

	tc.run(t)
}

func Test_demo_2(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/authorizations",
			lookupPath: "/authorizations",
			paramKey:   []string{""},
			paramValue: []string{""},
		},
	}

	tc.run(t)
}

func Test_demo_3(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/authorizations",
			lookupPath: "/authorizations",
			paramKey:   []string{""},
			paramValue: []string{""},
		},
		{
			insertPath: "/authorizations/:id",
			lookupPath: "/authorizations/hello",
			paramKey:   []string{"id"},
			paramValue: []string{"hello"},
		},
	}

	tc.run(t)
}
