package pages

import (
	g "maragu.dev/gomponents"
	h "maragu.dev/gomponents/html"

	hpdata "mljr-web/projects/homepage/data"
	"mljr-web/projects/homepage/homelab"
	"mljr-web/ui/layout"
	"mljr-web/ui/overlay"
	"mljr-web/ui/primitive"
	"mljr-web/ui/special"
	"mljr-web/ui/token"
)

const perPage = 6

func Home(d hpdata.SiteData, a AnalyticsConfig, hl homelab.Snapshot) g.Node {
	li := d.LinkedIn
	featured := d.FeaturedProjects()
	rest := d.AllProjects()
	// Featured projects get their own spotlight section; the grid shows the rest.
	gridProjects := rest
	if len(featured) == 0 {
		gridProjects = append(featured, rest...)
	}
	totalProjects := len(d.GitHub)
	headExtra := append([]g.Node{
		h.Script(h.Src("/static/motion.min.js")),
		g.El("style", g.Raw(homepageCSS)),
	}, AnalyticsHead(a)...)

	return layout.PageShell(
		layout.PageProps{
			Title:       "Michael Reinegger — Portfolio",
			Description: "Networks & IT Security · Go · self-hosted — no JS framework, no CDN, no adtech tracking.",
			Theme:       token.ThemeSwissBrut,
			Mode:        token.ModeLight,
			HeadExtra:   headExtra,
		},
		special.ThemeToggleRoot(token.ThemeSwissBrut, token.ModeLight),

		primitive.ReadProgress(primitive.ReadProgressProps{Height: "8px", Color: "var(--accent)", ZIndex: 100}),

		siteNavbar(),

		h.Main(
			h.Style("position:relative"),
			AnimatedLogoBackground(),
			heroSection(li, totalProjects),
			statsSection(d),
			experienceSection(li),
			featuredSection(featured),
			projectsSection(gridProjects),
			githubSection(d),
			homelabSection(hl),
			stravaSection(d),
			skillsSection(),
			codeShowcaseSection(),
			contactSection(),
		),

		siteFooter(),

		overlay.Toaster(overlay.ToasterProps{}),
		overlay.Portal("portal"),

		h.Script(g.Raw(`(function(){
  if(typeof Motion==='undefined') return;
  var sections=['#experience','#featured','#projects','#opensource','#homelab','#activity','#skills','#under-the-hood','#contact'];
  sections.forEach(function(sel){
    var el=document.querySelector(sel);
    if(!el) return;
    var children=Array.from(el.querySelectorAll('[data-component="card"],[data-component="timeline"],[data-component="timeline-item"]'));
    if(!children.length) children=[el];
    Motion.inView(el,function(){
      Motion.animate(children,{opacity:[0,1],y:[18,0]},{delay:Motion.stagger(0.06),duration:0.4,easing:[0.25,0.46,0.45,0.94]});
    },{margin:'-80px'});
  });
})();`)),
	)
}

