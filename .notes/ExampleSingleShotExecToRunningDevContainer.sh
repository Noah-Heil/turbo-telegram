# Attach to running container
CID=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" -q);docker exec -it --workdir /workspaces/turbo-telegram -u vscode:vscode $CID bash

# Attach to running containers opencode session
CID=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" -q);docker exec -it --workdir /workspaces/turbo-telegram -u vscode:vscode $CID /home/vscode/.opencode/bin/opencode --session ses_3716ce6d2ffelIpHAv1b4j233e


# Run Opencode web in background
# This one worked
CID=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" -q);docker exec -d --workdir /workspaces/turbo-telegram -u vscode:vscode $CID /home/vscode/.opencode/bin/opencode web 

# This doesnt work
# CID=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" -q);docker exec -d --workdir /workspaces/turbo-telegram -u vscode:vscode $CID /home/vscode/.opencode/bin/opencode web --session ses_3716ce6d2ffelIpHAv1b4j233e


# Probably dont need this one?
CID=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" -q);docker exec -d --workdir /workspaces/turbo-telegram -u vscode:vscode $CID bash -c "nohup ./opencode web > opencode.log 2>&1 &"


# echo "--- Process Check ---"
CID=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" -q);docker exec $CID ps aux | grep -i "opencode" | grep -v grep

# echo "--- Port Listen Check ---"
CID=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" -q);docker exec $CID netstat -tulpn | grep LISTEN

# Try killing it nicely
CID=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" -q);docker exec -u vscode:vscode $CID kill $PID
CID=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" -q);docker exec -u vscode:vscode $CID pkill -f "/home/vscode/.opencode/bin/opencode web"
CID=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" -q);docker exec -u vscode:vscode $CID pkill -f "/home/vscode/.opencode/bin/opencode --session ses_3716ce6d2ffelIpHAv1b4j233e"

# Take old yeller behind the shed.
CID=$(docker ps --filter "label=devcontainer.local_folder=$(pwd)" -q);docker exec -u vscode:vscode $CID pkill -9 -f "opencode web"


/home/vscode/.opencode/bin/opencode web
