// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

// FeadDao is the manager for logic model data accessing
// and custom defined data operations functions management.
type FeadDao struct {
	gmvc.M                                      // M is the core and embedded struct that inherits all chaining operations from gdb.Model.
	DB      gdb.DB                              // DB is the raw underlying database management object.
	Table   string                              // Table is the table name of the DAO.
	Columns feadColumns // Columns contains all the columns of Table that for convenient usage.
}

// FeadColumns defines and stores column names for table fead.
type feadColumns struct {
	Id                      string //                                       
    ProfileId               string //                                       
    DatasetId               string //                                       
    CreateDate              string //                                       
    Status                  string // 整体是否：1有效 0 无效                
    Step                    string // 授权进行到哪一步，一共3步             
    StepStatus              string // 进行到哪一步的单独的状态1有效 0 无效  
    ErrReason               string // 如果错误，记录失败原因                
    UpdateDate              string //                                       
    ClientTokenId           string // 订阅后的参数，针对队列有用            
    MessagesSubscriptionId  string // 订阅后的参数，针对队列有用            
    Version                 string // 订阅后参数，版本后针对队列排查问题
}

var (
	// Fead is globally public accessible object for table fead operations.
	Fead = FeadDao{
		M:     g.DB("xray_report").Model("fead").Safe(),
		DB:    g.DB("xray_report"),
		Table: "fead",
		Columns: feadColumns{
			Id:                     "id",                        
            ProfileId:              "profile_id",                
            DatasetId:              "dataset_id",                
            CreateDate:             "create_date",               
            Status:                 "status",                    
            Step:                   "step",                      
            StepStatus:             "step_status",               
            ErrReason:              "err_reason",                
            UpdateDate:             "update_date",               
            ClientTokenId:          "client_token_id",           
            MessagesSubscriptionId: "messages_subscription_id",  
            Version:                "version",
		},
	}
)