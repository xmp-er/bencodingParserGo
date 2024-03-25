package main

import (
	"fmt"

	"github.com/xmp-er/bencodingParserGo/processors"
)

func main() {

	var input string = ""
	// scanner := bufio.NewScanner(os.Stdin)

	// for {
	// 	fmt.Print("Please enter the Bencoded value: ") //scanning the input value
	// 	scanner.Scan()
	// 	input += scanner.Text()
	// 	temp_input := strings.Split(input, " ")
	// 	if validators.IsValidInp(temp_input) { //validating if multiple inputs with spaces have been sent in which case, reprompt the user for input
	// 		break
	// 	} else {
	// 		input = ""
	// 		fmt.Println("Invalid input, please enter the Bencoded value as a single line with no space breaks")
	// 	}
	// }

	//decoding the Bencoded value

	//constructing a valid array of the same length as input string
	valid := make([]bool, len("dli42eei23ee"))

	// //marking the integers and strings
	processors.MarkStringAndInts("dli42eei23ee", &valid)
	// // fmt.Println(valid)
	fmt.Println(processors.DecodeList("dli42eei23ee", valid))
	fmt.Println(input)
	fmt.Println("ended")
}
