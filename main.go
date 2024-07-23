package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"

	"golang.org/x/term"
)

func rgb (i int) (int, int, int) {
	var f = 0.1
	return int(math.Sin(f*float64(i)+0)*127 + 128),
		int(math.Sin(f*float64(i)+2*math.Pi/3)*127 + 128),
		int(math.Sin(f*float64(i)+4*math.Pi/3)*127 + 128)
}

func print(output []rune) {
	for j := 0; j < len(output); j++ {
		r, g ,b := rgb(j)
		fmt.Printf("\033[38;2;%d;%d;%dm%c\033[0m", r, g, b, output[j])
	}

	fmt.Println()
}

type Quote struct {
	Id string `json:"_id"`
	Content string `json:"content"`
	Author string `json:"author"`

}

func fetchQuote () Quote {
	var quote []Quote
	resp, err := http.Get("https://api.quotable.io/quotes/random")
	if err != nil {
		log.Print("HTTP ERR: ", err)
	}
	defer resp.Body.Close()
	jsonErr := json.NewDecoder(resp.Body).Decode(&quote)
	
	

	if jsonErr != nil {
		log.Print(jsonErr)
	}
	return quote[0]
}

func main () {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		fmt.Println("Error getting terminal size:", err)
		return
	}

	halfway := width / 3
	info, _ := os.Stdin.Stat()
	var output []rune

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("Press Enter to generate quote...")
    }

    reader := bufio.NewReader(os.Stdin)
    for {
        input, _, err := reader.ReadRune()
        if err != nil && err == io.EOF {
            break
        }
		if input == '\n' {
			quote := fetchQuote()
			var coloredQuote string
			for i, ch := range quote.Content {
				r, g, b := rgb(i)
				char := fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", r, g, b, ch)
			
			coloredQuote += char
			}
			coloredQuote += "\033[0m"
			fmt.Println(coloredQuote)
			fmt.Println()

			fmt.Printf("%*s%s\n",  halfway, "",  quote.Author)
		}
    }
    print(output)
	
} 
