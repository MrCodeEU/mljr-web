package pages

type orgLocation struct {
	lat float64
	lng float64
}

// orgLocation2 maps organization names (EN or DE) to map coordinates.
func orgLocation2(org string) (orgLocation, bool) {
	m := map[string]orgLocation{
		"Dynatrace":                        {48.314822011928605, 14.305210279591181},
		"Land Oberösterreich":              {48.29144053833716, 14.289947711349612},
		"Johannes Kepler University Linz":  {48.3371, 14.3196},
		"Johannes Kepler Universität Linz": {48.3371, 14.3196},
		"ventopay gmbh":                    {48.3678, 14.5165},
		"Bosch":                            {48.307830506435174, 14.308910298429641},
		"Bosch Rexroth":                    {48.24768759781733, 14.196875921934444},
		"ENGEL":                            {48.15757939951312, 14.505377793376999},
		"HerzReha Bad Ischl":               {47.7111, 13.6239},
		"HTL Steyr":                        {48.04897895303849, 14.428611585052488},
	}
	loc, ok := m[org]
	return loc, ok
}

// orgLocationByID overrides orgLocation2 for entries that share an
// organization name but sit at different locations (e.g. the two ENGEL
// internships in St. Valentin and Dietach).
var orgLocationByID = map[string]orgLocation{
	"engel-intern-2017": {48.15757939951312, 14.505377793376999},
	"engel-intern-2018": {48.09266388035448, 14.436295460227777},
}

// orgLocationFor resolves a pin location, preferring an id-based override
// before falling back to the organization-name lookup.
func orgLocationFor(id, org string) (orgLocation, bool) {
	if loc, ok := orgLocationByID[id]; ok {
		return loc, true
	}
	return orgLocation2(org)
}

func thesisLogo(typ string) string {
	switch typ {
	case "HTL":
		return "/static/logos/htl-steyr.png"
	case "Bachelor", "Master", "Project", "Projekt":
		return "/static/logos/jku.jpg"
	}
	return ""
}
