package models

type Menu struct {
	Path      string   `json:"path"`
	Component string   `json:"component"`
	Redirect  string   `json:"redirect,omitempty"`
	Meta      Menumeta `json:"meta"`
	Children  []Menu   `json:"children,omitempty"`
	Name      string   `json:"name"`
}

type Menumeta struct {
	AlwaysShow bool     `json:"alwaysShow,omitempty"`
	Hidden     bool     `json:"hidden"`
	Icon       string   `json:"icon"`
	KeepAlive  bool     `json:"keepAlive"`
	Roles      []string `json:"roles"`
	Title      string   `json:"title"`
}
