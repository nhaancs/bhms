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

type divisionJSON struct {
	ID       int
	ParentID int
	Code     int
	Level    uint8
	Name     string
}

func toCoreDivision(d divisionJSON) (division.Divison, error) {
	return division.Divison{
		ID:       d.ID,
		ParentID: d.ParentID,
		Code:     d.Code,
		Level:    d.Level,
		Name:     d.Name,
	}, nil
}

func toCoreDivisions(divs []divisionJSON) ([]division.Divison, error) {
	divisions := make([]division.Divison, len(divs))
	var err error
	for i := range divs {
		divisions[i], err = toCoreDivision(divs[i])
		if err != nil {
			return nil, err
		}
	}

	return divisions, nil
}
