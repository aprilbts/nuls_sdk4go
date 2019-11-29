package nuls_sdk4go

import (
	"bytes"
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"strings"
	"time"

	"golang.org/x/crypto/ripemd160"
)

const normal_PRICE_PRE_1024_BYTES = 100000

//const CROSSTX_PRICE_PRE_1024_BYTES = 1000000

const (
	TxType_COIN_BASE = 1
	TxType_TRANSFER  = 2
	//TxType_ACCOUNT_ALIAS            = 3
	//TxType_REGISTER_AGENT           = 4
	TxType_DEPOSIT        = 5
	TxType_CANCEL_DEPOSIT = 6
	/*
		TxType_YELLOW_PUNISH            = 7
		TxType_RED_PUNISH               = 8
		TxType_STOP_AGENT               = 9
		TxType_CROSS_CHAIN              = 10
		TxType_REGISTER_CHAIN_AND_ASSET = 11
		TxType_DESTROY_CHAIN_AND_ASSET  = 12
		TxType_ADD_ASSET_TO_CHAIN       = 13
		TxType_REMOVE_ASSET_FROM_CHAIN  = 14
		TxType_CREATE_CONTRACT          = 15
		TxType_CALL_CONTRACT            = 16
		TxType_DELETE_CONTRACT          = 17
		TxType_CONTRACT_TRANSFER        = 18
		TxType_CONTRACT_RETURN_GAS      = 19
		TxType_CONTRACT_CREATE_AGENT    = 20
		TxType_CONTRACT_DEPOSIT         = 21
		TxType_pCONTRACT_CANCEL_DEPOSIT = 22
		TxType_CONTRACT_STOP_AGENT      = 23
		TxType_VERIFIER_CHANGE          = 24
		TxType_VERIFIER_INIT            = 25
	*/
)

const (
	//NULS_MAIN_NET_VERSION = 1

	NULS_MAIN_CHAIN_ID     = 1
	NULS_MAIN_CHAIN_PREFIX = "NULS"

	NULS_TESTNET_CHAIN_ID     = 2
	NULS_TESTNET_CHAIN_PREFIX = "tNULS"

	NULS_MAIN_ASSET_ID    = 1
	NULS_TESTNET_ASSET_ID = 1

	NULS_DEFAULT_NORMAL_TX_FEE_PRICE = 100000
	NULS_DEFAULT_OTHER_TX_FEE_PRICE  = 100000

	NULS_NORMAL_TX_LOCKED = 0 //非解锁交易
)

/**
// address : prefix + NULSLENGTHPREFIX + base58(chainId(uint16)+address_type(byte)+hash160 + xor)
// raw address bytes : chainId(uint16)+address_type(byte)+hash160
*/
func getAddressBytes(addressString string) ([]byte, error) {
	var result []byte
	var prefixlen = 0
	if strings.HasPrefix(addressString, NULS_MAIN_CHAIN_PREFIX) {
		prefixlen = len(NULS_MAIN_CHAIN_PREFIX)

	} else if strings.HasPrefix(addressString, NULS_TESTNET_CHAIN_PREFIX) {
		prefixlen = len(NULS_TESTNET_CHAIN_PREFIX)
	} else {
		return result, errors.New("can't regonize address prefix")
	}

	base58address := addressString[(prefixlen + 1):]
	addrssbytes := Base58Decode(base58address)
	if len(addrssbytes) < 4 {
		return result, errors.New("invalid address")
	}

	return addrssbytes[0 : len(addrssbytes)-1], nil
}

//#########################################

type CoinFrom struct {
	Address      string
	AssetChainId int //  NULS_MAIN_CHAIN_ID 或 NULS_TESTNET_CHAIN_ID
	AssetId      int //NULS_MAIN_ASSET_ID
	Amount       *big.Int
	Nonce        string //长度16个字符，nonce设为账户最近一笔发出的交易ID的后16个字符，如果是第一笔交易 nonce设为"0000000000000000",
	Locked       byte   //0
}

