package nft

import (
	"github.com/go-kratos/kratos/pkg/log"
	plugin "github.com/itering/subscan-plugin"
	"github.com/itering/subscan-plugin/router"
	"github.com/itering/subscan-plugin/storage"
	"github.com/itering/subscan/plugins/nft/dao"
	"github.com/itering/subscan/plugins/nft/model"
	"github.com/itering/subscan/plugins/nft/service"
	"github.com/itering/subscan/util"
	"github.com/shopspring/decimal"
	"strings"
)

var srv *service.Service

type Nft struct {
	d storage.Dao
}

func New() *Nft {
	return &Nft{}
}

func (a *Nft) InitDao(d storage.Dao) {
	srv = service.New(d)
	a.d = d
	a.Migrate()
}

func (a *Nft) InitHttp() []router.Http {
	return []router.Http{}
}

func (a *Nft) ProcessExtrinsic(*storage.Block, *storage.Extrinsic, []storage.Event) error {
	return nil
}

func (a *Nft) ProcessEvent(block *storage.Block, event *storage.Event, fee decimal.Decimal) error {
	log.Info("plugin ProcessEvent :%v", event.ModuleId)
	if event == nil {
		return nil
	}
	var paramEvent []storage.EventParam
	util.UnmarshalAny(&paramEvent, event.Params)

	switch strings.ToLower(event.EventId) {
	case strings.ToLower("AuctionCreated"):
		return dao.NewAuction(a.d, paramEvent, event, block)
	case strings.ToLower("AuctionCancel"):
		return dao.CancelAuction(a.d, paramEvent, event, block)
	case strings.ToLower("AuctionBid"):
		return dao.BidAuction(a.d, paramEvent, event, block)
	case strings.ToLower("AuctionSucceed"):
		return dao.FinishAuction(a.d, paramEvent, event, block)
	}


	return nil
}

func (a *Nft) SubscribeExtrinsic() []string {
	return nil
}

func (a *Nft) SubscribeEvent() []string {
	return []string{"nft"}
}

func (a *Nft) Version() string {
	return "0.1"
}

func (a *Nft) UiConf() *plugin.UiConfig {
	//conf := new(plugin.UiConfig)
	//	//conf.Init()
	//	////conf.Body.Api.Method = "post"
	//	////conf.Body.Api.Url = "api/plugin/balance/accounts"
	//	////conf.Body.Api.Adaptor = fmt.Sprintf(conf.Body.Api.Adaptor, "list")
	//	////conf.Body.Columns = []plugin.UiColumns{
	//	////	{Name: "address", Label: "address"},
	//	////	{Name: "nonce", Label: "nonce"},
	//	////	{Name: "balance", Label: "balance"},
	//	////	{Name: "lock", Label: "lock"},
	//	////}
	//	//return conf
	return nil
}

func (a *Nft) Migrate() {
	_ = a.d.AutoMigration(&model.Auction{})
	// As, the ws listen to ChainNewHead and ChainFinalizedHead, so the plugin trigger twice.
	// we can add the uniq index to void some action twice
	_ = a.d.AddUniqueIndex(&model.Auction{}, "aid_bn_ei_status", "auction_id", "block_num", "event_idx", "status")
}
