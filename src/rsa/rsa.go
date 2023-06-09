package rsa

import (
	"math/big"
)

func Sign(msg []string, privateKey string) string {
	// convert number to bigInt type
	n_from := new(big.Int)
	n_from.SetString(msg[1], 10) //base 10
	d := new(big.Int)
	d.SetString(privateKey, 10) //base 10
	msgToSign := []byte(msg[0])
	sign := encrypt(n_from, d, msgToSign)
	return sign.String()
}

func Verify(msg []string, signature string) bool {
	nfrom := new(big.Int)
	nfrom.SetString(msg[1], 10) //base 10
	efrom := new(big.Int)
	efrom.SetString(msg[2], 10) //base 10
	sign := new(big.Int)
	sign.SetString(signature, 10)
	decrypted := decrypt(nfrom, efrom, sign)

	validSign := msg[0] // this could be from: publickey, but now it is 'signedTransaction'

	return string(decrypted) == validSign
}

/* RSA Decryption. */
func decrypt(n *big.Int, d *big.Int, c *big.Int) []byte {
	// Compute m = c^d mod n
	msgAsNumber := big.NewInt(0).Exp(c, d, n)

	// Transform the number into a byte array
	msgAsByteArray := msgAsNumber.Bytes()

	return msgAsByteArray
}

/* RSA Encryption. */
func encrypt(n *big.Int, e *big.Int, m []byte) *big.Int {
	// Transform the message into a number
	msg := big.NewInt(0).SetBytes(m)

	// Compute c = m^e mod n
	cipher := big.NewInt(0).Exp(msg, e, n)

	return cipher
}
