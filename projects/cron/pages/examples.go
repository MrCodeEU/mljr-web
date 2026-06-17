package pages

type CronExample struct {
	Key  string
	Name string
	Expr string
}

type CronExampleGroup struct {
	Name  string
	Items []CronExample
}

var ExampleGroups = []CronExampleGroup{
	{
		Name: "Common",
		Items: []CronExample{
			{"every-minute", "Every minute", "* * * * *"},
			{"every-5min", "Every 5 minutes", "*/5 * * * *"},
			{"every-hour", "Every hour", "@hourly"},
			{"every-day", "Every day at midnight", "@daily"},
			{"every-week", "Every Sunday", "@weekly"},
			{"every-month", "1st of every month", "@monthly"},
			{"every-year", "Once a year", "@yearly"},
		},
	},
	{
		Name: "Business Hours",
		Items: []CronExample{
			{"workday-9am", "Weekdays at 9 AM", "0 9 * * 1-5"},
			{"workday-5pm", "Weekdays at 5 PM", "0 17 * * 1-5"},
			{"lunch", "Lunch reminder (Mon–Fri 12:30)", "30 12 * * 1-5"},
			{"standup", "Daily standup 9:15 AM", "15 9 * * 1-5"},
		},
	},
	{
		Name: "Maintenance",
		Items: []CronExample{
			{"backup-nightly", "Nightly backup at 3 AM", "0 3 * * *"},
			{"cleanup-weekly", "Weekly cleanup Sunday 2 AM", "0 2 * * 0"},
			{"monthly-report", "Monthly report 1st at 8 AM", "0 8 1 * *"},
			{"quarterly", "Quarterly (1st Jan/Apr/Jul/Oct)", "0 0 1 1,4,7,10 *"},
		},
	},
	{
		Name: "With Seconds",
		Items: []CronExample{
			{"every-10sec", "Every 10 seconds", "*/10 * * * * *"},
			{"every-30sec", "Every 30 seconds", "*/30 * * * * *"},
			{"health-check", "Health check every 15s", "*/15 * * * * *"},
		},
	},
	{
		Name: "Advanced",
		Items: []CronExample{
			{"every-interval", "Every 2 hours", "0 */2 * * *"},
			{"twice-daily", "Twice daily (8 AM & 8 PM)", "0 8,20 * * *"},
			{"last-friday", "First Mon of month 9 AM", "0 9 1-7 * 1"},
			{"every-15-biz", "Every 15 min, 9–17 weekdays", "*/15 9-17 * * 1-5"},
		},
	},
}

func FindExample(key string) (CronExample, bool) {
	for _, g := range ExampleGroups {
		for _, e := range g.Items {
			if e.Key == key {
				return e, true
			}
		}
	}
	return CronExample{}, false
}
