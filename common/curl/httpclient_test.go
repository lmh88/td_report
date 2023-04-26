package curl

import (
	"testing"
)

type MResult struct {
	A int `json:"a"`
	B int `json:"b"`
}

func TestPost(t *testing.T) {
	//client:= NewHttpclient()
	//var url = "http://localhost/index.php?r=site/t"
	//
	//var str []MResult
	//err:=client.SetIsDebug(true).GetResult(&str).HttpPost(url,nil)
	// if err != nil {
	//	 t.Log(err.Error())
	//	 fmt.Println("hahahah")
	//	 t.Failed()
	// }
	//
	// fmt.Println(str)

}

func TestGet(t *testing.T) {
	//client:= NewHttpclient()
	//var url = "http://localhost/index.php?r=site/t"
	//var str []MResult
	//err:=client.SetIsDebug(true).GetResult(&str).HttpGet(url,nil)
	//if err != nil {
	//	t.Log(err.Error())
	//	fmt.Println("hahahah")
	//	t.Failed()
	//} else {
	//	fmt.Println("=========")
	//	fmt.Println(str)
	//}

}
