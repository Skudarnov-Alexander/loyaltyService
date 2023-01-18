package luhn

func CalculateLuhn(number int) bool {
        var check int

        lastNumber := number % 10
	checkNumber := checksum(number)

        if checkNumber == 0 {
                return check == lastNumber 
        }
        
	return 10 - checkNumber == lastNumber 

}

func checksum(number int) int {
	var luhn int

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
	return luhn % 10
}