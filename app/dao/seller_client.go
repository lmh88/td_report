// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"td_report/app/dao/internal"
)

// sellerClientDao is the manager for logic model data accessing
// and custom defined data operations functions management. You can define
// methods on it to extend its functionality as you wish.
type sellerClientDao struct {
	internal.SellerClientDao
}

var (
	// SellerClient is globally public accessible object for table seller_client operations.
	SellerClient = sellerClientDao{
		internal.SellerClient,
	}
)

// Fill with you ideas below.