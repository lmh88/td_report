package v2

import (
	"fmt"
	"td_report/app/repo"
	"td_report/app/service/report_v2/get_token"
	"td_report/app/service/report_v2/varible"
	"testing"
)

func TestProfileFilter(t *testing.T) {
	filter := make([]string, 0)
	filter = append(filter, "1480811868525160")
	list, err := repo.NewSellerProfileRepository().GetProfileAndRefreshTokenById(filter)
	//list, err := repo.NewSellerProfileRepository().GetProfileAndRefreshTokenByFilter(filter)
	fmt.Println(list[0].ProfileId, list[0].ClientId, list[0].Tag)
	fmt.Println(len(list), err)
	t.Log("over")
}

func TestRefreshAccess(t *testing.T) {
	crt := varible.ClientRefreshToken{
		ClientId: "amzn1.application-oa2-client.be4e30f5a0b14f488677728ec04c12e0",
		ClientSecret: "c8c0f06cd18861b9cd5b4cd228bf475ca08e7f6f7af4dad384e0772296a3dd24",
		RefreshToken: "abc",
	}
	//r1, r2, r3 := get_token.RefreshToken(&crt)
	//fmt.Println(r1)
	//fmt.Println(r2)
	//fmt.Println(r3)

	t1, e1 := get_token.GetAccessToken(&crt)
	fmt.Println(t1)
	fmt.Println(e1)
}

func TestGetClient(t *testing.T) {
	r1 := repo.NewSellerClientRepository().GetAll()
	fmt.Println(r1)
	t.Log("ov")
}