func (coinfrom CoinFrom) serializeToStream(buffer *bytes.Buffer) error {
	addressbytes, err := getAddressBytes(coinfrom.Address)
	if err != nil {
		return err
	}
	err = bufferWriteBytesWithLength(buffer, addressbytes, len(addressbytes))
	if err != nil {
		return err
	}
	err = bufferWriteUInt16(buffer, (uint16)(coinfrom.AssetChainId))
	if err != nil {
		return err
	}
	err = bufferWriteUInt16(buffer, (uint16)(coinfrom.AssetId))
	if err != nil {
		return err
	}
	err = bufferWriteBigInt(buffer, coinfrom.Amount)
	if err != nil {
		return err
	}
	noncebytes, err := hex.DecodeString(coinfrom.Nonce)
	if err != nil {
		return err
	}
	err = bufferWriteBytesWithLength(buffer, noncebytes, len(noncebytes))
	if err != nil {
		return err
	}
	err = buffer.WriteByte(coinfrom.Locked)
	return err
}

type CoinTo struct {
	Address      string
	AssetChainId int
	AssetId      int
	Amount       *big.Int
	LockTime     int64 //0
}

func (cointo CoinTo) serializeToStream(buffer *bytes.Buffer) error {
	addressbytes, err := getAddressBytes(cointo.Address)
	if err != nil {
		return err
	}
	err = bufferWriteBytesWithLength(buffer, addressbytes, len(addressbytes))
	if err != nil {
		return err
	}
	err = bufferWriteUInt16(buffer, (uint16)(cointo.AssetChainId))
	if err != nil {
		return err
	}
	err = bufferWriteUInt16(buffer, (uint16)(cointo.AssetId))
	if err != nil {
		return err
	}
	err = bufferWriteBigInt(buffer, cointo.Amount)
	if err != nil {
		return err
	}
	err = bufferWriteInt64(buffer, cointo.LockTime)
	return err
}

type TransferDto struct {
	Inputs  *list.List //CoinFrom list
	Outputs *list.List // CoinTo list
	Time    int64      //可为0
	Remark  string     //转账信息 可不填
}

type Transaction struct {
	Type                 int // TxType_....
	CoinData             []byte
	TxData               []byte
	Time                 int64
	TransactionSignature []byte
	Remark               []byte
	Hash                 [32]byte
	//BlockHeight          []byte
	//Status               int
	//Size                 int
	//InBlockIndex         int
	Inputs  *list.List //CoinFrom list
	Outputs *list.List // CoinTo list
}

func (trx Transaction) SerializeForHash() (bytearray []byte, err error) {
	var buffer bytes.Buffer
	err = bufferWriteUInt16(&buffer, (uint16)(trx.Type))
	if err != nil {
		return bytearray, err
	}
	err = bufferWriteUInt32(&buffer, (uint32)(trx.Time))
	if err != nil {
		return bytearray, err
	}
	err = bufferWriteBytesWithLength(&buffer, trx.Remark, len(trx.Remark))
	if err != nil {
		return bytearray, err
	}
	err = bufferWriteBytesWithLength(&buffer, trx.TxData, len(trx.TxData))
	if err != nil {
		return bytearray, err
	}
	err = bufferWriteBytesWithLength(&buffer, trx.CoinData, len(trx.CoinData))
	if err != nil {
		return bytearray, err
	}
	return buffer.Bytes(), err
}

func (trx Transaction) Serialize() ([]byte, error) {
	var bytearray []byte
	r, err := trx.SerializeForHash()
	if err != nil {
		return bytearray, err
	}
	buffer := bytes.NewBuffer(r)
	err = bufferWriteBytesWithLength(buffer, trx.TransactionSignature, len(trx.TransactionSignature))
	return buffer.Bytes(), err
}

func (ptrx *Transaction) generateTrxHash() error {
	data, err := ptrx.SerializeForHash()
	if err != nil {
		return err
	}
	h1 := sha256.Sum256(data)
	h2 := sha256.Sum256(h1[:])
	ptrx.Hash = h2
	return nil
}

