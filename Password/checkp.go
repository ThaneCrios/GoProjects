package main

import "fmt"

func СheckPassword(pass string) (ans bool) {
	if len(pass) < 6 {
		fmt.Println("Пароль ненадёжен. Минимальная длина не менее 6 символов.")
		return
	}

	var checkUpper, checkLower, checkDigits, checkSpecial bool
	for _, i := range pass {
		if i >= 'a' && i <= 'z' {
			checkLower = true
		} else {
			if i >= 'A' && i <= 'Z' {
				checkUpper = true
			}
		}
		if i >= '0' && i <= '9' {
			checkDigits = true
		}
		if i >= '!' && i <= '.' {
			checkSpecial = true
		}
	}

	if !checkLower {
		fmt.Println("Пароль ненадёжен. Добавьте строчные символы.")
	}
	if !checkUpper {
		fmt.Println("Пароль ненадёжен. Добавьте заглавные символы.")
	}
	if !checkDigits {
		fmt.Println("Пароль ненадёжен. Добавьте цифры.")
	}
	if !checkSpecial {
		fmt.Println("Пароль ненадёжен. Добавьте спец. символы.")
	}
	if checkUpper && checkLower && checkDigits && checkSpecial {
		ans = true
	}
	return ans
}
