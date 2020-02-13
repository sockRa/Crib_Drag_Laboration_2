package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {

	//Step 0: Read the ciphertext file and save to ciphers
	ciphers := ReadCipherTextFromFile()
	//Retrive the target cipher which is last in the ciphers array
	targetCipher := DecodeHexBytes([]byte(ciphers[len(ciphers)-1]))

	guess := AskForGuess()

	for i := 0; i < len(ciphers)-1; i++ {
		cipher := DecodeHexBytes([]byte(ciphers[i]))
		cipherXOR := StringXOR(cipher, targetCipher)

		fmt.Printf("\n\n~~~~~~~~~~~~[ Message %d XOR Message %d ]~~~~~~~~~~~~\n", i+1, len(ciphers))
		CribSearch(cipherXOR, guess)
	}

	main()
}

//CribSearch will XOR crib with target and print the equivilent ASCII
func CribSearch(cipherText, crib []byte) {

	res := StringXOR(cipherText[:len(crib)], crib)
	fmt.Printf("=> %s\t", res)
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
		temp[i] = a[i]
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

// AskForGuess asks the user to submit a guess of what word the ciphertext may contain
func AskForGuess() []byte {

	fmt.Println("\n\nGuess a word that might appear in one of the ciphertext's")
	fmt.Print("-> ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := []byte(scanner.Text())

	return text
}

//ReadCipherTextFromFile will read the given ciphertexts and add them to a string array.
func ReadCipherTextFromFile() []string {

	file, err := os.Open("Ciphers")

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
