package dao

import (
	"fmt"

	"github.com/go-kratos/kratos/pkg/log"
	"github.com/itering/scale.go/types"
	"github.com/itering/subscan/model"
	"github.com/itering/subscan/util"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/storage"
	"github.com/itering/substrate-api-rpc/storageKey"
	"github.com/itering/substrate-api-rpc/websocket"
)

const (
	wsNamesCreate = 7
)

func (d *Dao) CreateNames(txn *GormDB, ce *model.ChainEvent, blockHash string) error {
	log.Info(ce.EventId)
	if ce.EventId == "NameRegistered" {
		return d.CreateName(txn, ce, blockHash)
	}

	if ce.EventId == "NameUpdated" {
		return d.UpdateName(txn, ce, blockHash)
	}

	return nil
}

func (d *Dao) CreateName(txn *GormDB, ce *model.ChainEvent, blockHash string) error {
	params := []map[string]interface{}{}
	util.UnmarshalAny(&params, ce.Params)
	name := util.StringFromInterfaceSlice(params[0]["value"])
	name_bytes := []byte(name)
	tcompact32 := types.CompactU32{}
	data := tcompact32.Encode(len(name_bytes)).Data
	data = append(data, []byte(name)...)
	nameHex := util.BytesToHex(data)
	key := storageKey.EncodeStorageKey("Names", "Names", util.ToString(nameHex))
	v := &rpc.JsonRpcResult{}
	if err := websocket.SendWsRequest(nil, v, rpc.StateGetStorage(wsNamesCreate, util.AddHex(util.ToString(key.EncodeKey)), blockHash)); err != nil {
		return fmt.Errorf("websocket send error: %v", err)
	}
	value := ""
	owner := ""
	expiration := 0
	if dataHex, err := v.ToString(); err == nil {
		if dataHex == "" {
			log.Info("get dataHex failure, error=[%v]\n, block_num=[%v]\n, EventIdx=[%v]", err, ce.BlockNum, ce.EventIdx)
		} else {
			dataD, err := storage.Decode(dataHex, key.ScaleType, nil)
			if err != nil {
				log.Info("get Decode failure, error=[%v]\n, block_num=[%v]\n, EventIdx=[%v]", err, ce.BlockNum, ce.EventIdx)
			} else {
				result := dataD.ToMapInterface()
				value = util.StringFromInterfaceSlice(result["value"])
				owner, _ = result["owner"].(string)
				expiration = util.IntFromInterface(result["expiration"])
			}
		}
	}

	e := model.Name{
		EventIndex:   ce.EventIndex,
		BlockNum:     ce.BlockNum,
		ExtrinsicIdx: ce.ExtrinsicIdx,
		Name:         name,
		Value:        value,
		Owner:        owner,
		EventIdx:     ce.EventIdx,
		Expiration:   expiration,
	}
	query := txn.Create(&e)
	if query.RowsAffected == 0 {
		log.Info("=====create names failed========" + query.Error.Error())
	}
	if query.Error != nil {
		log.Info("=====create names failed========" + query.Error.Error())
	}

	return d.checkDBError(query.Error)
}

func (d *Dao) UpdateName(txn *GormDB, ce *model.ChainEvent, blockHash string) error {
	params := []map[string]interface{}{}
	util.UnmarshalAny(&params, ce.Params)
	name := util.StringFromInterfaceSlice(params[0]["value"])
	nameData := params[1]["value"].(map[string]interface{})
	owner := util.ToString(nameData["owner"])
	expiration := util.IntFromInterface(nameData["expiration"])
	value := util.StringFromInterfaceSlice(nameData["value"])

	var mname model.Name
	if txn.Where("name = ?", name).First(&mname).RecordNotFound() {
		e := model.Name{
			EventIndex:   ce.EventIndex,
			BlockNum:     ce.BlockNum,
			ExtrinsicIdx: ce.ExtrinsicIdx,
			Name:         name,
			Value:        value,
			Owner:        owner,
			EventIdx:     ce.EventIdx,
			Expiration:   expiration,
		}
		query := txn.Create(&e)
		if query.RowsAffected == 0 {
			log.Info("=====update names failed========" + query.Error.Error())
		}
		if query.Error != nil {
			log.Info("=====update names failed========" + query.Error.Error())
		}

		return d.checkDBError(query.Error)
	} else {
		if ce.BlockNum > mname.BlockNum || (ce.BlockNum == mname.BlockNum && ce.EventIdx > mname.EventIdx) {
			query := txn.Debug().Model(&mname).Updates(map[string]interface{}{"event_index": ce.EventIndex, "block_num": ce.BlockNum, "extrinsic_idx": ce.ExtrinsicIdx, "value": value, "owner": owner, "expiration": expiration, "event_idx": ce.EventIdx})
			if query.RowsAffected == 0 {
				log.Info("=====update names failed========" + query.Error.Error())
			}
			if query.Error != nil {
				log.Info("=====update names failed========" + query.Error.Error())
			}
			return d.checkDBError(query.Error)
		}
	}
	return nil
}