func (ptrx *Transaction) generateCoinData() error {
	var buffer bytes.Buffer
	fromCount := ptrx.Inputs.Len()
	err := bufferWriteEncodedVarInt(&buffer, int64(fromCount))
	if err != nil {
		return err
	}
	for e := ptrx.Inputs.Front(); e != nil; e = e.Next() {
		coinfrom, ok := (e.Value).(CoinFrom)
		if ok {
			err = coinfrom.serializeToStream(&buffer)
			if err != nil {
				return err
			}
		} else {
			return errors.New("transferDto.Inputs element type is not CoinFrom type!!!")
		}
	}
	toCount := ptrx.Outputs.Len()
	err = bufferWriteEncodedVarInt(&buffer, int64(toCount))
	if err != nil {
		return err
	}
	for e := ptrx.Outputs.Front(); e != nil; e = e.Next() {
		cointo, ok := (e.Value).(CoinTo)
		if ok {
			err = cointo.serializeToStream(&buffer)
			if err != nil {
				return err
			}
		} else {
			return errors.New("transferDto.Outputs element type is not CoinTo type!!!")
		}
	}
	ptrx.CoinData = buffer.Bytes()
	return nil
}

// one address -> another address
func CreateTransferTrx(fromAddress string, nonce string, toAddress string, amount *big.Int, fee *big.Int, chainId int, assetId int, remarkInfo string) (Transaction, error) {
	inputamount := big.NewInt(0)
	inputamount.Add(fee, amount)
	input := CoinFrom{
		Address:      fromAddress,
		AssetChainId: chainId,
		AssetId:      assetId,
		Amount:       inputamount,
		Nonce:        nonce, //交易顺序号，前一笔交易的hash的后8个字节的hex编码(16字符),第一笔交易为"0000000000000000"
		Locked:       NULS_NORMAL_TX_LOCKED,
	}

	output := CoinTo{
		Address:      toAddress,
		AssetChainId: chainId,
		AssetId:      assetId,
		Amount:       amount,
		LockTime:     0,
	}

	inputs := list.New()
	inputs.PushBack(input)

	outputs := list.New()
	outputs.PushBack(output)

	transferDto := TransferDto{
		Inputs:  inputs,
		Outputs: outputs,
		Time:    time.Now().Unix(),
		Remark:  remarkInfo,
	}
	return CreateTransferTrxInDetail(transferDto)
}

/*fee = coinfroms_amount_sum - cointo_amount_sum*/
func CreateTransferTrxInDetail(transferDto TransferDto) (Transaction, error) {
	var remarkbytes []byte
	if len(transferDto.Remark) > 0 {
		remarkbytes = []byte(transferDto.Remark)
	}
	transferTime := transferDto.Time
	if transferTime == 0 {
		transferTime = time.Now().Unix()
	}
	//bytesbuffer := bytes.NewBuffer()
	trx := Transaction{
		Time:    transferTime,
		Remark:  remarkbytes,
		Type:    TxType_TRANSFER,
		Inputs:  transferDto.Inputs,
		Outputs: transferDto.Outputs,
	}
	//coin data
	err := trx.generateCoinData()
	if err != nil {
		return trx, err
	}

	err = trx.generateTrxHash()
	if err != nil {
		return trx, err
	}
	if !trx.checkTrxFeeEnough(transferDto.Inputs, transferDto.Outputs) {
		return trx, errors.New("trx fee not enough")
	}
	return trx, err
}

