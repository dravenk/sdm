package main

var (
	projectPath = "projects/drupal"
	tplSettings = "settings.php"
	tplCompose  = "docker-compose.yaml"
	appPath     = "drupal"

	cmdInput string

	InputUp     = "up"
	InputInit   = "init"
	InputDown   = "down"
	InputRemove = "remove"

	// Used for flags.
	cfgFile string
)
