package luhn

import (
	"strconv"
)

func CheckOrder(orderID string) (bool, error) {
        num, err := strconv.Atoi(orderID)
	if err != nil {
		return false, err
	}
	
        luhn := num % 10
        num = num / 10

	for i := 0; num > 0; i++ {
		cur := num % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		num = num / 10
                
	}
        
	return luhn % 10 == 0, nil
}
