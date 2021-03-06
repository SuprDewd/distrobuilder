package managers

// NewApt creates a new Manager instance.
func NewApt() *Manager {
	return &Manager{
		command: "apt-get",
		flags: ManagerFlags{
			clean: []string{
				"clean",
			},
			global: []string{
				"-y",
			},
			install: []string{
				"install",
			},
			remove: []string{
				"remove", "--auto-remove",
			},
			refresh: []string{
				"update",
			},
			update: []string{
				"dist-upgrade",
			},
		},
	}
}
