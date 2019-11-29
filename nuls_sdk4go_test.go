package nuls_sdk4go

import (
	"container/list"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestCreateTransferTrxInDetail(t *testing.T) {
	fmt.Println("TestCreateTransferTrxInDetail  begin")
	expect_trxhex := "0200bbb5e05d0872656d61726b3232008c01170100016f776e306fd399498d50d83ef0d2bdf71a1922d001000100404b4c000000000000000000000000000000000000000000000000000000000008c413ee9158f97881000117010001b6acd7d3a620475fd309aa7b5fc49056d10adacc01000100a0c44a00000000000000000000000000000000000000000000000000000000000000000000000000"
	var tx_time int64 = 1575007675
	prevtxid := "305b6ad9c553775acbb82eb354fb9f3d4201f8ef143eafcdc413ee9158f97881"
	nonce := prevtxid[(len(prevtxid) - 16):]

	s1 := "2f40c87367d648874bd369442cdef92ef8a933b1560759629364fa6278063ade"
	s2 := "90a2030662e196f990300e9f5eb3d403f7245aae9db66b0a461072832429f0b5"
	pri, err := BuildPrivkey(s1)
	//pri, err := nuls_sdk4go.GeneratePriKey()
	pubbyte := GetEncodePubkey(&(pri.PublicKey), true)
	fromaddress := GetNulsAddressByPub(NULS_MAIN_CHAIN_ID, pubbyte)
	fmt.Printf("fromaddress:%s\n", fromaddress)
	fromprihex, err := GetHexPrivKey(pri.D)
	fmt.Printf("fromprihex:%s\n", fromprihex)

	//topri, err := nuls_sdk4go.GeneratePriKey()
	topri, err := BuildPrivkey(s2)
	topubbyte := GetEncodePubkey(&(topri.PublicKey), true)
	toaddress := GetNulsAddressByPub(NULS_MAIN_CHAIN_ID, topubbyte)
	fmt.Printf("toaddress:%s\n", toaddress)
	toprihex, err := GetHexPrivKey(topri.D)
	fmt.Printf("toprihex:%s\n", toprihex)

	input := CoinFrom{
		Address:      fromaddress,
		AssetChainId: NULS_MAIN_CHAIN_ID,
		AssetId:      NULS_MAIN_ASSET_ID,
		Amount:       big.NewInt(5000000),
		Nonce:        nonce, //交易顺序号，前一笔交易的hash的后8个字节的hex编码(16字符),第一笔交易为"0000000000000000"
		Locked:       NULS_NORMAL_TX_LOCKED,
	}

	output := CoinTo{
		Address:      toaddress,
		AssetChainId: NULS_MAIN_CHAIN_ID,
		AssetId:      NULS_MAIN_ASSET_ID,
		Amount:       big.NewInt(4900000),
		LockTime:     0,
	}

	inputs := list.New()
	inputs.PushBack(input)

	outputs := list.New()
	outputs.PushBack(output)

	transferDto := TransferDto{
		Inputs:  inputs,
		Outputs: outputs,
		Time:    tx_time,
		Remark:  "remark22",
	}

	trx, err := CreateTransferTrxInDetail(transferDto)
	if err != nil {
		t.Error(err.Error())
	}

	hashbytes := (trx.Hash)[:]
	fmt.Println("hashbytes:")
	fmt.Println(hashbytes)

	//hashbytes, err := hex.DecodeString(trxhash)

	bigI := pri.D
	hexpriv, err := GetHexPrivKey(bigI)

	fmt.Printf("from private int:%s\n", bigI.String())
	fmt.Printf("from private hex:%s\n", hexpriv)

	var prilist list.List
	prilist.PushBack(hexpriv)
	trx.SignTrx(&prilist)

	trxbytesss, err := trx.SerializeForHash()
	trxhex := hex.EncodeToString(trxbytesss)
	fmt.Println("signed trx hex:")
	fmt.Println(trxhex)

	if expect_trxhex != trxhex {
		t.Error("TestCreateTransferTrxInDetail :wrong trx")
	}
	fmt.Println("TestCreateTransferTrxInDetail  over..................")
}

func TestCreateTransferTrx(t *testing.T) {
	fmt.Println("TestCreateTransferTrx begin")
	expect_trxhex := "0200bbb5e05d0872656d61726b3232008c01170100016f776e306fd399498d50d83ef0d2bdf71a1922d001000100404b4c000000000000000000000000000000000000000000000000000000000008c413ee9158f97881000117010001b6acd7d3a620475fd309aa7b5fc49056d10adacc01000100a0c44a00000000000000000000000000000000000000000000000000000000000000000000000000"
	var tx_time int64 = 1575007675
	prevtxid := "305b6ad9c553775acbb82eb354fb9f3d4201f8ef143eafcdc413ee9158f97881"
	nonce := prevtxid[(len(prevtxid) - 16):]

	s1 := "2f40c87367d648874bd369442cdef92ef8a933b1560759629364fa6278063ade"
	s2 := "90a2030662e196f990300e9f5eb3d403f7245aae9db66b0a461072832429f0b5"
	pri, err := BuildPrivkey(s1)
	//pri, err := nuls_sdk4go.GeneratePriKey()
	pubbyte := GetEncodePubkey(&(pri.PublicKey), true)
	fromaddress := GetNulsAddressByPub(NULS_MAIN_CHAIN_ID, pubbyte)
	fmt.Printf("fromaddress:%s\n", fromaddress)
	fromprihex, err := GetHexPrivKey(pri.D)
	fmt.Printf("fromprihex:%s\n", fromprihex)

	//topri, err := nuls_sdk4go.GeneratePriKey()
	topri, err := BuildPrivkey(s2)
	topubbyte := GetEncodePubkey(&(topri.PublicKey), true)
	toaddress := GetNulsAddressByPub(NULS_MAIN_CHAIN_ID, topubbyte)
	fmt.Printf("toaddress:%s\n", toaddress)
	toprihex, err := GetHexPrivKey(topri.D)
	fmt.Printf("toprihex:%s\n", toprihex)

	amount := big.NewInt(4900000)
	fee := big.NewInt(NULS_DEFAULT_NORMAL_TX_FEE_PRICE)
	trx, err := CreateTransferTrx(fromaddress, nonce, toaddress, amount, fee, NULS_MAIN_CHAIN_ID, NULS_MAIN_ASSET_ID, "remark22")
	trx.Time = tx_time
	if err != nil {
		t.Error(err.Error())
	}

	hashbytes := (trx.Hash)[:]
	fmt.Println("hashbytes:")
	fmt.Println(hashbytes)

	bigI := pri.D
	hexpriv, err := GetHexPrivKey(bigI)

	fmt.Printf("from private int:%s\n", bigI.String())
	fmt.Printf("from private hex:%s\n", hexpriv)

	var prilist list.List
	prilist.PushBack(hexpriv)
	trx.SignTrx(&prilist)

	trxbytesss, err := trx.SerializeForHash()
	trxhex := hex.EncodeToString(trxbytesss)
	fmt.Println("signed trx hex:")
	fmt.Println(trxhex)

	if expect_trxhex != trxhex {
		t.Error("TestCreateTransferTrx :wrong trx")
	}
	fmt.Println("TestCreateTransferTrx  over..................")
}

func TestCreateDepositTrx(t *testing.T) {
	expect_trxhex := "0500e6bbe05d005700d0ed902e0000000000000000000000000000000000000000000000000000000100019f191c904e4e072c9d940c68120a19ec64f304feb9ae30b06d19d1740ac734c5284cb4a8eee6d2d0684684a0af234c03cff6e6bc8c01170100019f191c904e4e072c9d940c68120a19ec64f304fe01000100a056ef902e00000000000000000000000000000000000000000000000000000008c413ee9158f978810001170100019f191c904e4e072c9d940c68120a19ec64f304fe0100010000d0ed902e000000000000000000000000000000000000000000000000000000ff00000000000000"
	var tx_time int64 = 1575009254

	//prevtxid := "305b6ad9c553775acbb82eb354fb9f3d4201f8ef143eafcdc413ee9158f97881"
	//nonce := prevtxid[(len(prevtxid) - 16):]
	fromAddress := "NULSd6HgdemcQDAaiEJWq9ESMhtUHXRsNJ4KM"
	amount := big.NewInt(200000000000)
	fee := big.NewInt(100000)
	nonce := "c413ee9158f97881"
	agenthash := "b9ae30b06d19d1740ac734c5284cb4a8eee6d2d0684684a0af234c03cff6e6bc"

	deposittrx, err := CreateDepositTx(fromAddress, nonce, amount, agenthash, fee, NULS_MAIN_CHAIN_ID, NULS_MAIN_ASSET_ID)

	if err != nil {
		t.Error(err.Error())
	}
	deposittrx.Time = tx_time
	trxbytesss, err := deposittrx.SerializeForHash()
	if err != nil {
		t.Error(err.Error())
	}
	trxhex := hex.EncodeToString(trxbytesss)
	fmt.Printf("createdeposit SerializeForHash trx hex:%s\n", trxhex)

	if expect_trxhex != trxhex {
		t.Error("TestCreateDepositTrx :wrong trx")
	}
	fmt.Println("TestCreateDepositTrx  over..................")
}

func TestCreateWithdrawDepositTrx(t *testing.T) {
	expect_trxhex := "06009fbce05d0020a2e75a7b76ace482a3807da807c1bedb7dd282dd90f0aa7d454963f880cd9d4c8c011701000101646e391bfcd8d761c3e66b85b6d7fa35493ff70100010080b699552f00000000000000000000000000000000000000000000000000000008454963f880cd9d4cff011701000101646e391bfcd8d761c3e66b85b6d7fa35493ff701000100e02f98552f0000000000000000000000000000000000000000000000000000000000000000000000"
	var tx_time int64 = 1575009439
	amount := big.NewInt(203299600000)
	fromAddress := "NULSd6HgTxJyJUTPwaKjYn76oXpQWZK6sJAte"
	nonce := "454963f880cd9d4c"
	depositHash := "a2e75a7b76ace482a3807da807c1bedb7dd282dd90f0aa7d454963f880cd9d4c"
	fee := big.NewInt(NULS_DEFAULT_NORMAL_TX_FEE_PRICE)
	agenthash := "b9ae30b06d19d1740ac734c5284cb4a8eee6d2d0684684a0af234c03cff6e6bc"

	withtrx, err := CreateWithdrawDepositTx(fromAddress, nonce, amount, depositHash, agenthash, fee, NULS_MAIN_CHAIN_ID, NULS_MAIN_ASSET_ID)

	if err != nil {
		t.Error(err.Error())
	}
	withtrx.Time = tx_time
	//withtrx.SignTrx(&prilist)
	trxbytesss, err := withtrx.SerializeForHash()
	if err != nil {
		t.Error(err.Error())
	}
	trxhex := hex.EncodeToString(trxbytesss)
	fmt.Println("signed withdraw trx hex:")
	fmt.Println(trxhex)

	if expect_trxhex != trxhex {
		t.Error("TestCreateWithdrawDepositTrx :wrong trx")
	}
	fmt.Println("TestCreateWithdrawDepositTrx  over..................")
}
