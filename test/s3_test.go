package test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"td_report/app/bean"
	"td_report/common/amazon/s3"
	"td_report/common/tool"
	"td_report/pkg/logger"
	"td_report/pkg/save_file"
	"testing"
)

func TestS3Upload(t *testing.T) {
	client, err := s3.NewClient(g.Cfg().GetString("s3.bucket"))
	if err != nil {

	}
	var ctx context.Context = context.Background()
	key := "sp/day/campaigns/20220506/20220506_1898361487593447.csv"
	client.UploadSingle(ctx, "D:\\data\\report\\new_dsp_temp\\detail\\20220206\\20220206_20220210_1898361487593447.csv", key)
}

func TestS3Download(t *testing.T) {
	client, err := s3.NewClient(g.Cfg().GetString("s3.bucket"))
	if err != nil {

	}
	key := "dsp/day/detail/20220206/20220206_20220210_1898361487593447.csv"
	client.DownloadSingle(key, "20220206_20220210_1898361487593447.csv")
}

// raw_dev/amazon/sp_api/asin_sales/20220420/20220420_109904696947488_B081LMQZGF_039c5883f60f8a5e7846a9f02fba2256.json (1461 bytes, class STANDARD)
//raw_dev/amazon/sp_api/asin_sales/20220420/20220420_109904696947488_B081LMQZGF_0e60187dbb72402acd5245763351ef6b.json (1439 bytes, class STANDARD)
//raw_dev/amazon/sp_api/asin_sales/20220420/20220420_109904696947488_B081LPL772_496e18b8c186311a52190161f9af26c4.json (1506 bytes, class STANDARD)
//raw_dev/amazon/sp_api/asin_sales/20220420/20220420_109904696947488_B081LPL772_df0241b536a62a0c2040232be2558eaf.json (1445 bytes, class STANDARD)
//raw_dev/amazon/sp_api/asin_sales/20220420/20220420_109904696947488_B0829BYC5K_52f82732afc120dafd731328862fdad3.json (1449 bytes, class STANDARD)
//raw_dev/amazon/sp_api/asin_sales/20220420/20220420_109904696947488_B0829BYC5K_be7983da6c4f6fac74b3c5800e99a79d.json (1454 bytes, class STANDARD)
//raw_dev/amazon/sp_api/asin_sales/20220420/20220420_109904696947488_B095H2HF5B_4a13259ccbc024c3ba51d92bafd619c8.json (1435 bytes, class STANDARD)
//raw_dev/amazon/sp_api/asin_sales/20220420/20220420_109904696947488_B09HGDY9X3_c7d1eb31c35f4d8e3f1ccc0060c09cf0.json (1431 bytes, class STANDARD)
//raw_dev/amazon/sp_api/asin_sales/20220420/20220420_109904696947488_B09LCCB1H9_0b46fa63a44e3165ee5b1e686d00cd8f.json
func TestListNum(t *testing.T) {
	client, err := s3.NewClient(g.Cfg().GetString("s3.bucket"))
	if err != nil {

	}
	dirpath := "sp_api/asin_sales/20220420/20220420_109904696947488_B081LMQZGF_039c5883f60f8a5e7846a9f02fba2256.json"
	num, err := client.GetFileNum(dirpath)
	fmt.Println(num, err)
}

func TestS3Listfile(t *testing.T) {
	client, err := s3.NewClient(g.Cfg().GetString("s3.bucket"))
	if err != nil {

	}

	//dirpath := "raw_dev/amazon/dsp/day/detail/20220206"
	dirpath := "raw_dev/amazon/sp_api/asin_sales/20220420"

	//dirpath := "raw/amazon/sb/month"
	client.ListFilesWithPrefix(dirpath)
}

func TestS3UploadDir(t *testing.T) {
	dir := "D:\\data\\sp_api\\asin_sales\\20220420\\109904696947488"
	client, err := s3.NewClient(g.Cfg().GetString("s3.bucket"))
	if err != nil {

	} else {

		var ctx = logger.Logger.NewTraceIDContext(context.Background(), fmt.Sprintf("%s", dir))
		client.UploadDir(ctx, dir, "json", "sp_api/asin_sales/20220420")
	}
}

func TestDayweek(t *testing.T) {
	date := "2022-05-03"
	year, w := tool.GetWeek(date)
	fmt.Println(year, w)
}

func TestSet(t *testing.T) {
	data := &bean.UploadS3Data{
		Key:  "sp/day/target/20220511/20220511_1391827227662581.gz",
		Path: "D:/data/report/new_sp/targets/20220511/20220511_1391827227662581.gz",
	}

	save_file.Saveerrordata(data)
}
