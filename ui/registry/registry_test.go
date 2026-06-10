package registry

import "testing"

func TestAllAndCategoriesAreSorted(t *testing.T) {
	resetRegistry(t)

	Register(&Component{Slug: "zeta", Name: "Zeta", Category: "Primitive"})
	Register(&Component{Slug: "alpha", Name: "Alpha", Category: "Primitive"})
	Register(&Component{Slug: "banner", Name: "Banner", Category: "Layout"})

	all := All()
	gotOrder := []string{all[0].Slug, all[1].Slug, all[2].Slug}
	wantOrder := []string{"banner", "alpha", "zeta"}
	for i := range wantOrder {
		if gotOrder[i] != wantOrder[i] {
			t.Fatalf("All() order = %v, want %v", gotOrder, wantOrder)
		}
	}

	gotCategories := Categories()
	wantCategories := []string{"Layout", "Primitive"}
	for i := range wantCategories {
		if gotCategories[i] != wantCategories[i] {
			t.Fatalf("Categories() = %v, want %v", gotCategories, wantCategories)
		}
	}
}

func TestDefaultPropsAndCombinations(t *testing.T) {
	component := &Component{
		Controls: []Control{
			{Name: "variant", Type: ControlEnum, Options: []string{"default", "danger"}, Default: "default"},
			{Name: "size", Type: ControlEnum, Options: []string{"sm", "lg"}, Default: "sm"},
			{Name: "disabled", Type: ControlBool, Default: "false"},
		},
	}

	defaults := DefaultProps(component)
	if defaults["variant"] != "default" || defaults["size"] != "sm" || defaults["disabled"] != "false" {
		t.Fatalf("DefaultProps() = %v", defaults)
	}

	combos := Combinations(component)
	if len(combos) != 4 {
		t.Fatalf("len(Combinations()) = %d, want 4", len(combos))
	}
	for _, combo := range combos {
		if combo["disabled"] != "false" {
			t.Fatalf("combo lost bool default: %v", combo)
		}
		if combo["variant"] == "" || combo["size"] == "" {
			t.Fatalf("combo missing enum value: %v", combo)
		}
	}

	if got := ComboLabel(component, map[string]string{"variant": "danger", "size": "lg"}); got != "danger · lg" {
		t.Fatalf("ComboLabel() = %q, want danger · lg", got)
	}
}

func TestRegisterOverwritesSlug(t *testing.T) {
	resetRegistry(t)

	Register(&Component{Slug: "button", Name: "Old", Category: "Primitive"})
	Register(&Component{Slug: "button", Name: "New", Category: "Primitive"})

	component, ok := Get("button")
	if !ok {
		t.Fatal("Get(button) ok = false, want true")
	}
	if component.Name != "New" {
		t.Fatalf("component name = %q, want New", component.Name)
	}
}

func resetRegistry(t *testing.T) {
	t.Helper()

	mu.Lock()
	previous := reg
	reg = map[string]*Component{}
	mu.Unlock()

	t.Cleanup(func() {
		mu.Lock()
		reg = previous
		mu.Unlock()
	})
}
