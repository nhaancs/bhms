// Package divisionjson ...
package divisionjson

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/nhaancs/bhms/business/core/division"
	"github.com/nhaancs/bhms/foundation/logger"
)

var (
	//go:embed division.json
	divJSON string
)

// =============================================================================

type Store struct {
	log    *logger.Logger
	level1 []divisionJSON         // divisions level 1
	level2 map[int][]divisionJSON // lv1ID => division lv2
	level3 map[int][]divisionJSON // lv2ID => division lv3
	allMap map[int]divisionJSON   // ID => division
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger) (*Store, error) {
	var (
		id        int // autoincrement id
		provinces []province
		s         = Store{
			log:    log,
			level2: make(map[int][]divisionJSON),
			level3: make(map[int][]divisionJSON),
			allMap: make(map[int]divisionJSON),
		}
	)

	// get from file
	if err := json.Unmarshal([]byte(divJSON), &provinces); err != nil {
		return nil, fmt.Errorf("provinces:unmarshal:%w", err)
	}
	if len(provinces) == 0 {
		return nil, fmt.Errorf("empty provinces")
	}

	// convert json data to store
	s.level1 = make([]divisionJSON, len(provinces))
	for i := range provinces {
		// level 1
		if err := provinces[i].validate(); err != nil {
			return nil, err
		}

		id++
		lv1ID := id
		d := divisionJSON{
			ID:       lv1ID,
			Name:     provinces[i].Name,
			Level:    1,
			Code:     provinces[i].Code,
			ParentID: 0,
		}
		s.level1[i] = d
		s.allMap[lv1ID] = d

		// level 2
		s.level2[lv1ID] = make([]divisionJSON, len(provinces[i].Districts))
		for j := range provinces[i].Districts {
			if err := provinces[i].Districts[j].validate(); err != nil {
				return nil, err
			}

			id++
			lv2ID := id
			di := divisionJSON{
				ID:       lv2ID,
				Name:     provinces[i].Districts[j].Name,
				Level:    2,
				Code:     provinces[i].Districts[j].Code,
				ParentID: lv1ID,
			}
			s.level2[lv1ID][j] = di
			s.allMap[lv2ID] = di

			// level 3
			s.level3[lv2ID] = make([]divisionJSON, len(provinces[i].Districts[j].Wards))
			for k := range provinces[i].Districts[j].Wards {
				if err := provinces[i].Districts[j].Wards[k].validate(); err != nil {
					return nil, err
				}

				id++
				lv3ID := id
				war := divisionJSON{
					ID:       lv3ID,
					Name:     provinces[i].Districts[j].Wards[k].Name,
					Level:    3,
					Code:     provinces[i].Districts[j].Wards[k].Code,
					ParentID: lv2ID,
				}
				s.level3[lv2ID][k] = war
				s.allMap[lv3ID] = war
			}
		}
	}
	return &s, nil
}

func (s *Store) QueryByID(ctx context.Context, id int) (division.Division, error) {
	div, exist := s.allMap[id]
	if !exist {
		return division.Division{}, fmt.Errorf("querybyid: %w", division.ErrNotFound)
	}
	return toCoreDivision(div)
}

func (s *Store) QueryByParentID(ctx context.Context, id int) ([]division.Division, error) {
	parent, err := s.QueryByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("querybyparentid parent: %w", division.ErrNotFound)
	}
	switch parent.Level {
	case 2:
		divisions, exist := s.level2[id]
		if !exist {
			return nil, fmt.Errorf("querybyparentid level 2: %w", division.ErrNotFound)
		}
		return toCoreDivisions(divisions)
	case 3:
		divisions, exist := s.level3[id]
		if !exist {
			return nil, fmt.Errorf("querybyparentid level 3: %w", division.ErrNotFound)
		}
		return toCoreDivisions(divisions)
	default:
		return nil, fmt.Errorf("invalid division level: %d", parent.Level)
	}
}

func (s *Store) QueryLevel1s(ctx context.Context) ([]division.Division, error) {
	return toCoreDivisions(s.level1)
}
