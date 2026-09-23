package constants

type PlayerPosition struct {
	Code  string
	Title string
}

type MapPlayerPosition map[string]*PlayerPosition

func (m MapPlayerPosition) GetValue(code string) *PlayerPosition {
	val, ok := m[code]
	if !ok {
		return nil
	}

	return val
}

var DictPlayerPosition = make(MapPlayerPosition)

func registerPlayerPosition(code, title string) *PlayerPosition {
	pp := &PlayerPosition{
		Code:  code,
		Title: title,
	}

	DictPlayerPosition[code] = pp
	return pp
}

var (
	GoalKeeperPosition = registerPlayerPosition("K", "Penjaga Gawang")
	BackPosition       = registerPlayerPosition("B", "Pemain Bertahan")
	MidPosition        = registerPlayerPosition("M", "Pemain Tengah")
	StrikerPosition    = registerPlayerPosition("S", "Pemain Penyerang")
)
