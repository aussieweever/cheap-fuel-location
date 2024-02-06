package xmlmodel

import "encoding/xml"

type PositionTag struct {
	XMLName xml.Name `xml:"wpt"`
	Lat     string   `xml:"lat,attr"`
	Lon     string   `xml:"lon,attr"`
	Name    string   `xml:"name"`
}
