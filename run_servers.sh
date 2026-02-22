#!/bin/bash
# run_servers.sh
# Starts Load Balancer, API, and UI servers

# Define paths (adjust if necessary)
LOADBALANCER_PATH="/home/sam/Documents/3l0racle/odin"
API_PATH="/home/sam/Documents/3l0racle/odin"
UI_PATH="/home/sam/Documents/3l0racle/odin"

echo "🚀 Starting all servers..."


# Start API Server
cd "$API_PATH" || exit 1
./api --mode=DEV --ports=5001 --port=5000 --tls=true  &
API_PID=$!
echo "API Server running https on 5001 and HTTP on 5000 (PID: $API_PID)"

# Start UI Server (port 3000)
cd "$UI_PATH" || exit 1
./server --mode=DEV --ports=4001 --port=4000 --tls=true  &
UI_PID=$!
echo "UI Loki Server running https on 4001 and HTTP on 4000 (PID: $API_PID)"
echo "Waiting for UI to be ready..."

# Optional: Wait until UI is actually up (up to 30s)
for i in {1..30}; do
  if curl -s http://localhost:4000 > /dev/null; then
    echo "✅ UI Server is up!"
    break
  fi
  sleep 1
done


# Start Load Balancer (port 9090)
cd "$LOADBALANCER_PATH" || exit 1
./lb &
LB_PID=$!
echo "Load Balancer running HTTPS on 4040 and HTTP on 4041 (PID: $LB_PID)"

# Save PIDs for later termination
echo "$LB_PID" > "$PID_DIR/.lb_pid"
echo "$API_PID" > "$PID_DIR/.api_pid"
echo "$UI_PID" > "$PID_DIR/.ui_pid"

echo "✅ All servers started successfully!"
