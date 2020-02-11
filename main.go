package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
)

var isStringAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString

func main() {

	//Step 0: Read the ciphertext file and save to ciphers
	ciphers := ReadCipherTextFromFile()
	//Retrive the target cipher which is last in the ciphers array
	targetCipher := DecodeHexBytes([]byte(ciphers[len(ciphers)-1]))

	//Step 1: Enter a guess
	guess := []byte("hej") //AskForGuess()

	for i := 0; i < len(ciphers)-1; i++ {
		//XOR cipher[i] with targetCipher
		cipher := DecodeHexBytes([]byte(ciphers[i]))
		cipherXOR := StringXOR(cipher, targetCipher)
		fmt.Printf("Hex cipher: %s\n", EncodeHexBytes(cipherXOR))

		//fmt.Printf("\n%s XOR %s => %s\n", EncodeHexBytes(cipher), EncodeHexBytes(targetCipher), EncodeHexBytes(cipherXOR))

		//Crib drag the guess on cipherXOR
		fmt.Printf("~~~~~~~~~~~~[ Message %d XOR Message %d ]~~~~~~~~~~~~\n\n", i+1, len(ciphers))
		CribSearch(cipherXOR, guess)
	}

	main()
}

//CribSearch will XOR crib with target and print the equivilent ASCII
func CribSearch(cipherText, crib []byte) {

	/*if len(cipherText) == len(crib) {
		res = StringXOR(cipherText, crib)
		fmt.Printf("\nFinal message => %s\n NOTE: Enter the final message to get the plaintext from the other cyphertext\n", res)
	}*/

	//Calculate the maximum iterations needed for completion
	//maxIterations := math.Ceil(float64(len(cipherText)) / float64(len(crib)))

	var leftSlice int
	var rightSlice int
	//var rightCribSlice int

	for i := 0; i < (len(cipherText) - len(crib) + 1); i++ {

		// Initial assignment for choosing the slice-length of the cipher text, based on the length of the crib word.
		if i == 0 {
			leftSlice = 0
			rightSlice = len(crib)
		} else {
			// We want the slice to move in equally large steps as the crib-length
			// this is a preperation for the XOR operation with the crib value
			leftSlice++  //= float64(len(crib) * i)
			rightSlice++ //= float64(len(crib) * (i + 1))

			if int(rightSlice) > len(cipherText) {
				rightSlice = len(cipherText)
			}
		}

		//right crib slice assigment
		/*if leftSlice+len(crib) > len(cipherText) {
			rightCribSlice = (leftSlice + len(crib)) - len(cipherText)
		} else {
			rightCribSlice = len(crib)
		}*/

		cipherTextSlice := cipherText[int(leftSlice):int(rightSlice)]
		//cribSlice := crib[:int(rightCribSlice)]
		res := StringXOR(cipherTextSlice, crib)
		fmt.Printf("[%d] => %s\n", i, res)
	}
}

//StringXOR will perfom XOR on two byte arrays and return the result
func StringXOR(a, b []byte) []byte {

	if len(a) > len(b) {
		b = AddZeroes(b, len(a)-len(b))
	} else if len(b) > len(a) {
		a = AddZeroes(a, len(b)-len(a))
	}

	result := make([]byte, len(a))
	for i := 0; i < len(a); i++ {
		result[i] = a[i] ^ b[i]
	}

	return result
}

//AddZeroes adds zeros to the array and returns the result
func AddZeroes(a []byte, size int) []byte {

	temp := make([]byte, len(a)+size)

	for i := 0; i < len(a); i++ {
		temp = append(temp[:i+size], a[i])
	}
	return temp
}

//DecodeHexBytes will convert from hex to bytes and return the result
func DecodeHexBytes(hexBytes []byte) []byte {
	result := make([]byte, hex.DecodedLen(len(hexBytes)))
	_, err := hex.Decode(result, hexBytes)
	if err != nil {
		return nil
	}
	return result
}

//EncodeHexBytes convert bytes to hex
func EncodeHexBytes(input []byte) []byte {
	result := make([]byte, hex.EncodedLen(len(input)))
	hex.Encode(result, input)
	return result
}

// AskForGuess asks the user to submit a guess of what word the ciphertext may contain
func AskForGuess() []byte {

	fmt.Println("\nGuess a word that might appear in one of the ciphertext's")
	fmt.Print("-> ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := []byte(scanner.Text())

	//return EncodeHexBytes(text)
	return text
}

//ReadCipherTextFromFile will read the given ciphertexts and add them to a string array.
func ReadCipherTextFromFile() []string {

	file, err := os.Open("Ciphers_debug")

	if err != nil {
		fmt.Printf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	var ciphertext []string

	for scanner.Scan() {
		ciphertext = append(ciphertext, scanner.Text())
	}

	file.Close()

	return ciphertext
}
