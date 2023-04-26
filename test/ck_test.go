package test

import (
	"testing"
)

func TestCkInsert(t *testing.T) {
	//SnowNode, _ := snowflake.NewNode(1)
	//err, ckclient := ck.NewCkClient("default")
	//if err != nil {
	//	fmt.Println(err)
	//} else {
	//
	//	conn := ckclient.GetConn()
	//	ctx := clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
	//		"max_block_size": 10,
	//	}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
	//		fmt.Println("progress: ", p)
	//	}), clickhouse.WithProfileInfo(func(p *clickhouse.ProfileInfo) {
	//		fmt.Println("profile info: ", p)
	//	}))
	//
	//	if err := conn.Ping(ctx); err != nil {
	//		if exception, ok := err.(*clickhouse.Exception); ok {
	//			fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
	//		}
	//		fmt.Println(err.Error())
	//		return
	//	}
	//
	//	fileds := []string{"id", "advertiser_id", "marketplace_id", "dataset_id", "idempotency_id", "attributed_sales_1d_same_sku",
	//		"attributed_conversions_1d", "attributed_sales_1d", "attributed_conversions_7d", "attributed_conversions_14d_same_sku", "attributed_conversions_14d",
	//		"time_window_start", "attributed_sales_7d", "attributed_conversions_30d", "attributed_units_ordered_14d_same_sku", "attributed_units_ordered_30d", "attributed_sales_7d_same_sku",
	//		"attributed_units_ordered_14d", "attributed_units_ordered_7d", "nattributed_units_ordered_7d_same_sku", "attributed_conversions_30d_same_sku", "ad_group_id",
	//		"placement", "attributed_units_ordered_1d", "attributed_sales_30d", "attributed_units_ordered_30d_same_sku", "attributed_units_ordered_1d_same_sku",
	//		"currency", "ad_id", "attributed_conversions_1d_same_sku", "attributed_sales_14d", "attributed_conversions_7d_same_sku", "attributed_sales_30d_same_sku", "campaign_id",
	//		"keyword_id", "attributed_sales_14d_same_sku", "create_date",
	//	}
	//
	//	filedsstr := strings.Join(fileds, ",")
	//	batch, err := conn.PrepareBatch(ctx, fmt.Sprintf("INSERT INTO sp_conversion (%s)", filedsstr))
	//	if err != nil {
	//		fmt.Println(err.Error())
	//		return
	//	}
	//
	//	if err := batch.Append(SnowNode.Generate().Int64(), fmt.Sprintf("value_%d", i), time.Now()); err != nil {
	//		fmt.Println(err.Error())
	//		return
	//	}
	//
	//	if err := batch.Send(); err != nil {
	//		fmt.Println(err.Error())
	//		return
	//	}
	//}
}

func TestCkSelect(t *testing.T) {

}
