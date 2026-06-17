package pages

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type CronInput struct {
	Expression string
	Count      int
	WithSec    bool
}

type FieldInfo struct {
	Name  string
	Value string
	Desc  string
}

type CronResult struct {
	Err        string
	Expression string
	Fields     []FieldInfo
	Next       []time.Time
	Human      string
}

func EvalCron(inp CronInput) CronResult {
	expr := strings.TrimSpace(inp.Expression)
	if expr == "" {
		return CronResult{Err: "enter a cron expression"}
	}
	count := inp.Count
	if count <= 0 {
		count = 10
	}

	var parser cron.Parser
	withSec := inp.WithSec || looksLikeSixField(expr)
	if withSec {
		parser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	} else {
		parser = cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	}

	sched, err := parser.Parse(expr)
	if err != nil {
		return CronResult{Err: err.Error(), Expression: expr}
	}

	nexts := make([]time.Time, 0, count)
	t := time.Now()
	for i := 0; i < count; i++ {
		t = sched.Next(t)
		if t.IsZero() {
			break
		}
		nexts = append(nexts, t)
	}

	fields := parseFields(expr, withSec)
	human := humanize(expr, withSec)

	return CronResult{
		Expression: expr,
		Fields:     fields,
		Next:       nexts,
		Human:      human,
	}
}

func looksLikeSixField(expr string) bool {
	if strings.HasPrefix(expr, "@") {
		return false
	}
	parts := strings.Fields(expr)
	return len(parts) >= 6
}

func parseFields(expr string, withSec bool) []FieldInfo {
	if strings.HasPrefix(expr, "@") {
		return nil
	}
	parts := strings.Fields(expr)

	names := []string{"Minute", "Hour", "Day", "Month", "Weekday"}
	if withSec {
		names = []string{"Second", "Minute", "Hour", "Day", "Month", "Weekday"}
	}

	fields := make([]FieldInfo, 0, len(names))
	for i, name := range names {
		val := "*"
		if i < len(parts) {
			val = parts[i]
		}
		fields = append(fields, FieldInfo{
			Name:  name,
			Value: val,
			Desc:  describeField(name, val),
		})
	}
	return fields
}

func describeField(name, val string) string {
	if val == "*" {
		return "every " + strings.ToLower(name)
	}
	if val == "?" {
		return "any"
	}
	if strings.Contains(val, "/") {
		parts := strings.SplitN(val, "/", 2)
		step := parts[1]
		base := parts[0]
		unit := strings.ToLower(name)
		if base == "*" || base == "0" {
			return fmt.Sprintf("every %s %s", step, unit)
		}
		return fmt.Sprintf("every %s %s starting at %s", step, unit, base)
	}
	if strings.Contains(val, "-") {
		return fmt.Sprintf("%s %s–%s", strings.ToLower(name), val[:strings.Index(val, "-")], val[strings.Index(val, "-")+1:])
	}
	if strings.Contains(val, ",") {
		return strings.ToLower(name) + " " + val
	}
	return fieldValueLabel(name, val)
}

func fieldValueLabel(name, val string) string {
	switch name {
	case "Month":
		months := map[string]string{
			"1": "Jan", "2": "Feb", "3": "Mar", "4": "Apr",
			"5": "May", "6": "Jun", "7": "Jul", "8": "Aug",
			"9": "Sep", "10": "Oct", "11": "Nov", "12": "Dec",
		}
		if m, ok := months[val]; ok {
			return m
		}
	case "Weekday":
		days := map[string]string{
			"0": "Sun", "1": "Mon", "2": "Tue", "3": "Wed",
			"4": "Thu", "5": "Fri", "6": "Sat", "7": "Sun",
		}
		if d, ok := days[val]; ok {
			return d
		}
	}
	return "at " + val
}

func humanize(expr string, withSec bool) string {
	switch expr {
	case "@yearly", "@annually":
		return "Once a year (Jan 1 at midnight)"
	case "@monthly":
		return "Once a month (1st at midnight)"
	case "@weekly":
		return "Once a week (Sunday at midnight)"
	case "@daily", "@midnight":
		return "Once a day at midnight"
	case "@hourly":
		return "Once an hour"
	}
	if strings.HasPrefix(expr, "@every ") {
		return "Every " + strings.TrimPrefix(expr, "@every ")
	}

	parts := strings.Fields(expr)
	offset := 0
	if withSec {
		offset = 1
	}
	if len(parts) < 5+offset {
		return ""
	}

	minute := parts[offset]
	hour := parts[1+offset]
	dom := parts[2+offset]
	month := parts[3+offset]
	dow := parts[4+offset]

	var sb strings.Builder

	// Time part
	if hour == "*" && minute == "*" {
		sb.WriteString("Every minute")
	} else if hour == "*" {
		sb.WriteString(fmt.Sprintf("At minute %s of every hour", minute))
	} else if minute == "*" {
		sb.WriteString(fmt.Sprintf("Every minute of hour %s", hour))
	} else {
		sb.WriteString(fmt.Sprintf("At %s:%s", padTwo(hour), padTwo(minute)))
	}

	if withSec && offset == 1 {
		sec := parts[0]
		if sec != "0" && sec != "*" {
			sb.WriteString(fmt.Sprintf(" (second %s)", sec))
		}
	}

	// Day part
	if dom != "*" && dom != "?" && dow == "*" {
		sb.WriteString(fmt.Sprintf(", on day %s of the month", dom))
	} else if dow != "*" && dow != "?" && dom == "*" {
		sb.WriteString(fmt.Sprintf(", on %s", expandDow(dow)))
	} else if dom != "*" && dow != "*" {
		sb.WriteString(fmt.Sprintf(", on day %s or %s", dom, expandDow(dow)))
	}

	if month != "*" {
		sb.WriteString(fmt.Sprintf(", in %s", expandMonth(month)))
	}

	return sb.String()
}

func padTwo(s string) string {
	if len(s) == 1 {
		return "0" + s
	}
	return s
}

func expandDow(s string) string {
	days := map[string]string{
		"0": "Sunday", "1": "Monday", "2": "Tuesday", "3": "Wednesday",
		"4": "Thursday", "5": "Friday", "6": "Saturday", "7": "Sunday",
		"MON": "Monday", "TUE": "Tuesday", "WED": "Wednesday",
		"THU": "Thursday", "FRI": "Friday", "SAT": "Saturday", "SUN": "Sunday",
	}
	if d, ok := days[strings.ToUpper(s)]; ok {
		return d
	}
	return s
}

func expandMonth(s string) string {
	months := map[string]string{
		"1": "January", "2": "February", "3": "March", "4": "April",
		"5": "May", "6": "June", "7": "July", "8": "August",
		"9": "September", "10": "October", "11": "November", "12": "December",
	}
	if m, ok := months[s]; ok {
		return m
	}
	return s
}
