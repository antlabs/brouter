package brouter

import (
	"testing"
)

func Test_github_Param1(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/teams/:id/repos",
			lookupPath: "/teams/antlabs/repos",
			paramKey:   []string{"id"},
			paramValue: []string{"antlabs"},
		},
		{
			insertPath: "/teams/:id/repos/:owner/:repo",
			lookupPath: "/teams/antlabs-aaa/repos/guonaihong/baserouter-aaa",
			paramKey:   []string{"id", "owner", "repo"},
			paramValue: []string{"antlabs-aaa", "guonaihong", "baserouter-aaa"},
		},
		{
			insertPath: "/repos/:owner/:repo/pulls/:number/files",
			lookupPath: "/repos/guonaihong/baserouter/pulls/1/files",
			paramKey:   []string{"owner", "repo", "number"},
			paramValue: []string{"guonaihong", "baserouter", "1"},
		},
		{
			insertPath: "/repos/:owner/:repo/pulls/:number/merge",
			lookupPath: "/repos/NaihongGuo/deepcopy/pulls/2/merge",
			paramKey:   []string{"owner", "repo", "number"},
			paramValue: []string{"NaihongGuo", "deepcopy", "2"},
		},
		{
			insertPath: "/repos/:owner/:repo/pulls/:number/comments",
			lookupPath: "/repos/guonh/timer/pulls/3/comments",
			paramKey:   []string{"owner", "repo", "number"},
			paramValue: []string{"guonh", "timer", "3"},
		},
	}

	tc.run(t)
}

func Test_github_Param2(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/authorizations/:id",
			lookupPath: "/authorizations/12",
			paramKey:   []string{"id"},
			paramValue: []string{"12"},
		},
		{
			insertPath: "/applications/:client_id/tokens",
			lookupPath: "/applications/client_id-aaa/tokens",
			paramKey:   []string{"client_id"},
			paramValue: []string{"client_id-aaa"},
		},
		{
			insertPath: "/applications/:client_id/tokens/:access_token",
			lookupPath: "/applications/client_id-bbb/tokens/access_token-aaa",
			paramKey:   []string{"client_id", "access_token"},
			paramValue: []string{"client_id-bbb", "access_token-aaa"},
		},
	}

	tc.run(t)
}

func Test_github_Param3(t *testing.T) {

	tc := testCases{
		{
			insertPath: "/teams/:id",
			lookupPath: "/teams/antlabs",
			paramKey:   []string{"id"},
			paramValue: []string{"antlabs"},
		},
		{
			insertPath: "/teams/:id/members/:user",
			lookupPath: "/teams/antlabs/members/guonaihong",
			paramKey:   []string{"id", "user"},
			paramValue: []string{"antlabs", "guonaihong"},
		},
	}

	tc.run(t)
}

// tail里面是长的，insert里面是短的
func Test_github_Param4(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/repos/:owner/:repo/commits/:what/comments",
			lookupPath: "/repos/guonaihong/baserouter/commits/wokao/comments",
			paramKey:   []string{"owner", "repo", "what"},
			paramValue: []string{"guonaihong", "baserouter", "wokao"},
		},
		{
			insertPath: "/repos/:owner/:repo/commits/:what",
			lookupPath: "/repos/guonaihong/baserouter/commits/wokao",
			paramKey:   []string{"owner", "repo"},
			paramValue: []string{"guonaihong", "baserouter"},
		},
	}

	tc.run(t)
}

func Test_github_Param5(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/authorizations",
			lookupPath: "/authorizations",
			paramKey:   []string{""},
			paramValue: []string{""},
		},
		{
			insertPath: "/authorizations/:id",
			lookupPath: "/authorizations/12",
			paramKey:   []string{"id"},
			paramValue: []string{"12"},
		},
		{
			insertPath: "/applications/:client_id/tokens/:access_token",
			lookupPath: "/applications/client_id-bbb/tokens/access_token-aaa",
			paramKey:   []string{"client_id", "access_token"},
			paramValue: []string{"client_id-bbb", "access_token-aaa"},
		},
		{
			insertPath: "/repos/:owner/:repo/events",
			lookupPath: "/repos/guonaihong/baserouter/events",
			paramKey:   []string{"owner", "repo"},
			paramValue: []string{"guonaihong", "baserouter"},
		},
		{
			insertPath: "/orgs/:org/events",
			lookupPath: "/orgs/antlabs/events",
			paramKey:   []string{"org"},
			paramValue: []string{"antlabs"},
		},
	}

	tc.run(t)
}

func Test_github_Param6(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/authorizations",
			lookupPath: "/authorizations",
			paramKey:   []string{""},
			paramValue: []string{""},
		},
		{
			insertPath: "/authorizations/:id",
			lookupPath: "/authorizations/123",
			paramKey:   []string{"id"},
			paramValue: []string{"123"},
		},
		{
			insertPath: "/repos/:owner/:repo/events",
			lookupPath: "/repos/antlabs/baserouter/events",
			paramKey:   []string{"owner", "repo"},
			paramValue: []string{"antlabs", "baserouter"},
		},
	}

	tc.run(t)
}

