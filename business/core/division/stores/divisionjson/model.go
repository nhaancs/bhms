package divisionjson

import (
	"fmt"
	"github.com/nhaancs/bhms/business/core/division"
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

type divisionJSON struct {
	ID       int
	ParentID int
	Code     int
	Level    uint8
	Name     string
}

func toCoreDivision(o divisionJSON) (division.Division, error) {
	return division.Division{
		ID:       o.ID,
		ParentID: o.ParentID,
		Code:     o.Code,
		Level:    o.Level,
		Name:     o.Name,
	}, nil
}

func toCoreDivisions(os []divisionJSON) ([]division.Division, error) {
	divisions := make([]division.Division, len(os))
	var err error
	for i := range os {
		divisions[i], err = toCoreDivision(os[i])
		if err != nil {
			return nil, err
		}
	}

	return divisions, nil
}
