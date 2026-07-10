package user

type User struct {
	Name  string
	Level int32
	Stats []UserStats
	Xp    int64
	Gold  int64
	Rank  string
}

type UserStats struct {
	Name  string
	Value int32
}
