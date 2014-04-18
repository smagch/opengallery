package main

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isHex(r rune) bool {
	return isDigit(r) || ('a' <= r && 'f' >= r) || ('A' <= r && 'F' >= r)
}

// 8, 13, 23, 18 is '-'
// 6ba7b814-9dad-11d1-80b4-00c04fd430c8
func IsUUID(s string) bool {
	if len(s) != 36 {
		return false
	}
	var count int
	for i, r := range s {
		// should not be multi byte char
		if count != i {
			return false
		}
		switch count {
		case 8, 13, 18, 23:
			if r != '-' {
				return false
			}
		default:
			if !isHex(r) {
				return false
			}
		}
		count += 1
	}
	return true
}
