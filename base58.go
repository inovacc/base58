package base58

import (
	"fmt"
	"math/big"
)

var (
	bigIntermediateRadix = big.NewInt(430804206899405824) // 58**10
	alphabet             = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	b58table             = tableInitializer()
)

func tableInitializer() [256]byte {
	var table [256]byte
	for i := range table {
		table[i] = 0xFF
	}
	for i, char := range []byte(alphabet) {
		table[char] = byte(i)
	}
	return table
}

type Encoding struct {
	encode []byte
}

// StdEncoding is the standard base58 encoding, as defined in
// RFC 2119.
var StdEncoding = NewEncoding()

func NewEncoding() *Encoding {
	enc := new(Encoding)
	copy(enc.encode[:], alphabet)
	return enc
}

// Encode converts a slice of bytes to a base58 string.
func (enc *Encoding) Encode(input []byte) string {
	output := make([]byte, 0)
	num := new(big.Int).SetBytes(input)
	mod := new(big.Int)
	var primitiveNum int64
	for num.Sign() > 0 {
		num.DivMod(num, bigIntermediateRadix, mod)
		primitiveNum = mod.Int64()
		for i := 0; (num.Sign() > 0 || primitiveNum > 0) && i < 10; i++ {
			output = append(output, alphabet[primitiveNum%58])
			primitiveNum /= 58
		}
	}
	output = appendZeroBytes(output, input)
	return string(reverseByteOrder(output))
}

// AppendZeroBytes adds leading zero bytes for precise decoding.
func appendZeroBytes(output, input []byte) []byte {
	for i := 0; i < len(input) && input[i] == 0; i++ {
		output = append(output, alphabet[0])
	}
	return output
}

// ReverseByteOrder reverses the byte order of a byte slice.
func reverseByteOrder(output []byte) []byte {
	for i := 0; i < len(output)/2; i++ {
		output[i], output[len(output)-1-i] = output[len(output)-1-i], output[i]
	}
	return output
}

func (enc *Encoding) EncodeToString(src []byte) string {
	return enc.Encode(src)
}

// Decode converts a base58 string to a slice of bytes, and returns
// an error if the input is not valid base58 string.
func (enc *Encoding) Decode(input string) (output []byte, err error) {
	result, err := calculateResult(input)
	if err != nil {
		return output, err
	}
	tmpBytes := result.Bytes()
	numZeros := countZeros(input)
	length := numZeros + len(tmpBytes)
	output = make([]byte, length)
	copy(output[numZeros:], tmpBytes)
	return
}

// CalculateResult calculates the `big.Int` result for Decode function.
func calculateResult(input string) (*big.Int, error) {
	result := big.NewInt(0)
	tmpBig := new(big.Int)

	for i := 0; i < len(input); {
		var a, m int64 = 0, 58
		var err error
		a, m, i, err = createTmpVariables(input, a, m, i)
		if err != nil {
			return result, err
		}
		result.Mul(result, tmpBig.SetInt64(m))
		result.Add(result, tmpBig.SetInt64(a))
	}
	return result, nil
}

// CreateTmpVariables creates intermediate variables for Decode function.
func createTmpVariables(input string, a int64, m int64, i int) (int64, int64, int, error) {
	for f := true; i < len(input) && (f || i%10 != 0); i++ {
		tmp := b58table[input[i]]
		if tmp == 255 {
			msg := "invalid Base58 input string at character \"%c\", position %d"
			return a, m, i, fmt.Errorf(msg, input[i], i)
		}
		a = a*58 + int64(tmp)
		if !f {
			m *= 58
		}
		f = false
	}
	return a, m, i, nil
}

// CountZeros counts the leading zeros in a string.
func countZeros(input string) int {
	var numZeros int
	for numZeros = 0; numZeros < len(input); numZeros++ {
		if input[numZeros] != '1' {
			break
		}
	}
	return numZeros
}

func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	return enc.Decode(s)
}