func CreateWithdrawDepositTx(fromAddress string, nonce string, amount *big.Int, depositHash string, agenthash string, fee *big.Int, chainId int, assetId int) (Transaction, error) {
	outputamount := big.NewInt(0)
	outputamount.Sub(amount, fee)
	coinFrom := CoinFrom{
		Address:      fromAddress,
		AssetChainId: chainId,
		AssetId:      assetId,
		Amount:       amount, //包含手续费，最小为2000Nuls（即200000000000）
		Nonce:        nonce,  //交易顺序号，前一笔交易的hash的后8个字节的hex编码(16字符),第一笔交易为"0000000000000000"
		Locked:       0xFF,
	}
	coinFroms := list.New()
	coinFroms.PushBack(coinFrom)

	coinTos := list.New()
	cointo := CoinTo{
		Address:      coinFrom.Address,
		AssetChainId: coinFrom.AssetChainId,
		AssetId:      coinFrom.AssetId,
		Amount:       outputamount,
		LockTime:     0,
	}
	coinTos.PushBack(cointo)

	trx := Transaction{
		Time:    time.Now().Unix(),
		Type:    TxType_CANCEL_DEPOSIT,
		Inputs:  coinFroms,
		Outputs: coinTos,
	}

	err := trx.generateCoinData()
	if err != nil {
		return trx, err
	}

	//tx data
	trx.TxData, err = hex.DecodeString(depositHash)
	if err != nil {
		return trx, err
	}

	err = trx.generateTrxHash()
	if err != nil {
		return trx, err
	}

	if !trx.checkTrxFeeEnough(coinFroms, coinTos) {
		return trx, errors.New("trx fee not enough")
	}
	return trx, err
}

func CreateDepositTx(fromAddress string, nonce string, amount *big.Int, agenthash string, fee *big.Int, chainId int, assetId int) (Transaction, error) {
	inputamount := big.NewInt(0)
	inputamount.Add(fee, amount)
	coinFrom := CoinFrom{
		Address:      fromAddress,
		AssetChainId: chainId,
		AssetId:      assetId,
		Amount:       inputamount, //包含手续费，最小为2000Nuls（即200000000000）
		Nonce:        nonce,       //交易顺序号，前一笔交易的hash的后8个字节的hex编码(16字符),第一笔交易为"0000000000000000"
		Locked:       NULS_NORMAL_TX_LOCKED,
	}
	coinFroms := list.New()
	coinFroms.PushBack(coinFrom)

	coinTos := list.New()
	cointo := CoinTo{
		Address:      coinFrom.Address,
		AssetChainId: coinFrom.AssetChainId,
		AssetId:      coinFrom.AssetId,
		Amount:       amount,
		LockTime:     0xFF,
	}
	coinTos.PushBack(cointo)

	trx := Transaction{
		Time:    time.Now().Unix(),
		Type:    TxType_DEPOSIT,
		Inputs:  coinFroms,
		Outputs: coinTos,
	}

	err := trx.generateCoinData()
	if err != nil {
		return trx, err
	}

	//tx data
	var buffer bytes.Buffer
	bufferWriteBigInt(&buffer, amount)
	addressbytes, err := getAddressBytes(fromAddress)
	if err != nil {
		return trx, err
	}
	bufferWriteBytes(&buffer, addressbytes, len(addressbytes))
	agentBytes, err := hex.DecodeString(agenthash)
	if err != nil {
		return trx, err
	}
	bufferWriteBytes(&buffer, agentBytes, len(agentBytes))
	trx.TxData = buffer.Bytes()

	err = trx.generateTrxHash()
	if err != nil {
		return trx, err
	}

	if !trx.checkTrxFeeEnough(coinFroms, coinTos) {
		return trx, errors.New("trx fee not enough")
	}
	return trx, err
}

func getXor(body []byte) byte {
	var xor byte = 0x00
	for i := 0; i < len(body); i++ {
		xor ^= body[i]
	}
	return xor
}

