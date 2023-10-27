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

type divFromJSON struct {
	Code      uint32 `json:"code"`
	Name      string `json:"name"`
	Districts []struct {
		Code  uint32 `json:"code"`
		Name  string `json:"name"`
		Wards []struct {
			Code uint32 `json:"code"`
			Name string `json:"name"`
		} `json:"wards"`
	} `json:"districts"`
}

// =============================================================================

type Storer interface {
}

type store struct {
	Level1 map[uint32]div
	Level2 map[uint32]div
	Level3 map[uint32]div
}

type div struct {
	ID       uint32
	ParentID uint32
	Code     uint32
	Level    uint8
	Name     string
}

// Core manages the set of APIs for user access.
type Core struct {
	store Storer
	log   *logger.Logger
}

func NewCore(log *logger.Logger) (*Core, error) {
	var s = store{
		Level1: make(map[uint32]div),
		Level2: make(map[uint32]div),
		Level3: make(map[uint32]div),
	}

	var divData []divFromJSON
	if err := json.Unmarshal([]byte(divJSON), &divData); err != nil {
		return nil, fmt.Errorf("division:init:unmarshal:%w", err)
	}
	if len(divData) == 0 {
		return nil, fmt.Errorf("division:init:empty data")
	}

	return &Core{
		store: s,
		log:   log,
	}, nil
}
