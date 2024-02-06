package xmlmodel

import "encoding/xml"

type GpxTag struct {
	XMLName  xml.Name    `xml:"gpx"`
	Version  string      `xml:"version,attr"`
	Creator  string      `xml:"creator,attr"`
	Position PositionTag `xml:"wpt"`
}