// address : prefix + NULSLENGTHPREFIX + base58(chainId(uint16)+address_type(byte)+hash160 + xor)
func GetNulsAddressByPub(chainId int, pubkey []byte) string {
	NULSLENGTHPREFIX := [...]string{"", "a", "b", "c", "d", "e"}
	DEFAULT_ADDRESS_TYPE := byte(1)
	h1 := sha256.Sum256(pubkey)
	h2 := ripemd160.New()
	h2.Write(h1[:])
	hash160 := h2.Sum(nil)

	var buffer bytes.Buffer
	bufferWriteUInt16(&buffer, uint16(chainId))
	buffer.WriteByte(DEFAULT_ADDRESS_TYPE)
	bufferWriteBytes(&buffer, hash160, len(hash160))

	xor := getXor(buffer.Bytes())
	buffer.WriteByte(xor)
	rawaddress := Base58Encode(buffer.Bytes())
	var prefix string
	if chainId == NULS_MAIN_CHAIN_ID {
		prefix = NULS_MAIN_CHAIN_PREFIX
	} else if chainId == NULS_TESTNET_CHAIN_ID {
		prefix = NULS_TESTNET_CHAIN_PREFIX
	}
	address := prefix + NULSLENGTHPREFIX[len(prefix)] + rawaddress
	return address
}

const kB = 1024

/**
 * 根据交易大小计算需要交纳的手续费
 * According to the transaction size calculate the handling fee.
 * @param size 交易大小/size of the transaction
 */
func getNormalTxFee(size int) *big.Int {
	price := big.NewInt(normal_PRICE_PRE_1024_BYTES)
	sizekb := int64(size / kB)
	kbs := big.NewInt(sizekb)
	fee := new(big.Int)
	fee.Mul(kbs, price)
	//fee := big.Int.Mul(kbs, price)

	if size%kB != 0 {
		fee.Add(fee, price)
	}
	return fee
}

const p2PHKSignature_SERIALIZE_LENGTH = 110

func getNeedSignAddressMap(coinFroms *list.List) map[string]int {
	addresses := make(map[string]int)
	for e := coinFroms.Front(); e != nil; e = e.Next() {
		coinfrom, ok := (e.Value).(CoinFrom)
		_, ok = addresses[coinfrom.Address]
		if !ok {
			addresses[coinfrom.Address] = 1
		}
	}
	//count := len(addresses)
	return addresses
}

// sigCount : 签名的数量
func (ptrx *Transaction) CaculateTxFee(sigCount int) (*big.Int, error) {
	fee := big.NewInt(0)
	if sigCount < 0 {
		return fee, errors.New("sigCount must > 0")
	}
	bytess, err := ptrx.SerializeForHash()
	if err != nil {
		return fee, err
	}
	size := len(bytess)
	size = size + sigCount*p2PHKSignature_SERIALIZE_LENGTH
	fee = getNormalTxFee(size)
	return fee, nil
}

func (ptrx *Transaction) checkTrxFeeEnough(coinFroms *list.List, coinTos *list.List) bool {
	fee, err := ptrx.CaculateTxFee(len(getNeedSignAddressMap(coinFroms)))
	if err != nil {
		return false
	}
	fromSum := big.NewInt(0)
	toSum := big.NewInt(0)
	for e := coinFroms.Front(); e != nil; e = e.Next() {
		coinfrom, ok := (e.Value).(CoinFrom)
		if ok == false {
			return false
		}
		fromSum = fromSum.Add(fromSum, coinfrom.Amount)
	}

	for e := coinTos.Front(); e != nil; e = e.Next() {
		cointo, ok := (e.Value).(CoinTo)
		if ok == false {
			return false
		}
		toSum = toSum.Add(toSum, cointo.Amount)
	}

	actualfee := big.NewInt(0)
	actualfee = actualfee.Sub(fromSum, toSum)

	if actualfee.Cmp(fee) < 0 {
		return false
	}
	return true
}

