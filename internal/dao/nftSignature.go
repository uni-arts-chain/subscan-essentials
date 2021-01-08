package dao

import (
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

type nftSignautre struct {
	Names      []byte `json:"names"`
	NamesOwner string `json:"names_owner"`
	SignTime   int    `json:"sign_time"`
	Memo       string `json:"memo"`
	Collection int    `json:"collection"`
	Item       int    `json:"item"`
	Expiration int    `json:"expiration"`
}

func (d *Dao) CreateNftSignature(txn *GormDB, ce *model.ChainEvent, blockHash string) error {
	params := []map[string]interface{}{}
	util.UnmarshalAny(&params, ce.Params)
	cid := params[0]["value"]
	iid := params[1]["value"]
	collectionId, _ := strconv.Atoi(util.ToString(cid))
	itemId, _ := strconv.Atoi(util.ToString(iid))
	key := storageKey.EncodeStorageKey("Nft", "SignatureList", util.IntToEncode64Hex(collectionId), util.IntToEncode64Hex(itemId))
	v := &rpc.JsonRpcResult{}
	if err := websocket.SendWsRequest(nil, v, rpc.StateGetStorage(wsNftCreate, util.AddHex(key.EncodeKey), blockHash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}

	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			log.Info("get dataHex failure, error=[%v]\n, block_num=[%v]\n, EventIdx=[%v]", err, ce.BlockNum, ce.EventIdx)
		} else {
			dataD, err := storage.Decode(dataHex, "Vec<SignatureAuthentication>", nil)
			if err != nil {
				log.Info("get Decode failure, dataHex=[%v]\n, ScaleType=[%v]\n", dataHex, key.ScaleType)
				log.Info("get Decode failure, error=[%v]\n, block_num=[%v]\n, EventIdx=[%v]", err, ce.BlockNum, ce.EventIdx)
			} else {
				infos := []nftSignautre{}
				json.Unmarshal([]byte(dataD), &infos)

				for _, info := range infos {
					names := util.ToString(info.Names)
					if ce.BlockNum == info.SignTime {
						e := model.NftSignature{
							EventIndex:   ce.EventIndex,
							BlockNum:     ce.BlockNum,
							ExtrinsicIdx: ce.ExtrinsicIdx,
							Names:        names,
							NamesOwner:   info.NamesOwner,
							CollectionId: collectionId,
							ItemId:       itemId,
							EventIdx:     ce.EventIdx,
							SignTime:     info.SignTime,
							Memo:         info.Memo,
						}
						if info.Expiration != 0 {
							e.Expiration = info.Expiration
						}
						query := txn.Debug().Create(&e)
						if query.RowsAffected == 0 {
							log.Info("=====create nft Signature failed========" + query.Error.Error())
						}
						if query.Error != nil {
							log.Info("=====create nft Signature failed========" + query.Error.Error())
							return d.checkDBError(query.Error)
						}
					}
				}
			}
		}
	}
	return nil
}
