package global

const (
	cV = "v1" // Version of container naming

	NERDCTL = "sudo nerdctl "
	LIST    = NERDCTL + "container ls --all --format '{{.Names}} {{.ID}} {{.Image}} {{.Status}}' | sort"
	CREATE  = NERDCTL + "run -d %s --name %s-%s-" + cV + " %s" // flags, name, instance, image
	DELETE  = NERDCTL + "rm %s-%s-" + cV + " -f"               // name, instance
	START   = NERDCTL + "start %s-%s-" + cV + ""               // name, instance
	STOP    = NERDCTL + "stop %s-%s-" + cV + ""                // name, instance
	RESTART = NERDCTL + "restart %s-%s-" + cV + ""             // name, instance
	PAUSE   = NERDCTL + "pause %s-%s-" + cV + ""               // name, instance
	UNPAUSE = NERDCTL + "unpause %s-%s-" + cV + ""             // name, instance
	EXEC    = NERDCTL + "exec %s-%s %s-" + cV + ""             // name, instance, command
	LOGS    = NERDCTL + "logs %s"                              // name
)
