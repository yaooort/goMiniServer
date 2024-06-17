package theme

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

// Declare conformity with Layout interface
var _ fyne.Layout = (*extendBoxLayout)(nil)

type extendBoxLayout struct {
	horizontal  bool
	extendIndex int
}

func NewHExtendBox(extendIndex int, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(NewHExtendBoxLayout(extendIndex), objects...)
}

// NewHExtendBoxLayout returns a horizontal box layout for stacking a number of child
// canvas objects or widgets left to right.
func NewHExtendBoxLayout(extendIndex int) fyne.Layout {
	return &extendBoxLayout{true, extendIndex}
}
func NewVExtendBox(extendIndex int, objects ...fyne.CanvasObject) *fyne.Container {
	return container.New(NewVExtendBoxLayout(extendIndex), objects...)
}

// NewVExtendBoxLayout returns a vertical box layout for stacking a number of child
// canvas objects or widgets top to bottom.
func NewVExtendBoxLayout(extendIndex int) fyne.Layout {
	return &extendBoxLayout{false, extendIndex}
}

func isVerticalSpacer(obj fyne.CanvasObject) bool {
	if spacer, ok := obj.(layout.SpacerObject); ok {
		return spacer.ExpandVertical()
	}

	return false
}

func isHorizontalSpacer(obj fyne.CanvasObject) bool {
	if spacer, ok := obj.(layout.SpacerObject); ok {
		return spacer.ExpandHorizontal()
	}

	return false
}

func (g *extendBoxLayout) isSpacer(obj fyne.CanvasObject) bool {
	// invisible spacers don't impact layout
	if !obj.Visible() {
		return false
	}

	if g.horizontal {
		return isHorizontalSpacer(obj)
	}
	return isVerticalSpacer(obj)
}

// Layout is called to pack all child objects into a specified size.
// For a VBoxLayout this will pack objects into a single column where each item
// is full width but the height is the minimum required.
// Any spacers added will pad the view, sharing the space if there are two or more.
func (g *extendBoxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	spacers := make([]fyne.CanvasObject, 0)
	total := float32(0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if g.isSpacer(child) {
			spacers = append(spacers, child)
			continue
		}
		if g.horizontal {
			total += child.MinSize().Width
		} else {
			total += child.MinSize().Height
		}
	}

	x, y := float32(0), float32(0)
	var extra float32
	if g.horizontal {
		extra = size.Width - total - (theme.Padding() * float32(len(objects)-len(spacers)-1))
	} else {
		extra = size.Height - total - (theme.Padding() * float32(len(objects)-len(spacers)-1))
	}
	extraCell := float32(0)
	if len(spacers) > 0 {
		extraCell = extra / float32(len(spacers))
	}
	lastSize := float32(0)
	for i, child := range objects {
		if !child.Visible() {
			continue
		}
		width := child.MinSize().Width
		height := child.MinSize().Height

		if i > g.extendIndex {
			if g.horizontal {
				lastSize += theme.Padding() + width
			} else {
				lastSize += theme.Padding() + height
			}
		}
	}
	for i, child := range objects {
		if !child.Visible() {
			continue
		}
		width := float32(0)
		height := float32(0)
		if i == g.extendIndex {
			width = size.Width - x - lastSize
		} else {
			width = child.MinSize().Width
		}
		if i == g.extendIndex {
			height = size.Height - y - lastSize
		} else {
			height = child.MinSize().Height
		}

		if g.isSpacer(child) {
			if g.horizontal {
				x += extraCell
			} else {
				y += extraCell
			}
			continue
		}
		child.Move(fyne.NewPos(x, y))

		if g.horizontal {
			x += theme.Padding() + width
			child.Resize(fyne.NewSize(width, size.Height))
		} else {
			y += theme.Padding() + height
			child.Resize(fyne.NewSize(size.Width, height))
		}
	}
}

// MinSize finds the smallest size that satisfies all the child objects.
// For a BoxLayout this is the width of the widest item and the height is
// the sum of of all children combined with padding between each.
func (g *extendBoxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	addPadding := false
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		if g.isSpacer(child) {
			continue
		}

		if g.horizontal {
			minSize.Height = fyne.Max(child.MinSize().Height, minSize.Height)
			minSize.Width += child.MinSize().Width
			if addPadding {
				minSize.Width += theme.Padding()
			}
		} else {
			minSize.Width = fyne.Max(child.MinSize().Width, minSize.Width)
			minSize.Height += child.MinSize().Height
			if addPadding {
				minSize.Height += theme.Padding()
			}
		}
		addPadding = true
	}
	return minSize
}
