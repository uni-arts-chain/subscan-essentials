package dao

import (
	"encoding/json"
	"fmt"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/nft/model"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/rpc"
	rpcStorage "github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/shopspring/decimal"
)

const (
	wsNftAuctionCreate = 10
)

func NewAuction(db storage.DB, paramEvent []storage.EventParam, event *storage.Event, block *storage.Block) error {
	auctionId := int(paramEvent[0].Value.(float64))
	collectionId := int(paramEvent[1].Value.(float64))
	itemId := int(paramEvent[2].Value.(float64))
	value := decimal.RequireFromString(util.ToString(paramEvent[3].Value))
	startPrice := decimal.RequireFromString(util.ToString(paramEvent[4].Value))
	sender := util.ToString(paramEvent[5].Value)
	key := storageKey.EncodeStorageKey("Nft", "AuctionList", util.IntToEncode64Hex(collectionId), util.IntToEncode64Hex(itemId))
	v := &rpc.JsonRpcResult{}
	if err := websocket.SendWsRequest(nil, v, rpc.StateGetStorage(wsNftAuctionCreate, util.AddHex(key.EncodeKey), block.Hash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			log.Info("get dataHex failure, error=[%v]\n, block_num=[%v]\n, EventIdx=[%v]", err, block.BlockNum, event.EventIdx)
		} else {
			dataD, err := rpcStorage.Decode(dataHex, key.ScaleType, nil)
			if err != nil {
				log.Info("get Decode failure, error=[%v]\n, block_num=[%v]\n, EventIdx=[%v]", err, block.BlockNum, event.EventIdx)
			} else {
				auction := model.Auction{
					AuctionId:    auctionId,
					EventIndex:   fmt.Sprintf("%d-%d", block.BlockNum, event.ExtrinsicIdx),
					EventIdx:     event.EventIdx,
					BlockNum:     block.BlockNum,
					ExtrinsicIdx: event.ExtrinsicIdx,
					CollectionId: collectionId,
					ItemId:       itemId,
					Value:        value,
					StartPrice:   startPrice,
					Sender:       sender,
					Status:       "create",
				}
				json.Unmarshal([]byte(dataD), &auction)

				query := db.Create(&auction)
				if query != nil {
					log.Info("=====create nft auction failed========" + query.Error())
				}
				return query
			}
		}
	}
	return nil
}

func CancelAuction(db storage.DB, paramEvent []storage.EventParam, event *storage.Event, block *storage.Block) error {
	auctionId := int(paramEvent[0].Value.(float64))
	collectionId := int(paramEvent[1].Value.(float64))
	itemId := int(paramEvent[2].Value.(float64))
	auction := model.Auction{
		AuctionId:    auctionId,
		EventIndex:  fmt.Sprintf("%d-%d", block.BlockNum, event.ExtrinsicIdx),
		EventIdx:     event.EventIdx,
		BlockNum:     block.BlockNum,
		ExtrinsicIdx: event.ExtrinsicIdx,
		CollectionId: collectionId,
		ItemId:       itemId,
		Status:       "cancel",
	}
	query := db.Create(&auction)
	if query != nil {
		log.Info("=====cancel nft auction failed========" + query.Error())
	}
	return query
}

func FinishAuction(db storage.DB, paramEvent []storage.EventParam, event *storage.Event, block *storage.Block) error {
	auctionId := int(paramEvent[0].Value.(float64))
	collectionId := int(paramEvent[1].Value.(float64))
	itemId := int(paramEvent[2].Value.(float64))
	value := decimal.RequireFromString(util.ToString(paramEvent[3].Value))
	bidPrice := decimal.RequireFromString(util.ToString(paramEvent[4].Value))
	owner := util.ToString(paramEvent[5].Value)
	auction := model.Auction{
		AuctionId:    auctionId,
		EventIndex:  fmt.Sprintf("%d-%d", block.BlockNum, event.ExtrinsicIdx),
		EventIdx:     event.EventIdx,
		BlockNum:     block.BlockNum,
		ExtrinsicIdx: event.ExtrinsicIdx,
		CollectionId: collectionId,
		ItemId:       itemId,
		Value:        value,
		Owner:        owner,
		BidPrice:     bidPrice,
		Status:       "finish",
	}
	query := db.Create(&auction)
	if query != nil {
		log.Info("=====trade nft auction failed========" + query.Error())
	}
	return query
}

func BidAuction(db storage.DB, paramEvent []storage.EventParam, event *storage.Event, block *storage.Block) error {
	auctionId := int(paramEvent[0].Value.(float64))
	collectionId := int(paramEvent[1].Value.(float64))
	itemId := int(paramEvent[2].Value.(float64))
	value := decimal.RequireFromString(util.ToString(paramEvent[3].Value))
	bidPrice := decimal.RequireFromString(util.ToString(paramEvent[4].Value))
	sender := util.ToString(paramEvent[5].Value)
	auction := model.Auction{
		AuctionId:    auctionId,
		EventIndex:  fmt.Sprintf("%d-%d", block.BlockNum, event.ExtrinsicIdx),
		EventIdx:     event.EventIdx,
		BlockNum:     block.BlockNum,
		ExtrinsicIdx: event.ExtrinsicIdx,
		CollectionId: collectionId,
		ItemId:       itemId,
		Value:        value,
		BidPrice:     bidPrice,
		Sender:       sender,
		Status:       "bid",
	}
	query := db.Create(&auction)
	if query != nil {
		log.Info("=====bid nft auction failed========" + query.Error())
	}
	return query
}
