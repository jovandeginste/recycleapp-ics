package main

import (
	"fmt"
	"log"
	"time"

	ics "github.com/arran4/golang-ical"
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

func (r *RecycleInfo) ToCalendar(org string) *ics.Calendar {
	log.Printf("Items: %d", len(r.Items))

	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)

	for _, i := range r.Items {
		i.AddToCalendar(cal, org)
	}

	return cal
}

func (r *RecycleItem) AddToCalendar(cal *ics.Calendar, org string) {
	if !r.IsCollection() {
		return
	}

	event := cal.AddEvent(r.ID)
	event.SetCreatedTime(r.Fraction.CreatedAt)
	event.SetDtStampTime(r.Timestamp)
	event.SetModifiedAt(r.Fraction.UpdatedAt)
	event.SetStartAt(r.Timestamp)
	event.SetEndAt(r.Timestamp)
	event.SetSummary(r.FractionName(lang))
	event.SetOrganizer(org)
}
