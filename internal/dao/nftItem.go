package dao

import (
	//"encoding/base64"
	//"encoding/json"
	//"fmt"

	//"encoding/json"

	"strconv"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	//"github.com/itering/substrate-api-rpc/rpc"
	//"github.com/itering/substrate-api-rpc/storage"
	//"github.com/itering/substrate-api-rpc/storageKey"
	//"github.com/itering/substrate-api-rpc/websocket"
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
	//v := &rpc.JsonRpcResult{}
	cid := params[0]["value"]
	iid := params[1]["value"]
	collectionId,_ := strconv.Atoi(util.ToString(cid))
	// todo why it is string
	itemId,_ := strconv.Atoi(util.ToString(iid))
	//key := storageKey.EncodeStorageKey("Nft", "NftItemList", util.IntToEncode64Hex(cid.(int)),  util.IntToEncode64Hex(iid.(int)))
	//log.Info("=====Start NFT STORAGE====0====" + ":" + key.EncodeKey)
	//if err := websocket.SendWsRequest(nil, v, rpc.StateGetStorage(wsNftCreate, key.EncodeKey, blockHash)); err != nil {
	//	return fmt.Errorf("websocket send error: %v", err)
	//}
	//log.Info("=====END NFT STORAGE========")
	properties := ""
	name := ""
	author := ""
	//if dataHex, err := v.ToString(); err == nil {
	//	if dataHex == "" {
	//		log.Error("base64 decode failure, error=[%v]\n", err)
	//	} else {
	//		dataD, err := storage.Decode(dataHex, key.ScaleType, nil)
	//		if err != nil {
	//			log.Info("=====Decode hex========")
	//		}
	//		result := dataD.ToMapString()
	//		propertiesD, _ := result["data"]
	//		bproperties, err := base64.StdEncoding.DecodeString(propertiesD)
	//		if err != nil {
	//			author = result["owner"]
	//		} else {
	//			properties = string(bproperties)
	//			infos := map[string]interface{}{}
	//			json.Unmarshal(bproperties, &infos)
	//			name = infos["name"].(string)
	//			author = infos["owner"].(string)
	//		}
	//	}
	//	log.Info("=====dataHex NFT STORAGE========")
	//	log.Info(dataHex)
	//}

	//log.Info("=====dataHex END STORAGE========")
	log.Info("=====start create nft item=======")
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
