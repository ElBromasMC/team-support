package book

import "time"

type DocumentType string

const (
	DNIType    DocumentType = "DNI"
	CarnetType DocumentType = "CARNET"
	OtherType  DocumentType = "OTHER"
)

type GoodType string

const (
	ProductType GoodType = "PRODUCT"
	ServiceType GoodType = "SERVICE"
)

type ComplaintType string

const (
	ReclamoType ComplaintType = "RECLAMO"
	QuejaType   ComplaintType = "QUEJA"
)

type Entry struct {
	Id        int
	CreatedAt time.Time
	// Complaining Customer
	DocumentType   DocumentType
	DocumentNumber string
	Name           string
	Address        string
	PhoneNumber    string
	Email          string
	ParentName     string
	// Contracted Good
	GoodType        GoodType
	GoodDescription string
	// Complaint Details
	ComplaintType        ComplaintType
	ComplaintDescription string
	// Actions
	ActionsDescription string
}

func (dt DocumentType) ToSlug() string {
	if dt == DNIType {
		return "DNI"
	} else if dt == CarnetType {
		return "Carne"
	} else if dt == OtherType {
		return "Otro"
	}
	return ""
}

func (gt GoodType) ToSlug() string {
	if gt == ProductType {
		return "Producto"
	} else if gt == ServiceType {
		return "Servicio"
	}
	return ""
}

func (ct ComplaintType) ToSlug() string {
	if ct == ReclamoType {
		return "Reclamo"
	} else if ct == QuejaType {
		return "Queja"
	}
	return ""
}