func Test_github_Param7(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/users/:user/events/public",
			lookupPath: "/users/guonaihong/events/public",
			paramKey:   []string{"user"},
			paramValue: []string{"guonaihong"},
		},
		{
			insertPath: "/users/:user/events/orgs/:org",
			lookupPath: "/users/guonaihong/events/orgs/antlabs",
			paramKey:   []string{"user", "org"},
			paramValue: []string{"guonaihong", "antlabs"},
		},
		{
			insertPath: "/feeds",
			lookupPath: "/feeds",
			paramKey:   []string{""},
			paramValue: []string{""},
		},
		{
			insertPath: "/repos/:owner/:repo/notifications",
			lookupPath: "/repos/guonaihong/baserouter/notifications",
			paramKey:   []string{"owner", "repo"},
			paramValue: []string{"guonaihong", "baserouter"},
		},
		{
			insertPath: "/notifications/threads/:id",
			lookupPath: "/notifications/threads/10",
			paramKey:   []string{"id"},
			paramValue: []string{"10"},
		},
		{
			insertPath: "/notifications/threads/:id/subscription",
			lookupPath: "/notifications/threads/20/subscription",
			paramKey:   []string{"id"},
			paramValue: []string{"20"},
		},
		{
			insertPath: "/repos/:owner/:repo/stargazers",
			lookupPath: "/repos/guonaihong/baserouter/stargazers",
			paramKey:   []string{"owner", "repo"},
			paramValue: []string{"guonaihong", "baserouter"},
		},
		{
			insertPath: "/users/:user/starred",
			lookupPath: "/users/guonaihong/starred",
			paramKey:   []string{"user"},
			paramValue: []string{"guonaihong"},
		},
		{
			insertPath: "/user/starred",
			lookupPath: "/user/starred",
			paramKey:   []string{""},
			paramValue: []string{""},
		},
		{
			insertPath: "/user/starred/:owner/:repo",
			lookupPath: "/user/starred/guonaihong/baserouter",
			paramKey:   []string{"owner"},
			paramValue: []string{"guonaihong"},
		},
		{
			insertPath: "/repos/:owner/:repo/pulls/comments",
			lookupPath: "/repos/guonaihong/baserouter/pulls/comments",
			paramKey:   []string{"owner", "repo"},
			paramValue: []string{"guonaihong", "baserouter"},
		},
	}

	tc.run(t)
}

func Test_github_Param8(t *testing.T) {
	tc := testCases{
		{
			insertPath: "/authorizations",
			lookupPath: "/authorizations",
			paramKey:   []string{""},
			paramValue: []string{""},
		},
		{
			insertPath: "/authorizations/:id",
			lookupPath: "/authorizations/12",
			paramKey:   []string{"id"},
			paramValue: []string{"12"},
		},
		{
			insertPath: "/applications/:client_id/tokens/:access_token",
			lookupPath: "/applications/12/tokens/access_token_haha",
			paramKey:   []string{"client_id", "access_token"},
			paramValue: []string{"12", "access_token_haha"},
		},
		{
			insertPath: "/events",
			lookupPath: "/events",
			paramKey:   []string{},
			paramValue: []string{},
		},
		{
			insertPath: "/repos/:owner/:repo/events",
			lookupPath: "/repos/guonaihong/baserouter/events",
			paramKey:   []string{"owner", "repo"},
			paramValue: []string{"guonaihong", "baserouter"},
		},
		{
			insertPath: "/networks/:owner/:repo/events",
			lookupPath: "/networks/guonaihong/baserouter/events",
			paramKey:   []string{"owner", "repo"},
			paramValue: []string{"guonaihong", "baserouter"},
		},
		{
			insertPath: "/orgs/:org/events",
			lookupPath: "/orgs/antlabs/events",
			paramKey:   []string{"org"},
			paramValue: []string{"antlabs"},
		},
		{
			insertPath: "/users/:user/received_events",
			lookupPath: "/users/guonaihong/received_events",
			paramKey:   []string{"user"},
			paramValue: []string{"guonaihong"},
		},
		{
			insertPath: "/users/:user/received_events/public",
			lookupPath: "/users/guonaihong/received_events/public",
			paramKey:   []string{"user"},
			paramValue: []string{"guonaihong"},
		},
		{
			insertPath: "/users/:user/events/public",
			lookupPath: "/users/guonaihong/events/public",
			paramKey:   []string{"user"},
			paramValue: []string{"guonaihong"},
		},
		{
			insertPath: "/notifications/threads/:id",
			lookupPath: "/notifications/threads/11",
			paramKey:   []string{"id"},
			paramValue: []string{"11"},
		},
	}

	tc.run(t)
}

func Test_github_Param9(t *testing.T) {
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
		{
			insertPath: "/applications/:client_id/tokens/:access_token",
			lookupPath: "/applications/guonaihong/tokens/token",
			paramKey:   []string{"client_id", "access_token"},
			paramValue: []string{"guonaihong", "token"},
		},
	}

	tc.run(t)
}
