package dto

import (
	"encoding/xml"
)

type PDVListe struct {
	XMLName xml.Name `xml:"pdv_liste"`
	PDVs    []PDV    `xml:"pdv"`
}

type PDV struct {
	ID        string  `xml:"id,attr"`
	Latitude  float64 `xml:"latitude,attr"`
	Longitude float64 `xml:"longitude,attr"`
	CP        string  `xml:"cp,attr"`
	Pop       string  `xml:"pop,attr"`

	Adresse string `xml:"adresse"`
	Ville   string `xml:"ville"`

	Horaires Horaires `xml:"horaires"`
	Services Services `xml:"services"`
	Prix     []Prix   `xml:"prix"`
}

type Horaires struct {
	Automate24 string `xml:"automate-24-24,attr"`
	Jours      []Jour `xml:"jour"`
}

type Jour struct {
	ID    int    `xml:"id,attr"`
	Nom   string `xml:"nom,attr"`
	Ferme string `xml:"ferme,attr"`

	Horaires []Horaire `xml:"horaire"`
}

type Horaire struct {
	Ouverture string `xml:"ouverture,attr"`
	Fermeture string `xml:"fermeture,attr"`
}

type Services struct {
	List []string `xml:"service"`
}

type Prix struct {
	Nom    string  `xml:"nom,attr"`
	ID     int     `xml:"id,attr"`
	Maj    string  `xml:"maj,attr"`
	Valeur float64 `xml:"valeur,attr"`
}
