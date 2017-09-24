package bootstrap3

type MenuLocation int

const (
	Right = MenuLocation(iota)
	Left
)

type Nav struct {
	Name  string
	Href  string
	menus map[MenuLocation]Menu
}

func NewNav(name, href string) *Nav {
	return &Nav{
		Name:  name,
		Href:  href,
		menus: make(map[MenuLocation]Menu),
	}
}

func (n *Nav) Add(location MenuLocation, id, title, icon string) MenuItem {
	m := &menuItem{
		id:    id,
		title: title,
		icon:  icon,
	}
	n.menus[location] = append(n.menus[location], m)
	return m
}

func (n *Nav) Menu(location MenuLocation) Menu {
	return n.menus[location]
}

type MenuItem interface {
	Attrs(...string) []string
	Sub(string, string, string, string) MenuItem
	Divider()
	HasSub() bool
	Level() int
	ID() string
	Title() string
	Icon() string
	Href() string
	SubMenu() []MenuItem
	IsDivider() bool
	SetID(string)
	SetTitle(string)
	SetIcon(string)
	SetHref(string)
}

type Menu []MenuItem

type menuItem struct {
	level   int
	id      string
	title   string
	icon    string
	href    string
	subMenu Menu
	divider bool
}

func (m *menuItem) Attrs(attrs ...string) []string {
	if m.id != "" {
		attrs = append(attrs, "id", m.id)
	}
	if m.href != "" {
		attrs = append(attrs, "href", m.href)
	}
	return attrs
}

func (m *menuItem) Sub(id, title, icon, href string) MenuItem {
	mi := &menuItem{
		level: m.level + 1,
		id:    id,
		title: title,
		icon:  icon,
		href:  href,
	}
	m.subMenu = append(m.subMenu, mi)
	return mi
}

func (m *menuItem) Divider() {
	m.subMenu = append(m.subMenu, &menuItem{
		level:   m.level + 1,
		divider: true,
	})
}

func (m *menuItem) HasSub() bool {
	return len(m.subMenu) > 0
}

func (m *menuItem) Level() int    { return m.level }
func (m *menuItem) ID() string    { return m.id }
func (m *menuItem) Title() string { return m.title }
func (m *menuItem) Icon() string  { return m.icon }
func (m *menuItem) Href() string  { return m.href }
func (m *menuItem) SubMenu() []MenuItem {
	subMenu := make([]MenuItem, len(m.subMenu))
	for i, mi := range m.subMenu {
		subMenu[i] = mi
	}
	return subMenu
}
func (m *menuItem) IsDivider() bool { return m.divider }

func (m *menuItem) SetID(id string)       { m.id = id }
func (m *menuItem) SetTitle(title string) { m.title = title }
func (m *menuItem) SetIcon(icon string)   { m.icon = icon }
func (m *menuItem) SetHref(href string)   { m.href = href }
