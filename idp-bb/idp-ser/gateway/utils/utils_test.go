package utils

import (
	"reflect"
	"testing"

)

func TestParseParams(t *testing.T) {
	type mocktest struct {
		url string
	}

	cases := []struct {
		testCase string
		want map[string]string
		testp *mocktest
	}{
		{
			testCase: "test with params",
			want: map[string]string{"test1":"t1","test2":"t2"},
			testp: &mocktest{
				url: "test.org/test?test1=t1&test2=t2",
			},
		},
		{
			testCase: "test without params",
			want: nil,
			testp: &mocktest{
				url: "test.org/test",
			},
		},
	}

	for _, tt := range cases {
		got, _:= ParseParams(tt.testp.url)
		if reflect.DeepEqual(got, tt.want) == false {
			t.Errorf("case: %v expected to return %q, got %q", tt.testCase,tt.want, got)
		}
	}

} 

func TestParseBasicAuthHeader(t *testing.T){
	ttheader := "Basic Y2xpZW50SUQ6Y2xpZW50U2VjcmV0"

	u,p := ParseBasicAuthHeader(ttheader)

	if u != "clientID" && p != "clientSecret"{
		t.Errorf("expected to return %q and %q, got %q and %q", u, p, u, p)
	}
	
}