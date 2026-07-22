package stats

const IntelligenceCode = "INT"
const StrengthCode = "STR"
const EnduranceCode = "END"
const WisdomCode = "WIS"

func LvlEncrease(lvl int, xp int) (int32, int64) {
	nextLvl := lvl * 100
	if xp >= nextLvl {
		updxp := xp - nextLvl
		return int32(lvl + 1), int64(updxp)
	}
	return int32(lvl), int64(xp)
}
