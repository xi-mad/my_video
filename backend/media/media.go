package media

import "encoding/xml"

type NFO struct {
	XMLName xml.Name `xml:"movie"`
	Plot    Inner    `xml:"plot"`
	Title   Inner    `xml:"title"`
	Number  string   `xml:"num"`
	Release string   `xml:"release"`
	Label   string   `xml:"label"`
	Actor   []Actor  `xml:"actor"`
	Tag     []Inner  `xml:"tag"`
}

type Inner struct {
	Inner string `xml:",innerxml"`
}

type Actor struct {
	Name  string `xml:"name"`
	Thumb string `xml:"thumb"`
}
