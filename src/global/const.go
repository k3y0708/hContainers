package global

const (
	NERDCTL = "sudo nerdctl "
	LIST    = NERDCTL + "container ls --all --format '{{.Names}} {{.ID}} {{.Image}} {{.Status}}'"
	CREATE  = NERDCTL + "run -d %s --name %s-001 %s" // flags, name, image
	DELETE  = NERDCTL + "rm %s -f"                   // name
	START   = NERDCTL + "start %s"                   // name
	STOP    = NERDCTL + "stop %s"                    // name
	RESTART = NERDCTL + "restart %s"                 // name
	PAUSE   = NERDCTL + "pause %s"                   // name
	UNPAUSE = NERDCTL + "unpause %s"                 // name
	EXEC    = NERDCTL + "exec %s %s"                 // name, command
	LOGS    = NERDCTL + "logs %s"                    // name
)
