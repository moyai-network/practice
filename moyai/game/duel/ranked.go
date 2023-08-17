package duel

import "math"

// eloEarnings returns the elo earned using two provided elos.
func eloEarnings(eloOne, eloTwo int32) int32 {
	increase := int32(10)
	if eloOne < 1000 {
		increase = 23
	} else if eloOne >= 1000 && eloOne < 1100 {
		increase = 17
	} else if eloOne >= 1100 && eloOne < 1200 {
		increase = 16
	} else if eloOne >= 1200 && eloOne < 1300 {
		increase = 14
	} else if eloOne >= 1300 && eloOne < 1400 {
		increase = 13
	} else if eloOne >= 1400 && eloOne < 1500 {
		increase = 12
	} else if eloOne >= 1500 && eloOne < 1600 {
		increase = 11
	} else if eloOne >= 1600 && eloOne < 1700 {
		increase = 10
	} else if eloOne >= 1700 && eloOne < 1800 {
		increase = 9
	} else if eloOne >= 1800 && eloOne < 1900 {
		increase = 8
	} else if eloOne >= 1900 && eloOne < 2000 {
		increase = 7
	} else if eloOne >= 2000 {
		increase = 6
	}

	difference := math.Abs(float64(eloOne) - float64(eloTwo))
	if eloOne < eloTwo {
		if difference >= 50 && difference < 100 {
			increase += 2
		} else if difference >= 100 && difference < 150 {
			increase += 4
		} else if difference >= 150 && difference < 200 {
			increase += 6
		} else if difference >= 200 && difference < 250 {
			increase += 8
		} else if difference >= 250 && difference < 300 {
			increase += 10
		} else if difference >= 300 {
			increase += 12
		}
	} else if eloOne > eloTwo {
		if difference >= 50 {
			increase -= 4
		} else if difference >= 100 && difference < 150 {
			increase -= 6
		} else if difference >= 150 && difference < 200 {
			increase -= 8
		} else if difference >= 200 && difference < 250 {
			increase -= 10
		} else if difference >= 250 && difference < 300 {
			increase -= 12
		} else if difference >= 300 {
			increase -= 14
		}
	}
	if increase <= 0 {
		return 1
	} else if increase > 30 {
		return 30
	}
	return increase
}

// eloLosings returns the elo loss using two provided elos.
func eloLosings(eloOne, eloTwo int32) int32 {
	decrease := int32(10)
	if eloOne < 1000 {
		decrease = 7
	} else if eloOne >= 1000 && eloOne < 1200 {
		decrease = 17
	} else if eloOne >= 1200 && eloOne < 1400 {
		decrease = 18
	} else if eloOne >= 1400 && eloOne < 1600 {
		decrease = 19
	} else if eloOne >= 1600 && eloOne < 1800 {
		decrease = 20
	} else if eloOne >= 1800 && eloOne < 2000 {
		decrease = 21
	} else if eloOne >= 2000 && eloOne < 2200 {
		decrease = 22
	} else if eloOne >= 2200 {
		decrease = 25
	}

	difference := math.Abs(float64(eloOne) - float64(eloTwo))
	if eloOne < eloTwo {
		if difference >= 50 && difference < 100 {
			decrease -= 2
		} else if difference >= 100 && difference < 150 {
			decrease -= 4
		} else if difference >= 150 && difference < 200 {
			decrease -= 6
		} else if difference >= 200 && difference < 250 {
			decrease -= 8
		} else if difference >= 250 {
			decrease -= 10
		}
	} else if eloOne > eloTwo {
		if difference > 50 && difference < 100 {
			decrease += 2
		} else if difference >= 100 && difference < 150 {
			decrease += 4
		} else if difference >= 150 && difference < 200 {
			decrease += 6
		} else if difference >= 200 && difference < 250 {
			decrease += 8
		} else if difference >= 250 {
			decrease += 10
		}
	}
	if decrease <= 0 {
		return 1
	} else if decrease > 30 {
		return 30
	}
	return decrease
}
