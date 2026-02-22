#!/bin/bash
# stop_servers.sh
# Stops Load Balancer, API, and UI servers

echo "Stopping servers..."

# Stop using saved PIDs if available
if [ -f .lb_pid ]; then
    kill $(cat .lb_pid) 2>/dev/null && echo "Stopped Load Balancer"
    rm .lb_pid
else
    sudo fuser -k 4040/tcp 4041/tcp 2>/dev/null
fi

if [ -f .api_pid ]; then
    kill $(cat .api_pid) 2>/dev/null && echo "Stopped API Server"
    rm .api_pid
else
    sudo fuser -k 5000/tcp 5001/tcp 2>/dev/null
fi

if [ -f .ui_pid ]; then
    kill $(cat .ui_pid) 2>/dev/null && echo "Stopped UI Server"
    rm .ui_pid
else
    sudo fuser -k 4000/tcp 4001/tcp 2>/dev/null
fi

echo "✅ All servers stopped."
