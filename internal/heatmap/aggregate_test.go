package heatmap

import (
	"testing"
	"time"

	"cm_open_api/internal/models"
)

func mustTime(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestAggregate(t *testing.T) {
	rows := []SourceRow{
		// инцидент i1: два сообщения (повторный анонс) — те же дома, тот же месяц
		{Month: mustTime("2026-01-01"), Service: "water", IncidentID: "i1",
			StreetKladr: "S1", CityKladr: "C1", CityName: "Тирасполь",
			StreetName: "Ленина", StreetType: "ул", HouseNumbers: "1,2"},
		{Month: mustTime("2026-01-01"), Service: "water", IncidentID: "i1",
			StreetKladr: "S1", CityKladr: "C1", CityName: "Тирасполь",
			StreetName: "Ленина", StreetType: "ул", HouseNumbers: "1,2"},
		// инцидент i2: та же улица, дом 1 + диапазон
		{Month: mustTime("2026-01-01"), Service: "water", IncidentID: "i2",
			StreetKladr: "S1", CityKladr: "C1", CityName: "Тирасполь",
			StreetName: "Ленина", StreetType: "ул", HouseNumbers: "1", HouseRanges: "10-11"},
		// инцидент i3: без домов вообще → street_incidents
		{Month: mustTime("2026-01-01"), Service: "water", IncidentID: "i3",
			StreetKladr: "S1", CityKladr: "C1", CityName: "Тирасполь",
			StreetName: "Ленина", StreetType: "ул"},
		// инцидент i4: только неразворачиваемый диапазон → street_incidents
		{Month: mustTime("2026-01-01"), Service: "water", IncidentID: "i4",
			StreetKladr: "S1", CityKladr: "C1", CityName: "Тирасполь",
			StreetName: "Ленина", StreetType: "ул", HouseRanges: "12а-12г"},
		// другая услуга — отдельная строка
		{Month: mustTime("2026-01-01"), Service: "electricity", IncidentID: "i5",
			StreetKladr: "S1", CityKladr: "C1", CityName: "Тирасполь",
			StreetName: "Ленина", StreetType: "ул", HouseNumbers: "1"},
		// другой месяц — отдельная строка
		{Month: mustTime("2026-02-01"), Service: "water", IncidentID: "i6",
			StreetKladr: "S2", CityKladr: "C1", CityName: "Тирасполь",
			StreetName: "Мира", StreetType: "пер", HouseNumbers: "3"},
		// инцидент i7: две строки — с домом и без; при наличии хотя бы одного
		// разрешимого дома в группе НЕ должен попасть в street_incidents
		{Month: mustTime("2026-02-01"), Service: "water", IncidentID: "i7",
			StreetKladr: "S2", CityKladr: "C1", CityName: "Тирасполь",
			StreetName: "Мира", StreetType: "пер", HouseNumbers: "5"},
		{Month: mustTime("2026-02-01"), Service: "water", IncidentID: "i7",
			StreetKladr: "S2", CityKladr: "C1", CityName: "Тирасполь",
			StreetName: "Мира", StreetType: "пер"},
	}

	got, total := Aggregate(rows)

	if total != 7 {
		t.Errorf("total incidents = %d, want 7", total)
	}
	want := []models.HeatmapRow{
		{Month: "2026-01", Service: "electricity", StreetKladr: "S1", CityKladr: "C1",
			CityName: "Тирасполь", StreetName: "Ленина", StreetType: "ул",
			HouseIncidents: map[string]int{"1": 1}},
		{Month: "2026-01", Service: "water", StreetKladr: "S1", CityKladr: "C1",
			CityName: "Тирасполь", StreetName: "Ленина", StreetType: "ул",
			HouseIncidents:  map[string]int{"1": 2, "2": 1, "10": 1, "11": 1},
			StreetIncidents: 2},
		{Month: "2026-02", Service: "water", StreetKladr: "S2", CityKladr: "C1",
			CityName: "Тирасполь", StreetName: "Мира", StreetType: "пер",
			HouseIncidents: map[string]int{"3": 1, "5": 1}},
	}
	if len(got) != len(want) {
		t.Fatalf("rows = %d, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		w, g := want[i], got[i]
		if g.Month != w.Month || g.Service != w.Service || g.StreetKladr != w.StreetKladr ||
			g.StreetIncidents != w.StreetIncidents || g.CityName != w.CityName ||
			g.StreetName != w.StreetName || g.StreetType != w.StreetType || g.CityKladr != w.CityKladr {
			t.Errorf("row %d = %+v, want %+v", i, g, w)
		}
		if len(g.HouseIncidents) != len(w.HouseIncidents) {
			t.Errorf("row %d houses = %v, want %v", i, g.HouseIncidents, w.HouseIncidents)
			continue
		}
		for h, c := range w.HouseIncidents {
			if g.HouseIncidents[h] != c {
				t.Errorf("row %d house %q = %d, want %d", i, h, g.HouseIncidents[h], c)
			}
		}
	}
}
