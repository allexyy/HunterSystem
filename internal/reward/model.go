package reward

const DifficultEasy = "Easy"
const DifficultNormal = "Normal"
const DifficultHard = "Hard"
const DifficultEpic = "Epic"

type Reward struct {
	XP   int32
	Gold int32
}

func Easy() Reward {
	return Reward{XP: 10, Gold: 5}
}

func Normal() Reward {
	return Reward{XP: 20, Gold: 10}
}

func Hard() Reward {
	return Reward{XP: 40, Gold: 25}
}

func Epic() Reward {
	return Reward{XP: 100, Gold: 60}
}

func GetRewardByDifficult(dif string) Reward {
	rewards := map[string]Reward{
		DifficultEasy:   Easy(),
		DifficultNormal: Normal(),
		DifficultHard:   Hard(),
		DifficultEpic:   Epic(),
	}
	return rewards[dif]
}
