package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"math"
	"os"
)

//var isStringAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString
//var calculateEntropy = false

func main() {

	//Step 0: Read the ciphertext file and save to ciphers
	ciphers := ReadCipherTextFromFile()
	//Retrive the target cipher which is last in the ciphers array
	targetCipher := DecodeHexBytes([]byte(ciphers[len(ciphers)-1]))

	//Step 1: Enter a guess
	guess := AskForGuess()

	for i := 0; i < len(ciphers)-1; i++ {
		//XOR cipher[i] with targetCipher
		cipher := DecodeHexBytes([]byte(ciphers[i]))
		cipherXOR := StringXOR(cipher, targetCipher) //, calculateEntropy)

		//Crib drag the guess on cipherXOR
		fmt.Printf("\n\n~~~~~~~~~~~~[ Message %d XOR Message %d ]~~~~~~~~~~~~\n", i+1, len(ciphers))
		CribSearch(cipherXOR, guess)
		fmt.Printf("%s\n", EncodeHexBytes(cipherXOR))
	}

	//main()
}

//CribSearch will XOR crib with target and print the equivilent ASCII
func CribSearch(cipherText, crib []byte) {

	var leftSlice int
	var rightSlice int
	//calculateEntropy = true
	loopLength := (len(cipherText) - len(crib))

	for i := 0; i <= loopLength; i++ {

		// Initial assignment for choosing the slice-length of the cipher text, based on the length of the crib word.
		if i == 0 {
			leftSlice = 0
			rightSlice = len(crib)
		} else {
			// We want the slice to move in equally large steps as the crib-length
			// this is a preperation for the XOR operation with the crib value
			leftSlice++
			rightSlice++

			if int(rightSlice) > len(cipherText) {
				rightSlice = len(cipherText)
			}
		}

		// Extract the prepared slice from the ciphertext
		cipherTextSlice := cipherText[int(leftSlice):int(rightSlice)]

		// Perform XOR with the ciphertext slice together with the crib.
		// Optional: Calculate the entropy for each word
		// If the entropy is low, it could suggest that the output is a word in english
		// and not jibberish.
		res := StringXOR(cipherTextSlice, crib) //, calculateEntropy)

		// Configure this to change the number of rows that you want to display.
		if i%3 == 0 && i != 0 {
			fmt.Printf("\n")
		}

		// Print result from crib drag with optional entropy
		fmt.Printf("\t[%d] => %s\t", i, res)

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

	// If this method StringXOR is called when performin crib drag
	// this bool should be true if the user want's to calculate
	// the entropy on the output after the crib drag operation.

	return result
}

//CalculateEntropy will calculate the average byte in the result
func CalculateEntropy(result []byte) float64 {
	var totalEntropy float64

	for i := 0; i < len(result)-1; i++ {
		//Filter out space in the calculation
		if int(result[i]) != 32 {
			// Calculate the difference between each byte
			totalEntropy += math.Abs(float64(result[i]) - float64(result[i+1]))
		}
	}

	// Return the total difference divided by the length of the string
	// if the score is low, it could suggest that the output is a possible word.
	return totalEntropy / float64(len(result))
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
