package dao

import (
	"strconv"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/shopspring/decimal"
)

func (d *Dao) CreateNftOrder(txn *GormDB, ce *model.ChainEvent) error {

	params := []map[string]interface{}{}
	util.UnmarshalAny(&params, ce.Params)
	cid := params[0]["value"]
	iid := params[1]["value"]
	value := params[2]["value"]
	price := params[3]["value"]
	sender := ""
	if len(params) > 4 {
		sender = util.ToString(params[4]["value"])
	}

	collectionId, _ := strconv.Atoi(util.ToString(cid))
	itemId, _ := strconv.Atoi(util.ToString(iid))
	e := model.NftOrder{
		EventIndex:   ce.EventIndex,
		BlockNum:     ce.BlockNum,
		ExtrinsicIdx: ce.ExtrinsicIdx,
		CollectionId: collectionId,
		ItemId:       itemId,
		Value:        decimal.RequireFromString(util.ToString(value)),
		Price:        decimal.RequireFromString(util.ToString(price)),
		EventIdx:     ce.EventIdx,
		Sender:       sender,
		Status:       "create",
	}

	query := txn.Create(&e)
	if query.Error != nil {
		log.Error("Nft order create failed:" + query.Error.Error())

	}
	return d.checkDBError(query.Error)
}

func (d *Dao) CreateNftOrderCancel(txn *GormDB, ce *model.ChainEvent) error {
	params := []map[string]interface{}{}
	util.UnmarshalAny(&params, ce.Params)
	cid := params[0]["value"]
	iid := params[1]["value"]
	collectionId, _ := strconv.Atoi(util.ToString(cid))
	itemId, _ := strconv.Atoi(util.ToString(iid))
	e := model.NftOrder{
		EventIndex:   ce.EventIndex,
		BlockNum:     ce.BlockNum,
		ExtrinsicIdx: ce.ExtrinsicIdx,
		CollectionId: collectionId,
		ItemId:       itemId,
		EventIdx:     ce.EventIdx,
		Status:       "cancel",
	}
	query := txn.Create(&e)
	return d.checkDBError(query.Error)
}

func (d *Dao) CreateNftOrderSucceed(txn *GormDB, ce *model.ChainEvent) error {
	params := []map[string]interface{}{}
	util.UnmarshalAny(&params, ce.Params)
	cid := params[0]["value"]
	iid := params[1]["value"]
	sender := params[2]["value"]
	collectionId, _ := strconv.Atoi(util.ToString(cid))
	itemId, _ := strconv.Atoi(util.ToString(iid))
	e := model.NftOrder{
		EventIndex:   ce.EventIndex,
		BlockNum:     ce.BlockNum,
		ExtrinsicIdx: ce.ExtrinsicIdx,
		CollectionId: collectionId,
		ItemId:       itemId,
		EventIdx:     ce.EventIdx,
		Sender:       util.ToString(sender),
		Status:       "succeed",
	}
	query := txn.Create(&e)
	return d.checkDBError(query.Error)
}
