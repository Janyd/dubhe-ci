package common

import (
	"strings"
)

type Page struct {
	Total   int64       `json:"total"`
	Current uint32      `json:"current"`
	Size    uint32      `json:"size"`
	Records interface{} `json:"records"`
	Desc    []string    `json:"desc"`
	Asc     []string    `json:"asc"`
	bean    interface{} `json:"-"`
	columns []string    `json:"-"`
}

func NewPageSort(current, size uint32, desc []string, asc []string, records interface{}) *Page {
	return &Page{Current: current, Size: size, Desc: desc, Asc: asc, Records: records}
}

func NewPage(current, size uint32, out interface{}) *Page {
	return &Page{
		Current: current,
		Size:    size,
		Records: out,
	}
}

func (p *Page) GetCurrent() uint32 {
	if p.Current <= 0 {
		return 1
	}
	return p.Current
}

func (p *Page) GetSize() uint32 {
	if p.Size <= 0 || p.Size > 1000 {
		return 10
	}

	return p.Size
}

func (p *Page) Offset() uint32 {
	return (p.Current - 1) * p.Size
}

func (p *Page) SetTotal(total int64) {
	p.Total = total
}

func (p *Page) Bind(records interface{}, bean interface{}) *Page {
	p.bean = bean
	p.Records = records
	return p
}

func (p *Page) Bean(bean interface{}) *Page {
	p.bean = bean
	return p
}

func (p *Page) GetBean() interface{} {
	return p.bean
}

func (p *Page) Columns(columns ...string) {
	p.columns = columns
}

func (p *Page) GetColumns() string {
	return strings.Join(p.columns, ", ")
}
