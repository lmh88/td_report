// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"td_report/app/dao/internal"
)

// campaignDao is the manager for logic model data accessing
// and custom defined data operations functions management. You can define
// methods on it to extend its functionality as you wish.
type campaignDao struct {
	internal.CampaignDao
}

var (
	// Campaign is globally public accessible object for table campaign operations.
	Campaign = campaignDao{
		internal.Campaign,
	}
)

// Fill with you ideas below.