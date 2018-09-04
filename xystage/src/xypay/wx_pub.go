package main

import (
	. "JsLib/JsConfig"
	. "JsLib/JsLogger"
	"JsLib/JsOrder"
	"JsLib/JsPay/wxpay"
	"common"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"ider"
	"log"
	"strconv"
	"strings"

	"time"
)

const (
	// appId  = CFG.DirectPay.AppId     // 微信公众平台应用ID
	// mchId  = CFG.DirectPay.MchId     // 微信支付商户平台商户号
	// apiKey = CFG.DirectPay.SecretKey // 微信支付商户平台API密钥

	// 微信支付商户平台证书路径
	certFile   = "cert/apiclient_cert.pem"
	keyFile    = "cert/apiclient_key.pem"
	rootcaFile = "cert/rootca.pem"
)

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var coder = base64.NewEncoding(base64Table)
var c *wxpay.Client = nil

func init_wx_pay() {
	c = wxpay.NewClient(CFG.DirectPay.AppId, CFG.DirectPay.MchId, CFG.DirectPay.SecretKey)

	// 附着商户证书
	err := c.WithCert(certFile, keyFile, rootcaFile)
	if err != nil {
		log.Fatal(err)
	}
}
func wx_pub_pay(order *common.ST_Order) (map[string]string, error) {
	b := make([]byte, 8)
	rand.Read(b)
	nonce_str := strings.ToUpper(coder.EncodeToString(b))
	order.AppId = c.AppId
	order.Mch_id = c.MchId
	order.Nonce_str = nonce_str
	params := make(wxpay.Params)
	// 查询企业付款接口请求参数
	params.SetString("appid", CFG.DirectPay.AppId)
	params.SetString("mch_id", CFG.DirectPay.MchId)
	params.SetString("device_info", "WEB")
	params.SetString("nonce_str", nonce_str) // 随机字符串
	params.SetString("body", order.Desc)
	params.SetString("out_trade_no", order.OrderID) // 商户订单号
	params.SetString("total_fee", strconv.Itoa(order.Amount))
	params.SetString("spbill_create_ip", order.TerminalIp)
	params.SetString("notify_url", CFG.DirectPay.WxPubPayCb)
	params.SetString("trade_type", "JSAPI")
	params.SetString("attach", order.UID)
	params.SetString("openid", order.OpenId)
	params.SetString("sign", c.Sign(params)) // 签名
	Info("Setting params=%v\n", params)
	//url := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	url := CFG.DirectPay.WxPubPayUrl

	ret, err := c.Post(url, params, true)
	if err != nil {
		Error(err.Error())
		return nil, err
	}

	Info("Backing Ret=%v\n", ret)

	charge := make(map[string]string)
	charge["appId"] = CFG.DirectPay.AppId
	charge["timeStamp"] = fmt.Sprintf("%d", time.Now().Unix())
	charge["nonceStr"] = ret["nonce_str"]
	charge["package"] = "prepay_id=" + ret["prepay_id"]
	charge["signType"] = "MD5"

	charge["paySign"] = c.Sign(charge)

	return charge, nil
}

func direct_transfer(transfer *JsOrder.ST_Transfer) (map[string]string, error) {
	b := make([]byte, 8)
	rand.Read(b)
	nonce_str := strings.ToUpper(coder.EncodeToString(b))

	transfer.Tid = strconv.FormatInt(time.Now().Unix(), 10)
	params := make(wxpay.Params)
	// 查询企业付款接口请求参数
	params.SetString("mch_appid", CFG.DirectPay.AppId)
	params.SetString("mchid", CFG.DirectPay.MchId)
	params.SetString("device_info", "WEB")
	params.SetString("nonce_str", nonce_str) // 随机字符串
	params.SetString("partner_trade_no", transfer.Tid)
	params.SetString("openid", transfer.OpenId)
	params.SetString("check_name", "NO_CHECK")
	params.SetString("amount", strconv.Itoa(transfer.Amount))
	params.SetString("desc", transfer.Desc)
	params.SetString("spbill_create_ip", CFG.DirectPay.SpbillCreateIp)

	params.SetString("sign", c.Sign(params)) // 签名

	//url := "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
	url := CFG.DirectPay.WxPubTransferUrl

	ret, err := c.Post(url, params, true)
	if err != nil {
		Error(err.Error())
		return nil, err
	}

	return ret, nil
}

func wx_pub_refund(order *common.ST_Order) (map[string]string, error) {

	b := make([]byte, 8)
	rand.Read(b)
	nonce_str := strings.ToUpper(coder.EncodeToString(b))

	id, err := ider.GenID()
	if err != nil {
		return nil, err
	}
	order.RefundId = id
	params := make(wxpay.Params)

	params.SetString("appid", CFG.DirectPay.AppId)
	params.SetString("mch_id", CFG.DirectPay.MchId)
	params.SetString("device_info", "WEB")
	params.SetString("nonce_str", nonce_str) // 随机字符串
	params.SetString("transaction_id", order.WxPayCb.Transaction_id)
	params.SetString("out_trade_no", order.OrderID) // 商户订单号
	params.SetString("out_refund_no", order.RefundId)
	params.SetString("total_fee", order.WxPayCb.Cash_fee)
	params.SetString("refund_fee", order.WxPayCb.Cash_fee)
	params.SetString("op_user_id", CFG.DirectPay.MchId)

	params.SetString("sign", c.Sign(params)) // 签名

	//url := "https://api.mch.weixin.qq.com/secapi/pay/refund"
	url := CFG.DirectPay.WxPubRefundUrl

	ret, err := c.Post(url, params, true)
	if err != nil {
		Error(err.Error())
		return nil, err
	}

	return ret, nil
}

func wx_bond_refund(order *common.ST_Order) (map[string]string, error) {

	b := make([]byte, 8)
	rand.Read(b)
	nonce_str := strings.ToUpper(coder.EncodeToString(b))

	id, err := ider.GenID()
	if err != nil {
		return nil, err
	}
	order.RefundId = id
	params := make(wxpay.Params)

	params.SetString("appid", CFG.DirectPay.AppId)
	params.SetString("mch_id", CFG.DirectPay.MchId)
	params.SetString("device_info", "WEB")
	params.SetString("nonce_str", nonce_str) // 随机字符串
	params.SetString("transaction_id", order.WxPayCb.Transaction_id)
	params.SetString("out_trade_no", order.OrderID) // 商户订单号
	params.SetString("out_refund_no", order.RefundId)
	params.SetString("total_fee", strconv.Itoa(order.Bond))
	params.SetString("refund_fee", strconv.Itoa(order.Bond))
	params.SetString("op_user_id", CFG.DirectPay.MchId)

	params.SetString("sign", c.Sign(params)) // 签名

	//url := "https://api.mch.weixin.qq.com/secapi/pay/refund"
	url := CFG.DirectPay.WxPubRefundUrl

	ret, err := c.Post(url, params, true)
	if err != nil {
		Error(err.Error())
		return nil, err
	}

	return ret, nil
}
