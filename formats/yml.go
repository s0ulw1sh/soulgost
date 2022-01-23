package formats

import (
	"strings"
	"github.com/s0ulw1sh/soulgost/hash"
)

type YmlCategory struct {
	Id    string `xml:"id,attr"`
	Pid   string `xml:"parentId,attr"`
	Name  string `xml:",chardata"`
}

type YmlOffer struct {
	Id        string `xml:"id,attr"`
	Available string `xml:"available,attr"`
	Category  string `xml:"categoryId"`
	Name      string `xml:"name"`
	Model     string `xml:"model"`
	Url       string `xml:"url"`
	Price     string `xml:"price"`
	Pic       string `xml:"picture"`
	VenName   string `xml:"vendor"`
	VenCode   string `xml:"vendorCode"`
	Desc      string `xml:"description"`
}

func (self *YmlOffer) IdHash() uint32 {
	return hash.MurMur([]byte(self.Id))
}

func (self *TYmlOffer) FullName() string {
	if self.Name != "" {
		return self.Name
	}

	if self.Model != "" {
		return self.Model
	}

	return ""
}

func (self *TYmlOffer) Prepare() bool {
	self.Id    = strings.TrimSpace(self.Id)
	self.Price = strings.TrimSpace(self.Price)

	return len(self.Id) > 0 && len(self.Price) > 0 && (len(self.Name) > 0 || len(self.Model) > 0)
}

func (self *TYmlOffer) CompositeHash() (hash string) {
	hash += hash.CRC16Hex([]byte(self.Id))
	hash += hash.CRC16Hex([]byte(self.Category))
	hash += hash.CRC16Hex([]byte(self.FullName()))
	hash += hash.CRC16Hex([]byte(self.Url))
	hash += hash.CRC16Hex([]byte(self.Desc))

	return
}