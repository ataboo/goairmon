package models

type FlashBag struct {
	Success []string
	Info    []string
	Error   []string
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
