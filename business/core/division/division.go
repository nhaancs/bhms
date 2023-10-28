package division

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/nhaancs/bhms/foundation/logger"
)

var (
	//go:embed division.json
	divJSON string
)

type province struct {
	Code      int        `json:"code"`
	Name      string     `json:"name"`
	Districts []district `json:"districts"`
}

func (p province) validate() error {
	if p.Code <= 0 {
		return fmt.Errorf("invalid province code %d", p.Code)
	}
	if len(p.Name) == 0 {
		return fmt.Errorf("province name is empty")
	}
	if len(p.Districts) == 0 {
		return fmt.Errorf("province districts is empty")
	}
	return nil
}

type district struct {
	Code  int    `json:"code"`
	Name  string `json:"name"`
	Wards []ward `json:"wards"`
}

func (d district) validate() error {
	if d.Code <= 0 {
		return fmt.Errorf("invalid district code %d", d.Code)
	}
	if len(d.Name) == 0 {
		return fmt.Errorf("district name is empty")
	}
	if len(d.Wards) == 0 {
		return fmt.Errorf("district wards is empty")
	}
	return nil
}

type ward struct {
	Code int    `json:"code"`
	Name string `json:"name"`
}

func (w ward) validate() error {
	if w.Code <= 0 {
		return fmt.Errorf("invalid ward code %d", w.Code)
	}
	if len(w.Name) == 0 {
		return fmt.Errorf("ward name is empty")
	}
	return nil
}

// =============================================================================

type Storer interface {
}

type store struct {
	Level1 []div
	Level2 map[int][]div
	Level3 map[int][]div
	Map    map[int]div
}

type div struct {
	ID       int
	ParentID int
	Code     int
	Level    uint8
	Name     string
}

// Core manages the set of APIs for user access.
type Core struct {
	store Storer
	log   *logger.Logger
}

func NewCore(log *logger.Logger) (*Core, error) {
	s, err := initStore()
	if err != nil {
		return nil, err
	}

	return &Core{
		store: s,
		log:   log,
	}, nil
}

func initStore() (*store, error) {
	var (
		id        int // autoincrement id
		provinces []province
		s         = store{
			Level2: make(map[int][]div),
			Level3: make(map[int][]div),
			Map:    make(map[int]div),
		}
	)

	// get from file
	if err := json.Unmarshal([]byte(divJSON), &provinces); err != nil {
		return nil, fmt.Errorf("provinces:unmarshal:%w", err)
	}
	if len(provinces) == 0 {
		return nil, fmt.Errorf("empty provinces")
	}

	// map to store
	for i := range provinces {
		// level 1
		if err := provinces[i].validate(); err != nil {
			return nil, err
		}

		id++
		lv1ID := id
		d := div{
			ID:       lv1ID,
			Name:     provinces[i].Name,
			Level:    1,
			Code:     provinces[i].Code,
			ParentID: 0,
		}
		s.Level1 = append(s.Level1, d)
		s.Map[lv1ID] = d

		// level 2
		for j := range provinces[i].Districts {
			if err := provinces[i].Districts[j].validate(); err != nil {
				return nil, err
			}

			id++
			lv2ID := id
			di := div{
				ID:       lv2ID,
				Name:     provinces[i].Districts[j].Name,
				Level:    2,
				Code:     provinces[i].Districts[j].Code,
				ParentID: lv1ID,
			}
			s.Level2[lv1ID] = append(s.Level2[lv1ID], di)
			s.Map[lv2ID] = di

			// level 3
			for k := range provinces[i].Districts[j].Wards {
				if err := provinces[i].Districts[j].Wards[k].validate(); err != nil {
					return nil, err
				}

				id++
				lv3ID := id
				war := div{
					ID:       lv3ID,
					Name:     provinces[i].Districts[j].Wards[k].Name,
					Level:    3,
					Code:     provinces[i].Districts[j].Wards[k].Code,
					ParentID: lv2ID,
				}
				s.Level3[lv2ID] = append(s.Level3[lv2ID], war)
				s.Map[lv3ID] = war
			}
		}
	}
	return &s, nil
}
