package model

type Street string

const (
	Preflop Street = "preflop"
	Flop    Street = "flop"
	Turn    Street = "turn"
	River   Street = "river"
)

type Action struct {
	HandID   string  `csv:"handId" json:"handId"`
	Idx      int     `csv:"idx" json:"idx"`
	Street   Street  `csv:"street" json:"street"`
	Actor    string  `csv:"actor" json:"actor"`
	Type     string  `csv:"type" json:"type"`
	SizeBB   float64 `csv:"sizeBB" json:"sizeBB"`
	PotBB    float64 `csv:"potBB" json:"potBB"`
	ToCallBB float64 `csv:"toCallBB" json:"toCallBB"`
	SPR      float64 `csv:"spr" json:"spr"`
}

type Hand struct {
	HandID         string  `csv:"handId" json:"handId"`
	Site           string  `csv:"site" json:"site"`
	Game           string  `csv:"game" json:"game"`
	SBbb           float64 `csv:"sbBB" json:"sbBB"`
	BBbb           float64 `csv:"bbBB" json:"bbBB"`
	Currency       string  `csv:"currency" json:"currency"`
	StartedAt      string  `csv:"startedAt" json:"startedAt"`
	TableName      string  `csv:"table" json:"table"`
	MaxPlayers     int     `csv:"maxPlayers" json:"maxPlayers"`
	Hero           string  `csv:"hero" json:"hero"`
	HeroCards      string  `csv:"heroCards" json:"heroCards"`
	FlopCards      string  `csv:"flop" json:"flop"`
	TurnCard       string  `csv:"turn" json:"turn"`
	RiverCard      string  `csv:"river" json:"river"`
	Showdown       bool    `csv:"showdown" json:"showdown"`
	TotalPotBB     float64 `csv:"totalPotBB" json:"totalPotBB"`
	RakeBB         float64 `csv:"rakeBB" json:"rakeBB"`
	HeroCollectedB float64 `csv:"heroCollectedBB" json:"heroCollectedBB"`
	HeroInvestedB  float64 `csv:"heroInvestedBB" json:"heroInvestedBB"`
	RedlineBB      float64 `csv:"redlineBB" json:"redlineBB"`
	BluelineBB     float64 `csv:"bluelineBB" json:"bluelineBB"`
	RealProfitB    float64 `csv:"realProfitBB" json:"realProfitBB"`
	EVProfitB      float64 `csv:"evProfitBB" json:"evProfitBB"`
	BoardType      string  `csv:"boardType" json:"boardType"`
	SPRFlop        float64 `csv:"sprFlop" json:"sprFlop"`
}