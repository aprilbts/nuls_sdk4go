package nuls_sdk4go

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"encoding/hex"
	"math/big"

	"github.com/fomichev/secp256k1"
)

func GeneratePriKey() (*ecdsa.PrivateKey, error) {
	//return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return ecdsa.GenerateKey(secp256k1.SECP256K1(), rand.Reader)
}

type sigStruct struct {
	R *big.Int
	S *big.Int
}

/*
对digest签名
返回asn1 der编码的签名
*/
func sign(priv *ecdsa.PrivateKey, digest []byte) ([]byte, error) {
	return priv.Sign(rand.Reader, digest, nil)
}

func VerifySig(pub *ecdsa.PublicKey, hash []byte, sig []byte) bool {
	var sigs sigStruct
	_, err := asn1.Unmarshal(sig, &sigs)
	if err != nil {
		return false
	}
	return ecdsa.Verify(pub, hash, sigs.R, sigs.S)
}

/*
公钥为x,y2个big int坐标
x,y满足 y^2 = x^3 + 7
很容易就可以将它转换为压缩形式。我们只需省略掉y并改掉它的前缀。这个新前缀是根据y来决定的：前缀是 02表示y是偶数值，前缀是 03表示y是奇数值。
*/
func GetEncodePubkey(pub *ecdsa.PublicKey, needCompress bool) []byte {
	var resultbytes []byte
	var flag byte

	xbytes := pub.X.Bytes()
	var ybytes []byte
	if needCompress == true {
		if pub.Y.Bit(0) == 0 {
			flag = 2
		} else {
			flag = 3
		}
	} else {
		flag = 4
		ybytes = pub.Y.Bytes()
	}

	arrbytes := make([][]byte, 3)
	arrbytes[0] = []byte{flag}
	arrbytes[1] = xbytes
	arrbytes[2] = ybytes

	resultbytes = bytes.Join(arrbytes, []byte{})
	return resultbytes
}

func BuildPrivkey(privhex string) (*ecdsa.PrivateKey, error) {
	pribytes, err := hex.DecodeString(privhex)
	d := new(big.Int).SetBytes(pribytes)
	aprivKey := new(ecdsa.PrivateKey)
	aprivKey.D = d

	aprivKey.Curve = secp256k1.SECP256K1()
	aprivKey.PublicKey.X, aprivKey.PublicKey.Y = aprivKey.Curve.ScalarBaseMult(aprivKey.D.Bytes())
	return aprivKey, err
}
