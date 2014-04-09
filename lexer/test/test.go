package main

import (
    "lexer"
    "fmt"
)

func main() {

    lex := lexer.Lex("test","some text /* a comment */ more text ")
    
    for {
        item := lex.NextItem()
        fmt.Println(item)
        if item.IsEOF() {
            break
        }
    }
}