package dao

import (
	//"encoding/base64"
	//"encoding/json"
	//"fmt"

	//"encoding/json"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/websocket"
)

const (
	wsNftCreate = 6
)

func (d *Dao) CreateNft(txn *GormDB, ce *model.ChainEvent, blockHash string) error {
	log.Info(ce.EventId)
	if ce.EventId == "ItemCreated" {
		return d.CreateNftItem(txn, ce, blockHash)
	}

	if ce.EventId == "ItemOrderCreated" {
		return d.CreateNftOrder(txn, ce)
	}

	if ce.EventId == "ItemOrderCancel" {
		return d.CreateNftOrderCancel(txn, ce)
	}

	if ce.EventId == "ItemOrderSucceed" {
		return d.CreateNftOrderSucceed(txn, ce)
	}
	return nil
}

func (d *Dao) CreateNftItem(txn *GormDB, ce *model.ChainEvent, blockHash string) error {
	params := []map[string]interface{}{}
	util.UnmarshalAny(&params, ce.Params)
	cid := params[0]["value"]
	iid := params[1]["value"]
	collectionId,_ := strconv.Atoi(util.ToString(cid))
	itemId,_ := strconv.Atoi(util.ToString(iid))
	key := storageKey.EncodeStorageKey("Nft", "NftItemList", util.IntToEncode64Hex(collectionId),  util.IntToEncode64Hex(itemId))
	v := &rpc.JsonRpcResult{}
	if err := websocket.SendWsRequest(nil, v, rpc.StateGetStorage(wsNftCreate,  util.AddHex(key.EncodeKey), blockHash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	properties := ""
	name := ""
	author := ""
	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			log.Info("get dataHex failure, error=[%v]\n, block_num=[%v]\n, EventIdx=[%v]", err, ce.BlockNum, ce.EventIdx)
		} else {
			dataD, err := storage.Decode(dataHex, key.ScaleType, nil)
			if err != nil {
				log.Info("get Decode failure, error=[%v]\n, block_num=[%v]\n, EventIdx=[%v]", err, ce.BlockNum, ce.EventIdx)
			} else {
				result := dataD.ToMapString()
				properties, _ := result["Data"]
				bproperties, err := base64.StdEncoding.DecodeString(properties)
				if err != nil {
					log.Info("get base64 Decode failure, error=[%v]\n, block_num=[%v]\n, EventIdx=[%v]", err, ce.BlockNum, ce.EventIdx)
					infos := map[string]interface{}{}
					err = json.Unmarshal([]byte(properties), &infos)
					if err == nil {
						name = infos["name"].(string)
						author = infos["author"].(string)
					}
				} else {
					infos := map[string]interface{}{}
					err = json.Unmarshal(bproperties, &infos)
					if err != nil {
						log.Info(err.Error())
					}
					name = infos["name"].(string)
					author = infos["author"].(string)
				}
			}
		}
	}

	e := model.NftItem{
		EventIndex:   ce.EventIndex,
		BlockNum:     ce.BlockNum,
		ExtrinsicIdx: ce.ExtrinsicIdx,
		Name:         name,
		CollectionId: collectionId,
		ItemId:      itemId,
		EventIdx:     ce.EventIdx,
		Sender:       "",
		Author:       author,
		Status:       "create",
		Properties:   properties,
	}
	query := txn.Create(&e)
	if query.RowsAffected == 0 {
		log.Info("=====create nft item failed========" + query.Error.Error())
	}
	if query.Error != nil {
		log.Info("=====create nft item failed========" + query.Error.Error())
	}

	return d.checkDBError(query.Error)
}
