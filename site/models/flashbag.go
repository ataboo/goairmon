package models

import "encoding/json"

type FlashBag struct {
	Success []string `json:"success"`
	Info    []string `json:"info"`
	Error   []string `json:"error"`
}

func (f *FlashBag) Decode(strVal string) error {
	return json.Unmarshal([]byte(strVal), f)
}

func (f *FlashBag) Encode() string {
	bytes, _ := json.Marshal(f)

	return string(bytes)
}

func (f *FlashBag) HasSuccess() bool {
	return len(f.Success) > 0
}

func (f *FlashBag) HasInfo() bool {
	return len(f.Info) > 0
}

func (f *FlashBag) HasError() bool {
	return len(f.Error) > 0
}
