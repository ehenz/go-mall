package pay

import (
	"mall-api/userop-web/global"

	"github.com/smartwalle/alipay/v3"
)

func AlipayPagePay(orderId, totalAmount string) string {
	c := global.SrvConfig.Alipay
	client, err := alipay.New(c.Appid, c.PrivateKey, false)
	if err != nil {
		panic(err)
	}
	err = client.LoadAliPayPublicKey(c.AliPublicKey)
	if err != nil {
		panic(err)
	}

	var p = alipay.TradePagePay{}
	p.NotifyURL = c.NotifyUrl
	p.ReturnURL = c.ReturnUrl
	p.Subject = "标题 - 订单支付"
	p.OutTradeNo = orderId
	p.TotalAmount = totalAmount
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		panic(err)
	}

	// 这个 payURL 即是用于支付的 URL，可将输出的内容复制，到浏览器中访问该 URL 即可打开支付页面。
	var payURL = url.String()

	return payURL
}
