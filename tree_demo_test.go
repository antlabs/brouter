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

func Test_demo_4(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/authorizations/clients/:client_id",
			lookupPath: "/authorizations/clients/guo_id",
			paramKey:   []string{"client_id"},
			paramValue: []string{"guo_id"},
		},
	}

	tc.run(t)
}

func Test_demo_5(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/applications/:client_id/tokens/:access_token",
			lookupPath: "/applications/id/tokens/guonaihong_token",
			paramKey:   []string{"client_id", "access_token"},
			paramValue: []string{"id", "guonaihong_token"},
		},
	}

	tc.run(t)
}

func Test_demo_6(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/repos/:owner/:repo/events",
			lookupPath: "/repos/guonaihong/brouter/events",
			paramKey:   []string{"owner", "repo"},
			paramValue: []string{"guonaihong", "brouter"},
		},
	}

	tc.run(t)
}
