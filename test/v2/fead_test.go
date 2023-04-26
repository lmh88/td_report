package v2

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"td_report/app/bean"
	"td_report/app/model"
	"td_report/app/repo"
	"td_report/app/service/fead"
	"td_report/cmd/report_job_v2"
	"testing"
)

func TestAddBatch(t *testing.T) {
	list := make([]*model.FeadSubscription, 0)
	list = append(list, &model.FeadSubscription{
		ProfileId: 123,
		DatasetId: "sp-traffic",
		CreateDate: gtime.Now(),
	})
	list = append(list, &model.FeadSubscription{
		ProfileId: 456,
		DatasetId: "sp-conversion",
		CreateDate: gtime.Now(),
	})
	repo.NewFeadSubscriptionRepository().AddBatch(list)
}

func TestGetAll(t *testing.T) {
	res := repo.NewFeadSubscriptionRepository().GetAll()
	//for _, item := range res {
	//
	//}
	fmt.Println(res)
}

func getProfileToken() *bean.ProfileToken {
	one := bean.ProfileToken{
		ProfileId: "917533703417640",
		Region: "NA",
		RefreshToken: "Atzr|IwEBIFKdXhERp0ljWkL1MUS14QhDhouEUqMSxtZ5dH1Gsb57kfhFgirybmi2CsYjyNNMe9AWBg0NRa58s5RFyotdtY1itEeKtZ5LB7BAYU66ckhtY4c0Pi2gYzEDSffMrJKzBftKyrveIEoW8ZYfYjkGSVbqz97NeBKSwDOZhyzcTBUhGikjd0fW5Cbdi_wLg_qAE0t9mBkUmSYJB0hQSyiHpSkVIiVphgq9GUt7zw8F2E5EAdtakm1wgjC0gS22j5uo8xMsZW7mOOgU10w2GrOtZaUlNDbjlGYlpI1N7e43WNP5Q3rMMdQN0LZfT4ELa7LRUn9JU-EKJ5rTGe2tzmxnQbl1X7hGt-5wIMKy4mUevrtkQJJOsIgZPj8m5DwWo8b_92_kGACQ_GJcjR8U0jTYPCsCWAuG4rgWocoW7vb2y4df1fCfd0dEucT-nF_DtIKm33S4H3uvfJlDCc2qwTdVaKcp",
		ClientId: "amzn1.application-oa2-client.be4e30f5a0b14f488677728ec04c12e0",
		ClientSecret: "c8c0f06cd18861b9cd5b4cd228bf475ca08e7f6f7af4dad384e0772296a3dd24",
	}
	return &one
}


func getByProfile(profileId string) *bean.ProfileToken {
	//dao.Profile.Fields("").Where("").LeftJoin()
	//dao.SellerProfile.Fields("")
	profiles := make([]string, 0)
	profiles = append(profiles, profileId)
	//data, err := repo.NewSellerProfileRepository().GetProfileAndRefreshTokenByFilter(profiles)
	one, _ := repo.NewSellerProfileRepository().GetProfileAndRefreshTokenById(profiles)
	fmt.Println(one[0])
	if len(one) > 0 {
		return one[0]
	}
	return nil
}

func TestVar1(t *testing.T) {
	need := fead.SqsMap
	delete(need, "sp-traffic")
	fmt.Println(need)
	fmt.Println(fead.SqsMap)
	t.Log("over")
}


func TestAddSub(t *testing.T) {
	fead.AddSubscription("557791085879491")
}

func TestDelSub(t *testing.T) {
	//fead.CancelSubscription("")
}

func TestGetSubscriptionByProfile(t *testing.T) {
	profileToken := getByProfile("4136036831441171")
	//profileToken := getProfileToken()
	data, _ := fead.GetSubscriptionByProfile(profileToken)
	fmt.Println(data)
}

func TestSubSomeUser(t *testing.T) {
	arn := "arn:aws:sqs:us-east-1:207485092024:gary-sp-traffic"
	dataId := "sp-traffic"
	profileToken := getByProfile("3551774245623279")
	subId, err := fead.CreateSubscription(profileToken, arn, dataId)
	fmt.Println(subId, err)
}

/*
1854448698434635
3202951000657536
3686874236077302
451026374501123
3551774245623279
 */



func TestFeadConsumer(t *testing.T) {
	fead.Receive("gary-sp-traffic")
	t.Log("over")
}

func TestCmdFeadConsumer(t *testing.T) {
	//fead.Receive("gary-sp-traffic")

	report_job_v2.RootCmd.SetArgs([]string{"fead_consumer", "--fead_queue=gary-sp-traffic"})
	report_job_v2.Execute()
	t.Log("over")
}


func TestCmdFeadSub(t *testing.T) {
	//report_job_v2.RootCmd.SetArgs([]string{"fead_consumer", "--fead_queue=budget-usage"})

	//report_job_v2.RootCmd.SetArgs([]string{"fead_consumer", "--fead_queue=sp-conversion"})

	report_job_v2.RootCmd.SetArgs([]string{"fead_consumer", "--fead_queue=sp-traffic"})

	report_job_v2.Execute()

	t.Log("over")
}


func TestGetClientNum(t *testing.T) {
	profiles, err := repo.NewSellerProfileRepository().GetProfileAndRefreshTokenByFilter([]string{})
	fmt.Println(len(profiles), err)
	var (
		c1 int
		c2 int
	)
	for _, item := range profiles {
		if item.Region == "NA" {
			if item.Tag == 1 {
				c1++
			} else if item.Tag == 2 {
				c2++
			}
		}
	}
	fmt.Printf("c1:%d, c2:%d", c1, c2)
}

