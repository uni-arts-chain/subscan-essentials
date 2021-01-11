package model

import (
	"github.com/shopspring/decimal"
)

type Auction struct {
	Id           int             `json:"-"`
	AuctionId    int             `json:"auction_id"`
	EventIndex   string          `sql:"default: null;size:100;" json:"event_index"`
	BlockNum     int             `gorm:"column:block_num; default:null" json:"block_num" `
	ExtrinsicIdx int             `json:"extrinsic_idx"`
	CollectionId int             `gorm:"column:collection_id; default:null" json:"collection_id"`
	ItemId       int             `gorm:"column:item_id; default:null" json:"item_id"`
	Value        decimal.Decimal `gorm:"column:value; type:decimal(32,0);default:null" json:"value"`
	Owner        string          `gorm:"column:owner;default:null" json:"owner"`
	StartPrice   decimal.Decimal `gorm:"column:start_price; type:decimal(32,0);default:null" json:"start_price"`
	CurrentPrice decimal.Decimal `gorm:"column:current_price; type:decimal(32,0);default:null" json:"current_price"`
	Increment    decimal.Decimal `gorm:"column:increment; type:decimal(32,0);default:null" json:"increment"`
	BidPrice     decimal.Decimal `gorm:"column:bid_price; type:decimal(32,0);default:null" json:"bid_price"`
	StartTime    int             `gorm:"column:start_time; default:null" json:"start_time"`
	EndTime      int             `gorm:"column:end_time; default:null" json:"end_time"`
	EventIdx     int             `json:"event_idx"`
	Sender       string          `json:"sender"`
	Status       string          `gorm:"type:varchar(50)" json:"status"`
}

func (c Auction) TableName() string {
	return "auctions"
}
