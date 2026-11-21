package pages

import (
	"testing"

	hpdata "mljr-web/projects/homepage/data"
)

func TestMergeEarlyInternships(t *testing.T) {
	end := func(s string) *string { return &s }
	items := []hpdata.TimelineItem{
		{ID: "jku-tutor", Start: "2023-10"},
		{ID: "bosch-rexroth-intern", Organization: "Bosch Rexroth", Start: "2021-06", End: end("2021-06")},
		{ID: "herzreha-civil", Start: "2020-09"},
		{ID: "engel-intern-2018", Organization: "ENGEL", Start: "2018-07", End: end("2018-07")},
		{ID: "engel-intern-2017", Organization: "ENGEL", Start: "2017-07", End: end("2017-07")},
	}
	got := mergeEarlyInternships(items)
	if len(got) != 3 {
		t.Fatalf("got %d entries, want 3", len(got))
	}
	m := got[1]
	if m.ID != "early-internships" || m.Start != "2017" || *m.End != "2021" || m.Organization != "ENGEL · Bosch Rexroth" {
		t.Errorf("merged entry = %+v", m)
	}
	if got[0].ID != "jku-tutor" || got[2].ID != "herzreha-civil" {
		t.Errorf("order = %s, %s", got[0].ID, got[2].ID)
	}
}