const homepageCSS = `
/* homepage-specific responsive styles */
@keyframes pulse-dot {
  0%,100%{opacity:1;transform:scale(1)}
  50%{opacity:.5;transform:scale(1.5)}
}

html, body { overflow-x: hidden; }
main > section,
main > div:not(#logo-svg-hp-bg) {
  position: relative;
  z-index: 1;
}
#skills { overflow-x: hidden; }
#skills [data-component="marquee"] {
  max-width: 100vw;
  overflow-x: hidden;
}
.experience-mobile-timeline { display: none; }
.hero-stat-tile:last-child { border-right: none; }

/* Swiss-editorial section index numbers */
.section-num {
  font-size: clamp(2.4rem, 5vw, 3.6rem);
  font-weight: 900;
  line-height: 1;
  color: transparent;
  -webkit-text-stroke: 2px var(--ink);
  paint-order: stroke;
  opacity: .35;
  letter-spacing: -.02em;
  user-select: none;
}

/* heatmap SVG should shrink-wrap and scroll on small screens */
#opensource [data-component="heatmap-wrap"] svg { max-width: 100%; height: auto; }

@media (max-width: 900px) {
  .featured-grid { grid-template-columns: 1fr !important; }
  .oss-grid { grid-template-columns: 1fr !important; }
  .hood-grid { grid-template-columns: 1fr !important; }
  .hood-grid > div:first-child { position: static !important; }
  .homelab-grid { grid-template-columns: 1fr !important; }
}

/* ── Tablet (≤900px) ──────────────────────────────────────────── */
@media (max-width: 900px) {
  .hero-grid {
    grid-template-columns: 1fr !important;
    gap: var(--sp-8) !important;
  }
  .hero-bento [data-component="bento-grid"] {
    grid-template-columns: repeat(3,1fr) !important;
  }
  .bento-photo { min-height: 200px !important; }
}

/* ── Mobile landscape (≤768px) ────────────────────────────────── */
@media (max-width: 768px) {
  #hero { min-height: auto !important; padding: var(--sp-8) 0 !important; }
  #hero [data-component="container"],
  #hero-content,
  .hero-grid {
    width: 100% !important;
    max-width: 100% !important;
    min-width: 0 !important;
    overflow: hidden !important;
  }
  #hero-content {
    gap: var(--sp-4) !important;
  }
  #hero h1 {
    font-size: clamp(2rem, 14vw, 3rem) !important;
    overflow-wrap: anywhere !important;
  }
  #hero p {
    max-width: 100% !important;
    overflow-wrap: anywhere !important;
  }
  #hero [data-component="typewriter"] {
    white-space: normal !important;
    overflow-wrap: anywhere !important;
  }

  .hero-grid {
    grid-template-columns: 1fr !important;
    gap: var(--sp-5) !important;
  }

  /* hide Experience + Skills nav links on small screens */
  .nav-link-hide { display: none !important; }

  /* compact mobile bento: keep the context stats, remove the photo. */
  .hero-bento {
    display: block !important;
    width: 100% !important;
    max-width: 100% !important;
    min-width: 0 !important;
    overflow: hidden !important;
  }
  .hero-bento [data-component="bento-grid"] {
    grid-template-columns: repeat(2,minmax(0,1fr)) !important;
    grid-auto-rows: minmax(104px,auto) !important;
    height: auto !important;
    width: 100% !important;
    max-width: 100% !important;
  }
  .bento-photo { display: none !important; }
  .hero-bento [data-component="bento-item"] {
    grid-column: span 1 !important;
    grid-row: span 1 !important;
    min-width: 0 !important;
  }
  .hero-bento [data-component="bento-item"]:first-child {
    display: none !important;
  }
  .hero-bento [data-component="bento-item"]:last-child {
    grid-column: span 2 !important;
  }
  .hero-bento [data-component="card"] {
    min-width: 0 !important;
    height: 100% !important;
  }

  /* stats: 2 col */
  .hero-stat-grid {
    grid-template-columns: 1fr 1fr !important;
  }

  /* experience: single column */
  #experience [data-component="grid"] { grid-template-columns: 1fr !important; }
  .experience-snake { display: none !important; }
  .experience-mobile-timeline { display: block !important; }

  /* projects grid → single or 2 col */
  #projects [data-component="grid"] { grid-template-columns: repeat(auto-fill,minmax(280px,1fr)) !important; }

  .activity-grid { grid-template-columns: 1fr !important; }

  /* contact grid → single column */
  #contact [data-component="grid"] { grid-template-columns: 1fr !important; }

  /* logo scatter: hide on mobile (too distracting) */
  #logo-svg-hp-bg { display: none !important; }
}

/* ── Mobile portrait (≤480px) ─────────────────────────────────── */
@media (max-width: 480px) {
  .hero-ctas { flex-direction: column !important; }
  .hero-ctas a { width: 100% !important; }
  .hero-ctas a [data-component="button"] { width: 100% !important; justify-content: center !important; }

  /* stats: single column on very small */
  .hero-stat-grid {
    grid-template-columns: 1fr 1fr !important;
    gap: var(--sp-3) !important;
    padding: 0 var(--sp-4) !important;
  }
}
`
