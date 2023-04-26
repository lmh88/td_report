// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"td_report/app/dao/internal"
)

// feadSubscriptionDao is the manager for logic model data accessing and custom defined data operations functions management. 
// You can define custom methods on it to extend its functionality as you wish.
type feadSubscriptionDao struct {
	*internal.FeadSubscriptionDao
}

var (
	// FeadSubscription is globally public accessible object for table fead_subscription operations.
	FeadSubscription = feadSubscriptionDao{
		internal.NewFeadSubscriptionDao(),
	}
)

// Fill with you ideas below.