// ==========================================================================
// This is auto-generated by gf cli tool. DO NOT EDIT THIS FILE MANUALLY.
// ==========================================================================

package internal

import (
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gmvc"
)

// CampaignRuleDao is the manager for logic model data accessing
// and custom defined data operations functions management.
type CampaignRuleDao struct {
	gmvc.M                                      // M is the core and embedded struct that inherits all chaining operations from gdb.Model.
	DB      gdb.DB                              // DB is the raw underlying database management object.
	Table   string                              // Table is the table name of the DAO.
	Columns campaignRuleColumns // Columns contains all the columns of Table that for convenient usage.
}

// CampaignRuleColumns defines and stores column names for table campaign_rule.
type campaignRuleColumns struct {
	Id          string //                             
    CampaignId  string //                             
    AdGroupId   string // adGroupId                   
    StartDate   string //                             
    EndDate     string //                             
    RuleName    string //                             
    RuleType    string //                             
    TemplateId  string //                             
    Condition   string //                             
    Status      string // 规则默认开启1，关闭0,归档2  
    Result      string //                             
    Action      string //                             
    CreatedAt   string //                             
    LastUpdate  string //
}

var (
	// CampaignRule is globally public accessible object for table campaign_rule operations.
	CampaignRule = CampaignRuleDao{
		M:     g.DB("td_xplatform").Model("campaign_rule").Safe(),
		DB:    g.DB("td_xplatform"),
		Table: "campaign_rule",
		Columns: campaignRuleColumns{
			Id:         "id",          
            CampaignId: "campaignId",  
            AdGroupId:  "adGroupId",   
            StartDate:  "startDate",   
            EndDate:    "endDate",     
            RuleName:   "ruleName",    
            RuleType:   "ruleType",    
            TemplateId: "templateId",  
            Condition:  "condition",   
            Status:     "status",      
            Result:     "result",      
            Action:     "action",      
            CreatedAt:  "createdAt",   
            LastUpdate: "lastUpdate",
		},
	}
)