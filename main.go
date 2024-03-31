package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/xmp-er/bencodingParserGo/processors"
	"github.com/xmp-er/bencodingParserGo/validators"
)

// Encode_bencoded_val takes a byte slice of JSON and returns its bencoded string representation along with any error encountered.
func Encode_bencoded_val(input []byte) (string, error) {
	var decoded_res map[string]interface{} //decoded result

	//unmarshal the json into a map
	err := json.Unmarshal(input, &decoded_res)

	if err != nil {
		fmt.Println("Error unmarshaling the json value ", err)
		return "", err
	}

	//encode the map into bencoded
	result, err := processors.Encode_Dictionary(decoded_res)

	if err != nil {
		fmt.Println("Error decoding the dictionary ", err)
		return "", err
	}

	return result, nil
}

// Decode_bencoded_val_as_interface takes a bencoded string and decodes it into an array of interfaces.
func Decode_bencoded_val_as_interface(input string) ([]interface{}, error) {
	//constructing a valid array of the same length as input string
	valid := make([]bool, len(input))

	//marking the integers and strings
	processors.MarkStringAndInts(input, &valid)

	result, err := processors.Decode(input, valid)

	if err != nil {
		fmt.Println("Error in decoding the Bencoded value ", err)
		return nil, err
	}

	return result, nil
}

// Decode_bencoded_val_as_JSON takes a bencoded string, decodes it, and returns the result in JSON format.
func Decode_bencoded_val_as_JSON(input string) ([]byte, error) {
	result, err := Decode_bencoded_val_as_interface(input)

	if err != nil {
		fmt.Println("Error in decoding value ", err)
		return nil, err
	}

	json_res, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error in converting to json ", err)
		return nil, err
	}
	return json_res, nil
}

// Print_Decoded_bencoded_val prompts the user for a bencoded string, decodes it, and prints the result in JSON format.
func Print_Decoded_bencoded_val() {
	var input string = ""
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Please enter the Bencoded value: ") //scanning the input value
		scanner.Scan()
		input += scanner.Text()
		temp_input := strings.Split(input, " ")
		if validators.IsValidInp(temp_input) { //validating if multiple inputs with spaces have been sent in which case, reprompt the user for input
			break
		} else {
			input = ""
			fmt.Println("Invalid input, please enter the Bencoded value as a single line with no space breaks")
		}
	}

	//decoding the Bencoded value

	//constructing a valid array of the same length as input string
	valid := make([]bool, len(input))

	//marking the integers and strings
	processors.MarkStringAndInts(input, &valid)

	result, err := processors.Decode(input, valid)

	if err != nil {
		fmt.Println("Error in decoding the Bencoded value ", err)
		return
	}

	//converting to json
	json_res, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error in converting to json ", err)
		return
	}

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("The Decoded Bencoded string is :")
		fmt.Println(string(json_res))
	}
}
