package marionette

type Actions struct {
	Actions []*InputActions `json:"actions"`
}

func (a *Actions) Pointer(id, pointerType string) (ret *InputActions) {
	ret = &InputActions{
		Type: "pointer",
		Id:   id,
		Parameters: PointerParameters{
			PointerType: pointerType,
		},
	}
	a.Actions = append(a.Actions, ret)
	return
}

func (a *Actions) Wheel(id string) (ret *InputActions) {
	ret = &InputActions{
		Type: "wheel",
		Id:   id,
	}
	a.Actions = append(a.Actions, ret)
	return
}

func (a *Actions) Key() (ret *InputActions) {
	ret = &InputActions{
		Type: "key",
	}
	a.Actions = append(a.Actions, ret)
	return
}

type InputActions struct {
	Type       string `json:"type"`
	Id         string `json:"id,omitEmpty"`
	Parameters any    `json:"parameters,omitEmpty"`
	Actions    []any  `json:"actions"`
}

type PointerParameters struct {
	// PointerType is "mouse", "pen" or "touch"
	PointerType string `json:"pointerType,omitempty"`
}

func (ia *InputActions) Add(action any) {
	switch v := action.(type) {
	case Pause:
		v.Type = "pause"
		ia.Actions = append(ia.Actions, v)

	case PointerMove:
		v.Type = "pointerMove"
		ia.Actions = append(ia.Actions, v)
	case PointerUp:
		v.Type = "pointerUp"
		ia.Actions = append(ia.Actions, v)
	case PointerDown:
		v.Type = "pointerDown"
		ia.Actions = append(ia.Actions, v)

	case Scroll:
		v.Type = "scroll"
		ia.Actions = append(ia.Actions, v)

	case KeyUp:
		v.Type = "keyUp"
		ia.Actions = append(ia.Actions, v)
	case KeyDown:
		v.Type = "keyDown"
		ia.Actions = append(ia.Actions, v)

	default:
		panic("invalid action")
	}
}

type Pause struct {
	Type     string `json:"type"`
	Duration int    `json:"duration"`
}

type PointerMove struct {
	Type     string `json:"type"`
	Duration int    `json:"duration"`
	Origin   string `json:"origin,omitempty"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}
type PointerUp struct {
	Type   string `json:"type"`
	Button int    `json:"button"`
}
type PointerDown struct {
	Type   string `json:"type"`
	Button int    `json:"button"`
}

type Scroll struct {
	Type     string `json:"type"`
	Duration int    `json:"duration"`
	Origin   string `json:"origin,omitempty"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
	DeltaX   int    `json:"deltaX"`
	DeltaY   int    `json:"deltaY"`
}

type KeyDown struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
type KeyUp struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
