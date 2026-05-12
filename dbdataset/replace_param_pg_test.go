package dbdataset

import "testing"

func TestReplaceParamPG(t *testing.T) {
	cases := []struct {
		name        string
		sql         string
		param       string
		paramNumber int
		want        string
	}{
		{"end-of-sql", "select :id", ":id", 1, "select $1"},
		{"space", "select :id from t", ":id", 1, "select $1 from t"},
		{"comma", ":a,:b", ":a", 1, "$1,:b"},
		{"paren-open", "f(:id)", ":id", 1, "f($1)"},
		{"paren-close", "(:id)", ":id", 1, "($1)"},
		{"equals", "id=:id", ":id", 1, "id=$1"},
		{"jsonb-cast", ":payload::jsonb", ":payload", 1, "$1::jsonb"},
		{"int-cast", ":n::int", ":n", 1, "$1::int"},
		{"semicolon", ":id;", ":id", 1, "$1;"},
		{"plus", ":n+1", ":n", 1, "$1+1"},
		{"minus", ":n-1", ":n", 1, "$1-1"},
		{"newline", "select :id\nfrom t", ":id", 1, "select $1\nfrom t"},
		{"tab", "select :id\tfrom t", ":id", 1, "select $1\tfrom t"},
		{"prefix-no-match", ":payload_extra", ":payload", 1, ":payload_extra"},
		{"prefix-letter-no-match", ":payloadx", ":payload", 1, ":payloadx"},
		{"prefix-digit-no-match", ":id1", ":id", 1, ":id1"},
		{"twice", ":id+:id", ":id", 1, "$1+$1"},
		{"jsonb-cast-mid-sql", "update t set p = :payload::jsonb where id = :id", ":payload", 1, "update t set p = $1::jsonb where id = :id"},
		{"not-found", "select * from t", ":id", 1, "select * from t"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, _ := replaceParamPG(tc.sql, tc.param, tc.paramNumber)
			if got != tc.want {
				t.Errorf("replaceParamPG(%q, %q, %d)\n  got:  %q\n  want: %q", tc.sql, tc.param, tc.paramNumber, got, tc.want)
			}
		})
	}
}
