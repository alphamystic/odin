#!/bin/bash
# stop_servers.sh

# Adjust this to match the PID_DIR in your run script
PID_DIR="./pids"

echo "🛑 Stopping servers..."

stop_process() {
    local name=$1
    local pid_file="$PID_DIR/.$2_pid"
    local ports=$3

    if [ -f "$pid_file" ]; then
        PID=$(cat "$pid_file")
        if kill "$PID" 2>/dev/null; then
            echo "✅ Stopped $name (PID: $PID)"
        else
            echo "⚠️  $name (PID: $PID) was not running."
        fi
        rm "$pid_file"
    else
        echo "🔍 No PID file for $name, searching by ports: $ports"
        # OS-specific fallback for killing by port
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS path
            for port in $ports; do
                lsof -ti :$port | xargs kill -9 2>/dev/null
            done
        else
            # Linux path
            for port in $ports; do
                sudo fuser -k "$port/tcp" 2>/dev/null
            done
        fi
        echo "✅ Attempted to clear ports for $name"
    fi
}

stop_process "Load Balancer" "lb" "4040 4041"
stop_process "API Server" "api" "5000 5001"
stop_process "UI Server" "ui" "4000 4001"

echo "Done."