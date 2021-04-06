package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jordic/goics"
)

const collectionType = "collection"

type RecycleInfo struct {
	Items []RecycleItem `json:"items"`
	Page  int           `json:"page"`
	Total int           `json:"total"`
	Size  int           `json:"size"`
	Self  string        `json:"self"`
	Last  string        `json:"last"`
	First string        `json:"first"`

	Org string
}

type RecycleItem struct {
	ID        string           `json:"id"`
	Type      string           `json:"type"`
	Timestamp time.Time        `json:"timestamp"`
	Fraction  RecycleFraction  `json:"fraction"`
	Exception RecycleException `json:"exception"`
}

type RecycleFraction struct {
	Name      map[string]string `json:"name"`
	Color     string            `json:"color"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

type RecycleException struct {
	ReplacedBy struct {
		Type string `json:"type"`
	} `json:"replacedBy"`
	Replaces struct {
		Type string `json:"type"`
	} `json:"replaces"`
}

func (r *RecycleItem) IsCollection() bool {
	if r.Type != collectionType {
		return false
	}

	return r.Exception.IsCollection()
}

func (r *RecycleItem) FractionName(lang string) string {
	if name, ok := r.Fraction.Name[lang]; ok {
		return name
	}

	fmt.Printf("%#v\n", r)

	return "???"
}

func (r *RecycleException) Type() string {
	if r.ReplacedBy.Type == collectionType {
		return "replaced_by"
	}

	if r.Replaces.Type == collectionType {
		return "replaces"
	}

	return "normal"
}

func (r *RecycleException) IsCollection() bool {
	return r.Type() != "replaced_by"
}

func (r RecycleInfo) EmitICal() goics.Componenter {
	log.Printf("Items: %d", len(r.Items))

	cal := goics.NewComponent()

	cal.SetType("VCALENDAR")
	cal.AddProperty("VERSION", "2.0")
	cal.AddProperty("PRODID", "-//Recycling calendar")
	cal.AddProperty("METHOD", "REQUEST")

	for _, i := range r.Items {
		e := i.ToEvent(r.Org)
		if e != nil {
			cal.AddComponent(e)
		}
	}

	return cal
}

func (r *RecycleItem) ToEvent(org string) goics.Componenter {
	if !r.IsCollection() {
		return nil
	}

	s := goics.NewComponent()
	s.SetType("VEVENT")
	s.AddProperty("UID", r.ID)

	AddDateTimeField(s, "LAST-MODIFIED", r.Fraction.UpdatedAt)
	AddDateTimeField(s, "CREATED", r.Fraction.CreatedAt)
	AddDateTimeField(s, "DTSTAMP", r.Timestamp)

	k, v := goics.FormatDateField("DTSTART", r.Timestamp)
	s.AddProperty(k, v)

	s.AddProperty("COLOR", r.Fraction.Color)
	s.AddProperty("TRANSP", "TRANSPARENT")
	s.AddProperty("SUMMARY", r.FractionName(lang))
	s.AddProperty("ORGANIZER", org)
	s.AddProperty("TZID", "Europe/Brussels")

	s.AddComponent(r.Alarm())

	return s
}

func (r *RecycleItem) Alarm() goics.Componenter {
	s := goics.NewComponent()

	s.SetType("VALARM")
	s.AddProperty("TRIGGER;RELATED=START", "-PT6H")
	s.AddProperty("ACTION", "DISPLAY")
	s.AddProperty("DESCRIPTION", r.FractionName(lang))

	return s
}

func AddDateTimeField(cal goics.Componenter, name string, t time.Time) {
	k, v := goics.FormatDateTime(name, t)
	cal.AddProperty(k, v)
}
