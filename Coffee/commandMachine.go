package main

import (
	"fmt"
	"os"
	"strconv"
)

func ShowStat(money, water, milk, cofBeans, cups int) {
	fmt.Println("В кофемашине:")
	fmt.Println(strconv.Itoa(water) + " мл воды")
	fmt.Println(strconv.Itoa(milk) + " мл молока")
	fmt.Println(strconv.Itoa(cofBeans) + " г кофейных зёрен")
	fmt.Println(strconv.Itoa(cups) + " чашек")
	fmt.Println(strconv.Itoa(money) + " рублей")
}

func ShowSales(soldEspresso, soldLatte, soldCappuccino, earned int) {
	fmt.Println("Всего продано: " + strconv.Itoa(soldEspresso+soldLatte+soldCappuccino) + " напитков на " + strconv.Itoa(earned) + " рублей")
	fmt.Println(strconv.Itoa(soldEspresso) + " - espresso")
	fmt.Println(strconv.Itoa(soldLatte) + " - latte")
	fmt.Println(strconv.Itoa(soldCappuccino) + " - cappuccino")
}

func buyDrink(water, milk, cofBeans, cups int) (result bool) {
	fmt.Println("Что вы хотите купить? 1 - espresso, 2 - latte, 3 - cappuccino")
	var userNumber int
	fmt.Fscan(os.Stdin, &userNumber)
	switch userNumber {
	case 1:
		if water >= 250 && cofBeans >= 16 && cups > 1 {
			ourCoffeeMachine.primaryWater = ourCoffeeMachine.primaryWater - 250
			ourCoffeeMachine.primaryCofBeans = ourCoffeeMachine.primaryCofBeans - 16
			ourCoffeeMachine.primaryMoney = ourCoffeeMachine.primaryMoney + 60
			ourCoffeeMachine.primaryCups--
			earnedMoney = earnedMoney + 60
			soldEspressoCount++
			result = true
		} else if water < 250 || cofBeans < 16 || cups < 1 {
			result = false
		}
	case 2:
		if water >= 300 && cofBeans >= 20 && milk >= 76 && cups > 1 {
			ourCoffeeMachine.primaryWater = ourCoffeeMachine.primaryWater - 300
			ourCoffeeMachine.primaryCofBeans = ourCoffeeMachine.primaryCofBeans - 20
			ourCoffeeMachine.primaryMilk = ourCoffeeMachine.primaryMilk - 76
			ourCoffeeMachine.primaryMoney = ourCoffeeMachine.primaryMoney + 110
			ourCoffeeMachine.primaryCups--
			earnedMoney = earnedMoney + 110
			soldLatteCount++
			result = true
		} else if water < 300 || cofBeans < 20 || milk < 76 || cups > 1 {
			result = false
		}
	case 3:
		if water >= 200 && cofBeans >= 16 && milk >= 100 && cups > 1 {
			ourCoffeeMachine.primaryWater = ourCoffeeMachine.primaryWater - 200
			ourCoffeeMachine.primaryCofBeans = ourCoffeeMachine.primaryCofBeans - 16
			ourCoffeeMachine.primaryMilk = ourCoffeeMachine.primaryMilk - 100
			ourCoffeeMachine.primaryMoney = ourCoffeeMachine.primaryMoney + 140
			ourCoffeeMachine.primaryCups--
			earnedMoney = earnedMoney + 140
			soldCappuccinoCount++
			result = true
		} else if water < 200 || cofBeans < 16 || milk < 100 || cups > 1 {
			result = false
		}
	}
	return result
}

func fillResources() {
	var waterNeeded, milkNeeded, cofBeansNeeded, cupsNeeded int
	fmt.Println("Сколько мл воды Вы хотите добавить?")
	fmt.Fscan(os.Stdin, &waterNeeded)
	fmt.Println("Сколько мл молока Вы хотите добавить?")
	fmt.Fscan(os.Stdin, &milkNeeded)
	fmt.Println("Сколько грамм кофейных зёрен Вы хотите добавить?")
	fmt.Fscan(os.Stdin, &cofBeansNeeded)
	fmt.Println("Сколько пустых чашек Вы хотите добавить?")
	fmt.Fscan(os.Stdin, &cupsNeeded)
	if waterNeeded+ourCoffeeMachine.primaryWater >= 5000 {
		ourCoffeeMachine.primaryWater = 5000
	} else {
		ourCoffeeMachine.primaryWater = waterNeeded + ourCoffeeMachine.primaryWater
	}
	if milkNeeded+ourCoffeeMachine.primaryMilk >= 1000 {
		ourCoffeeMachine.primaryMilk = 1000
	} else {
		ourCoffeeMachine.primaryMilk = milkNeeded + ourCoffeeMachine.primaryMilk
	}
	if cofBeansNeeded+ourCoffeeMachine.primaryCofBeans >= 900 {
		ourCoffeeMachine.primaryCofBeans = 900
	} else {
		ourCoffeeMachine.primaryCofBeans = cofBeansNeeded + ourCoffeeMachine.primaryCofBeans
	}
	if cupsNeeded+ourCoffeeMachine.primaryCups >= 50 {
		ourCoffeeMachine.primaryCups = 50
	} else {
		ourCoffeeMachine.primaryCups = cupsNeeded + ourCoffeeMachine.primaryCups
	}
}

func takeMoney() {
	fmt.Println("Выдаю Вам " + strconv.Itoa(ourCoffeeMachine.primaryMoney) + " рублей")
	ourCoffeeMachine.primaryMoney = 0
}
