package global

const (
	HCLOUD_API_URL = "https://api.hetzner.cloud/v1"

	ENV_HETZNER_TOKEN = "HCLOUD_TOKEN"
	ENV_SSH_KEY_PATH  = "HCONTAINERS_SSH_KEY_PATH"

	cV    = "v1"             // Version of container naming
	cName = "%s-%s-%s-" + cV // name, portprefix, instance

	NERDCTL = "sudo nerdctl "
	LIST    = NERDCTL + "container ls --all --format '{{.Names}} {{.ID}} {{.Image}} {{.Status}}' | sort"
	CREATE  = NERDCTL + "run -d %s --name " + cName + " %s" // flags, name, instance, image
	DELETE  = NERDCTL + "rm " + cName + " -f"               // name, portprefix, instance
	START   = NERDCTL + "start " + cName                    // name, portprefix, instance
	STOP    = NERDCTL + "stop " + cName                     // name, portprefix, instance
	RESTART = NERDCTL + "restart " + cName                  // name, portprefix, instance
	PAUSE   = NERDCTL + "pause " + cName                    // name, portprefix, instance
	UNPAUSE = NERDCTL + "unpause " + cName                  // name, portprefix, instance
	EXEC    = NERDCTL + "exec " + cName + " %s"             // name, portprefix, instance, command
	LOGS    = NERDCTL + "logs " + cName                     // name, portprefix, instance, command
)
