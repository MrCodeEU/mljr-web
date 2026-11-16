package pages

type playgroundItem struct {
	Label   string
	DescKey string
	Href    string
	Icon    string
}

var playgroundItems = []playgroundItem{
	{Label: "AI Snake Live", DescKey: "games.ai_snake_live.card_desc", Href: "/games/ai-snake-live", Icon: "lucide:brain"},
	{Label: "AI Snake", DescKey: "games.ai_snake.card_desc", Href: "/games/ai-snake", Icon: "lucide:bot"},
	{Label: "Snake", DescKey: "games.snake.card_desc", Href: "/games/snake", Icon: "lucide:gamepad-2"},
	{Label: "Minenfeld", DescKey: "games.minenfeld.card_desc", Href: "/games/minenfeld", Icon: "lucide:bomb"},
	{Label: "Double Pendulum", DescKey: "games.pendulum.card_desc", Href: "/games/double-pendulum", Icon: "lucide:orbit"},
	{Label: "Fourier Drawing", DescKey: "games.fourier.card_desc", Href: "/games/fourier", Icon: "lucide:audio-waveform"},
	{Label: "Boids", DescKey: "games.boids.card_desc", Href: "/games/boids", Icon: "lucide:bird"},
	{Label: "Particle Life", DescKey: "games.particle_life.card_desc", Href: "/games/particle-life", Icon: "lucide:atom"},
	{Label: "Maze", DescKey: "games.maze.card_desc", Href: "/games/maze", Icon: "lucide:route"},
	{Label: "Sorting Visualizer", DescKey: "games.sorting.card_desc", Href: "/games/sorting", Icon: "lucide:bar-chart-3"},
}

// PlaygroundURLPaths returns the route paths of all playground items, for sitemap generation.
func PlaygroundURLPaths() []string {
	paths := make([]string, len(playgroundItems))
	for i, it := range playgroundItems {
		paths[i] = it.Href
	}
	return paths
}
