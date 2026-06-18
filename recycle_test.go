package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/jordic/goics"
)

type mockEmitter struct {
	comp goics.Componenter
}

func (m mockEmitter) EmitICal() goics.Componenter {
	return m.comp
}

func TestRecycleItemIsCollection(t *testing.T) {
	tests := []struct {
		item RecycleItem
		want bool
	}{
		{
			item: RecycleItem{
				Type: "collection",
			},
			want: true,
		},
		{
			item: RecycleItem{
				Type: "not-collection",
			},
			want: false,
		},
		{
			item: RecycleItem{
				Type: "collection",
				Exception: RecycleException{
					ReplacedBy: struct {
						Type string `json:"type"`
					}{Type: "collection"},
				},
			},
			want: false,
		},
		{
			item: RecycleItem{
				Type: "collection",
				Exception: RecycleException{
					Replaces: struct {
						Type string `json:"type"`
					}{Type: "collection"},
				},
			},
			want: true,
		},
	}

	for _, tc := range tests {
		got := tc.item.IsCollection()
		if got != tc.want {
			t.Errorf("IsCollection for type %s / replaces %s / replacedBy %s = %t; want %t",
				tc.item.Type, tc.item.Exception.Replaces.Type, tc.item.Exception.ReplacedBy.Type, got, tc.want)
		}
	}
}

func TestFractionName(t *testing.T) {
	item := RecycleItem{
		Fraction: RecycleFraction{
			Name: map[string]string{
				"nl": "Plastic",
				"fr": "Plastique",
			},
		},
	}

	if got := item.FractionName("nl"); got != "Plastic" {
		t.Errorf("expected Plastic, got %s", got)
	}

	if got := item.FractionName("de"); got != "???" {
		t.Errorf("expected ???, got %s", got)
	}
}

func TestToEvent(t *testing.T) {
	oldLang := lang
	lang = "nl"

	defer func() { lang = oldLang }()

	org := &Organization{
		Name: "TestOrg",
	}

	now := time.Date(2026, 6, 18, 12, 0, 0, 0, time.UTC)

	item := &RecycleItem{
		ID:        "item-123",
		Type:      "collection",
		Timestamp: now,
		Fraction: RecycleFraction{
			Name:      map[string]string{"nl": "Papier"},
			Color:     "blue",
			CreatedAt: now.Add(-1 * time.Hour),
			UpdatedAt: now,
		},
	}

	evt := item.ToEvent(org)
	if evt == nil {
		t.Fatal("expected non-nil event")
	}

	b := bytes.Buffer{}
	goics.NewICalEncode(&b).Encode(mockEmitter{comp: evt})

	gotStr := b.String()

	expectedSubstrings := []string{
		"BEGIN:VEVENT",
		"UID:item-123",
		"COLOR:blue",
		"SUMMARY:Papier",
		"ORGANIZER;CN=TESTORG:nomail",
		"TZID:Europe/Brussels",
		"BEGIN:VALARM",
		"TRIGGER;RELATED=START:-PT6H",
		"ACTION:DISPLAY",
		"DESCRIPTION:Papier",
		"END:VALARM",
		"END:VEVENT",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(gotStr, sub) {
			t.Errorf("expected generated ICS to contain %q, but it didn't. Got:\n%s", sub, gotStr)
		}
	}

	nonCollection := &RecycleItem{
		Type: "other",
	}
	if evtNonCol := nonCollection.ToEvent(org); evtNonCol != nil {
		t.Errorf("expected nil event for non-collection, got %v", evtNonCol)
	}
}

func TestEmitICal(t *testing.T) {
	oldLang := lang
	lang = "nl"

	defer func() { lang = oldLang }()

	org := &Organization{
		Name: "TestOrg",
	}

	now := time.Date(2026, 6, 18, 12, 0, 0, 0, time.UTC)

	info := RecycleInfo{
		Org: org,
		Items: []RecycleItem{
			{
				ID:        "item-abc",
				Type:      "collection",
				Timestamp: now,
				Fraction: RecycleFraction{
					Name:      map[string]string{"nl": "Glas"},
					Color:     "green",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
	}

	b := bytes.Buffer{}
	goics.NewICalEncode(&b).Encode(info)

	gotStr := b.String()

	expectedSubstrings := []string{
		"BEGIN:VCALENDAR",
		"VERSION:2.0",
		"PRODID:-//Recycling calendar",
		"METHOD:REQUEST",
		"BEGIN:VEVENT",
		"UID:item-abc",
		"SUMMARY:Glas",
		"END:VEVENT",
		"END:VCALENDAR",
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(gotStr, sub) {
			t.Errorf("expected calendar ICS to contain %q, but it didn't. Got:\n%s", sub, gotStr)
		}
	}
}

func TestRecycleInfoJSON(t *testing.T) {
	oldLang := lang
	lang = "nl"
	defer func() { lang = oldLang }()

	now := time.Date(2026, 6, 18, 12, 0, 0, 0, time.UTC)
	info := RecycleInfo{
		Org: &Organization{
			Name: "TestOrg",
		},
		Items: []RecycleItem{
			{
				ID:        "item-abc",
				Type:      "collection",
				Timestamp: now,
				Fraction: RecycleFraction{
					Name:      map[string]string{"nl": "Glas"},
					Color:     "green",
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
		},
	}

	jsonData, err := json.Marshal(info.ToJSONEvents())
	if err != nil {
		t.Fatalf("failed to marshal JSONEvents to JSON: %v", err)
	}

	gotStr := string(jsonData)
	expectedSubstrings := []string{
		`"summary":"Glas"`,
		`"date":"2026-06-18"`,
		`"color":"green"`,
	}

	for _, sub := range expectedSubstrings {
		if !strings.Contains(gotStr, sub) {
			t.Errorf("expected JSON to contain %q, but it didn't. Got:\n%s", sub, gotStr)
		}
	}

	unexpectedSubstrings := []string{
		`"uid"`,
		`"organizer"`,
		`"created"`,
		`"last_modified"`,
	}

	for _, sub := range unexpectedSubstrings {
		if strings.Contains(gotStr, sub) {
			t.Errorf("expected JSON to NOT contain %q, but it did. Got:\n%s", sub, gotStr)
		}
	}
}
