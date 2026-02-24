Installation & Configuration

Add the opencode feature to your Dev Container by including it in your devcontainer.json:

{
	"features": {
		"ghcr.io/danzilberdan/devcontainers/opencode:0": {}
	}
}

opencode centrally manages credentials for enabled inference providers in ~/.local/share/opencode/auth.json. In order to expose this file to opencode running in the Dev Container do the following.

# In your project dir
mkdir .opencode
cd .opencode
# Create a hard link to your host's opencode `auth.json` file.
ln ~/.local/share/opencode/auth.json

Add the following line to .gitignore so that the credentials will not be commited to git:

.opencode/

Add the following to your devcontainer.json to expose the credentials to the Dev Container:

    "mounts": [
        {
            "source": "${localWorkspaceFolder}/.opencode/auth.json",
            "target": "/mnt/opencode-auth.json",
            "type": "bind"
        }
    ],

Now opencode in the Dev Container will have access host configured inference providers.