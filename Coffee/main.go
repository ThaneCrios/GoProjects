package main

import (
	"fmt"
	"os"
)

var ourCoffeeMachine coffeeMachine
var soldEspressoCount = 0
var soldLatteCount = 0
var soldCappuccinoCount = 0
var earnedMoney = 0

type coffeeMachine struct {
	primaryMoney    int
	primaryWater    int
	primaryMilk     int
	primaryCofBeans int
	primaryCups     int
}

func main() {
	var userCommand string
	var machineWork = true
	ourCoffeeMachine.primaryMoney = 390
	ourCoffeeMachine.primaryWater = 540
	ourCoffeeMachine.primaryMilk = 400
	ourCoffeeMachine.primaryCofBeans = 120
	ourCoffeeMachine.primaryCups = 9
	for machineWork != false {
		fmt.Println("Введите команду (buy, fill, take, stat, exit):")
		fmt.Fscan(os.Stdin, &userCommand)
		if userCommand == "stat" {
			ShowStat(ourCoffeeMachine.primaryMoney, ourCoffeeMachine.primaryWater, ourCoffeeMachine.primaryMilk, ourCoffeeMachine.primaryCofBeans, ourCoffeeMachine.primaryCups)
			ShowSales(soldEspressoCount, soldLatteCount, soldCappuccinoCount, earnedMoney)
		}
		if userCommand == "buy" {
			if !buyDrink(ourCoffeeMachine.primaryWater, ourCoffeeMachine.primaryMilk, ourCoffeeMachine.primaryCofBeans, ourCoffeeMachine.primaryCups) {
				fmt.Println("Недостаточно ресурсов")
			} else {
				fmt.Println("Окей, сейчас сделаю")
			}
		}
		if userCommand == "fill" {
			fillResources()
			fmt.Println(ourCoffeeMachine.primaryWater)
			fmt.Println(ourCoffeeMachine.primaryMilk)
			fmt.Println(ourCoffeeMachine.primaryCofBeans)
			fmt.Println(ourCoffeeMachine.primaryCups)
		}
		if userCommand == "take" {
			takeMoney()
			fmt.Println(ourCoffeeMachine.primaryMoney)
		}
		if userCommand == "exit" {
			fmt.Println("Пока!")
			machineWork = false
		}
	}
}
