package entities

import (
	"encoding/xml"
)

type RequestEnvelope struct {
	XMLName   xml.Name `xml:"soap:Envelope"`
	XmlnsSoap string   `xml:"xmlns:soap,attr"`

	Body struct {
		XMLName  xml.Name `xml:"soap:Body"`
		XmlnsWsu string   `xml:"xmlns:wsu,attr"`
		WsuId    string   `xml:"wsu:Id,attr"`

		SendMessage struct {
			XMLName  xml.Name `xml:"ns2:SendMessage"`
			XmlnsNs2 string   `xml:"xmlns:ns2,attr"`

			Payload interface{}
		}
	}
}

type ResponseEnvelope struct {
	XMLName xml.Name
	Body    struct {
		XMLName             xml.Name
		SendMessageResponse struct {
			XMLName  xml.Name
			Response interface{} `xml:"response"`
		}
	}
}