func (ptrx *Transaction) SignTrx(priHexlist *list.List) error {
	var buffer bytes.Buffer
	if priHexlist.Len() <= 0 {
		return errors.New("priHexlist is empty")
	}
	addressesMap := getNeedSignAddressMap(ptrx.Inputs)
	for p := priHexlist.Front(); p != nil; p = p.Next() {
		priHex, ok := (p.Value).(string)
		if ok == false {
			return errors.New("priHexlist must string list")
		}
		priv, err := BuildPrivkey(priHex)
		if err != nil {
			return err
		}
		sigbyte, err := sign(priv, ptrx.Hash[:])
		if err != nil {
			return err
		}

		pubbytes := GetEncodePubkey(&priv.PublicKey, true)
		addr := GetNulsAddressByPub(1, pubbytes)

		_, ok = addressesMap[addr]
		if ok == false {
			return errors.New("no coinfrom address match this key")
		}
		buffer.WriteByte(byte(len(pubbytes)))
		bufferWriteBytes(&buffer, pubbytes, len(pubbytes))

		bufferWriteBytesWithLength(&buffer, sigbyte, len(sigbyte))

	}
	ptrx.TransactionSignature = buffer.Bytes()
	return nil
}

/*
type DepositDto struct {
	Adress    string
	Deposit   *big.Int
	AgentHash string
	Input     CoinFrom
}

type WithDrawDto struct {
	Adress      string
	DepositHash string
	Price       *big.Int //手续费单价 required = false
	AgentHash   string
	Input       CoinFrom
}
*/

/*
func CreateDepositTx(depositDto DepositDto) (Transaction, error) {

	//coin data
	coinFroms := list.New()
	coinFrom := depositDto.Input
	coinFroms.PushBack(coinFrom)
	coinTos := list.New()
	cointo := CoinTo{
		Address:      coinFrom.Address,
		AssetChainId: coinFrom.AssetChainId,
		AssetId:      coinFrom.AssetId,
		Amount:       depositDto.Deposit,
		LockTime:     0xFF,
	}
	coinTos.PushBack(cointo)

	trx := Transaction{
		Time:    time.Now().Unix(),
		Type:    TxType_DEPOSIT,
		Inputs:  coinFroms,
		Outputs: coinTos,
	}

	err := trx.generateCoinData()
	if err != nil {
		return trx, err
	}

	//tx data
	var buffer bytes.Buffer
	bufferWriteBigInt(&buffer, depositDto.Deposit)
	addressbytes, err := getAddressBytes(depositDto.Adress)
	if err != nil {
		return trx, err
	}
	bufferWriteBytes(&buffer, addressbytes, len(addressbytes))
	agentBytes, err := hex.DecodeString(depositDto.AgentHash)
	if err != nil {
		return trx, err
	}
	bufferWriteBytes(&buffer, agentBytes, len(agentBytes))
	trx.TxData = buffer.Bytes()

	err = trx.generateTrxHash()
	if err != nil {
		return trx, err
	}

	if !trx.checkTrxFeeEnough(coinFroms, coinTos) {
		return trx, errors.New("trx fee not enough")
	}
	return trx, err
}
*/
/*
func CreateWithdrawDepositTx(withDrawDto WithDrawDto) (Transaction, error) {
	//coin data
	coinFroms := list.New()
	coinFrom := withDrawDto.Input
	coinFrom.Locked = 0xFF
	coinFroms.PushBack(coinFrom)
	coinTos := list.New()
	toAmount := big.NewInt(0)
	cointo := CoinTo{
		Address:      coinFrom.Address,
		AssetChainId: coinFrom.AssetChainId,
		AssetId:      coinFrom.AssetId,
		Amount:       toAmount.Sub(coinFrom.Amount, withDrawDto.Price),
		LockTime:     0,
	}
	coinTos.PushBack(cointo)

	trx := Transaction{
		Time:    time.Now().Unix(),
		Type:    TxType_CANCEL_DEPOSIT,
		Inputs:  coinFroms,
		Outputs: coinTos,
	}

	err := trx.generateCoinData()
	if err != nil {
		return trx, err
	}

	//tx data
	trx.TxData, err = hex.DecodeString(withDrawDto.DepositHash)
	if err != nil {
		return trx, err
	}

	err = trx.generateTrxHash()
	if err != nil {
		return trx, err
	}

	if !trx.checkTrxFeeEnough(coinFroms, coinTos) {
		return trx, errors.New("trx fee not enough")
	}
	return trx, err
}
*/
