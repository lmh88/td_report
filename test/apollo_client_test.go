package test

import (
	"fmt"
	"td_report/common/config_client/apollo"
	"testing"
)

func TestApolloClient(t *testing.T) {
	data, err := apollo.Getconfig("test")
	fmt.Println(data, "==============")
	t.Log("============11122")
	if err != nil {
		t.Log(err)
	} else {

		t.Log(data)
	}
}
