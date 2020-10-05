package main

import (
	"fmt"
	"os"
)

func main() {
	var userPassword string
	fmt.Print("Введите предполагаемый пароль: ")
	fmt.Fscan(os.Stdin, &userPassword)
	if СheckPassword(userPassword) {
		fmt.Println("Пароль безопасен!")
	} else {
		fmt.Println("Введите новый пароль!")
	}
}
