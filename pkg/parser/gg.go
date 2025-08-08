package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/richard0326/poker-normalizer/pkg/model"
)

func rx(p string) *regexp.Regexp { return regexp.MustCompile(p) }

var (
	rxFlop      = rx(`\*\*\*\s*FLOP\s*\*\*\*\s*\[(.*?)\]`)
	rxTurn      = rx(`\*\*\*\s*TURN\s*\*\*\*.*\[(.*?)\]`)
	rxRiver     = rx(`\*\*\*\s*RIVER\s*\*\*\*.*\[(.*?)\]`)
	rxDealtTo   = rx(`Dealt to\s+(.+?)\s+\[(.*?)\]`)
	rxPostsSB   = rx(`^(.*)\s+posts small blind\s+\$?([\d\.]+)`)
	rxPostsBB   = rx(`^(.*)\s+posts big blind\s+\$?([\d\.]+)`)
	rxBet       = rx(`^(.*?):?\s+bet[s]?\s+\$?([\d\.]+)`)
	rxRaiseTo   = rx(`^(.*?):?\s+raise[s]?\s+to\s+\$?([\d\.]+)`)
	rxCall      = rx(`^(.*?):?\s+call[s]?\s+\$?([\d\.]+)`)
	rxCheck     = rx(`^(.*?):?\s+check[s]?`)
	rxFold      = rx(`^(.*?):?\s+fold[s]?`)
	rxCollected = rx(`^(.*)\s+collected\s+\$?([\d\.]+)\s+from pot`)
	rxTotalPot  = rx(`Total pot\s+\$?([\d\.]+)\s+\|\s+Rake\s+\$?([\d\.]+)`) 
	rxShowdown  = rx(`\*\*\*\s*SHOW ?DOWN\s*\*\*\*`)
)

type state struct {
	bbUnit       float64
	heroInvested float64
	heroCollected float64
	showdown     bool
	street       model.Street
	idx          int
	flop, turn, river string
	hero, heroCards string
	totalPot, rake float64
}

func streetFrom(line string, cur model.Street) model.Street {
	if rxFlop.MatchString(line) { return model.Flop }
	if rxTurn.MatchString(line) { return model.Turn }
	if rxRiver.MatchString(line) { return model.River }
	return cur
}

func toBB(bbUnit, v float64) float64 {
	if bbUnit <= 0 { return v }
	return v / bbUnit
}

func ParseOneHand(lines []string, hero string, handID string) (model.Hand, []model.Action) {
	st := &state{ street: model.Preflop }
	var actions []model.Action

	// 액션 추가 헬퍼
	push := func(actor, typ string, amt float64) {
		st.idx++
		a := model.Action{
			HandID: handID,
			Idx:    st.idx,
			Street: st.street,
			Actor:  strings.TrimSpace(actor),
			Type:   typ,
			SizeBB: toBB(st.bbUnit, amt),
		}
		actions = append(actions, a)
		// 히어로 투자금 누적 (증분 기준이 필요하면 상태머신으로 개선)
		if strings.EqualFold(actor, hero) && (typ == "bet" || typ == "call" || typ == "raise") {
			st.heroInvested += amt
		}
	}

	for _, raw := range lines {
		line := strings.TrimRight(raw, "\r\n")

		// 히어로 카드
		if m := rxDealtTo.FindStringSubmatch(line); m != nil {
			name := strings.TrimSpace(m[1])
			cards := strings.TrimSpace(m[2])
			if strings.EqualFold(name, hero) {
				st.hero = name
				st.heroCards = cards
			}
		}

		// 보드
		if m := rxFlop.FindStringSubmatch(line); m != nil {
			st.flop = strings.TrimSpace(m[1])
		}
		if m := rxTurn.FindStringSubmatch(line); m != nil {
			st.turn = strings.TrimSpace(m[1])
		}
		if m := rxRiver.FindStringSubmatch(line); m != nil {
			st.river = strings.TrimSpace(m[1])
		}

		// 스트리트 전이 / 쇼다운
		st.street = streetFrom(line, st.street)
		if rxShowdown.MatchString(line) {
			st.showdown = true
		}

		// 블라인드 포스트: bbUnit 세팅 + 히어로 투자 반영
		if m := rxPostsBB.FindStringSubmatch(line); m != nil {
			actor := strings.TrimSpace(m[1])
			if v, err := strconv.ParseFloat(m[2], 64); err == nil {
				if st.bbUnit == 0 {
					st.bbUnit = v
				}
				if strings.EqualFold(actor, hero) {
					st.heroInvested += v
				}
			}
		}
		if m := rxPostsSB.FindStringSubmatch(line); m != nil {
			actor := strings.TrimSpace(m[1])
			if v, err := strconv.ParseFloat(m[2], 64); err == nil {
				if strings.EqualFold(actor, hero) {
					st.heroInvested += v
				}
			}
		}

		// 액션들
		if m := rxCheck.FindStringSubmatch(line); m != nil {
			actor := strings.TrimSpace(m[1])
			push(actor, "check", 0)
		}
		if m := rxFold.FindStringSubmatch(line); m != nil {
			actor := strings.TrimSpace(m[1])
			push(actor, "fold", 0)
		}
		if m := rxBet.FindStringSubmatch(line); m != nil {
			actor := strings.TrimSpace(m[1])
			if v, err := strconv.ParseFloat(m[2], 64); err == nil {
				push(actor, "bet", v)
			}
		}
		if m := rxCall.FindStringSubmatch(line); m != nil {
			actor := strings.TrimSpace(m[1])
			if v, err := strconv.ParseFloat(m[2], 64); err == nil {
				push(actor, "call", v)
			}
		}
		if m := rxRaiseTo.FindStringSubmatch(line); m != nil {
			actor := strings.TrimSpace(m[1])
			if v, err := strconv.ParseFloat(m[2], 64); err == nil {
				// NOTE: 현재는 "to 금액" 전체를 투자로 간주. 증분으로 바꾸려면 상태추적 필요.
				push(actor, "raise", v)
			}
		}

		// 수금/팟
		if m := rxCollected.FindStringSubmatch(line); m != nil {
			actor := strings.TrimSpace(m[1])
			if v, err := strconv.ParseFloat(m[2], 64); err == nil {
				if strings.EqualFold(actor, hero) {
					st.heroCollected += v
				}
			}
		}
		if m := rxTotalPot.FindStringSubmatch(line); m != nil {
			if v, err := strconv.ParseFloat(m[1], 64); err == nil {
				st.totalPot = v
			}
			if v, err := strconv.ParseFloat(m[2], 64); err == nil {
				st.rake = v
			}
		}
	}

	// 핸드 요약
	h := model.Hand{
		HandID:         handID,
		Site:           "GG",
		Game:           "NLH",
		SBbb:           0.5, // 필요시 실제 SB 값으로 보정
		BBbb:           st.bbUnit,
		Currency:       "USD",
		TableName:      "",
		MaxPlayers:     0,
		Hero:           st.hero,
		HeroCards:      st.heroCards,
		FlopCards:      st.flop,
		TurnCard:       st.turn,
		RiverCard:      st.river,
		Showdown:       st.showdown,
		TotalPotBB:     toBB(st.bbUnit, st.totalPot),
		RakeBB:         toBB(st.bbUnit, st.rake),
		HeroCollectedB: toBB(st.bbUnit, st.heroCollected),
		HeroInvestedB:  toBB(st.bbUnit, st.heroInvested),
	}

	return h, actions
}
