// Package ui exposes Datastar attribute helpers consumed by components and
// pages. Centralizing the data-* spelling means any future Datastar attribute
// rename is a one-file change.
package ui

import (
	g "maragu.dev/gomponents"
)

// Datastar v1.x attribute syntax:
//   plugins with a key use a colon separator: data-on:click, data-bind:name
//   plugins without a key use the name directly: data-signals, data-show, data-text, data-effect
//   modifiers append with double-underscore: data-on:click__prevent__stop

// On binds a Datastar event handler.
// event = "click", "input:keydown", etc.
// Modifiers: On("click__debounce.300ms", expr), On("click__prevent", expr).
func On(event, expr string) g.Node { return g.Attr("data-on:"+event, expr) }

// Bind sets data-bind:<signal> for two-way input binding.
func Bind(signal string) g.Node { return g.Attr("data-bind:" + signal) }

// Text binds element textContent to a Datastar expression.
func Text(expr string) g.Node { return g.Attr("data-text", expr) }

// Show toggles element visibility based on a Datastar expression.
func Show(expr string) g.Node { return g.Attr("data-show", expr) }

// Signals declares Datastar signals at the element. Pass a JS object literal.
func Signals(obj string) g.Node { return g.Attr("data-signals", obj) }

// Signal seeds a single signal: data-signals-<name>="<value>".
// Note: individual signal seeding uses hyphen (no key separator needed).
func Signal(name, value string) g.Node { return g.Attr("data-signals-"+name, value) }

// DSAttr drives reactive attribute updates: data-attr="{'data-theme':$theme}".
func DSAttr(obj string) g.Node { return g.Attr("data-attr", obj) }

// Indicator marks a fetch-in-progress indicator signal.
func Indicator(signal string) g.Node { return g.Attr("data-indicator:" + signal) }

// Attrs is sugar for variadic pass-through.
func Attrs(nodes ...g.Node) []g.Node { return nodes }
