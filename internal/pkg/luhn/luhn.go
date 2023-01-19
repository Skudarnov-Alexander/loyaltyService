package luhn

import "fmt"

/*
func CalculateLuhn(number int) bool {
        lastNumber := number % 10
	checkNumber := Checksum(number)

        fmt.Printf("lastNumber %d", lastNumber)
        fmt.Printf("checkNumber %d", checkNumber)

        if checkNumber == 0 {
                return checkNumber == lastNumber
        }



	return 10 - checkNumber == lastNumber

}
*/
//354835541278
func Checksum(number int, len int) int {
        var isEvenNumDigits bool 

        if len % 2 == 0 {
                isEvenNumDigits = true
        } 

        _ = isEvenNumDigits
	
        luhn := number % 10
        //lastNumber := number % 10
        number = number / 10

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
                
	}
        fmt.Printf("luhn %d\n", luhn)
	return luhn % 10
}
