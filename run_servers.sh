#!/bin/bash
# run_servers.sh
# Starts Load Balancer, API, and UI servers

# Define paths (adjust if necessary)
#LOADBALANCER_PATH="/home/sam/Documents/3l0racle/odin"
#API_PATH="/home/sam/Documents/3l0racle/odin"
#UI_PATH="/home/sam/Documents/3l0racle/odin"
#LOADBALANCER_PATH="/Users/sodhiambo/Documents/projects/odin"
#API_PATH="/Users/sodhiambo/Documents/projects/odin"
#UI_PATH="/Users/sodhiambo/Documents/projects/odin"

#!/bin/bash
# run_servers.sh - Multi-OS Server Starter

# --- Configuration ---
# Update this path to your project root
PROJECT_PATH="/Users/sodhiambo/Documents/projects/odin"
PID_DIR="$PROJECT_PATH/pids"
LOG_DIR="$PROJECT_PATH/.data/logs"

mkdir -p "$PID_DIR"
mkdir -p "$LOG_DIR"

# --- OS Detection & Binary Selection ---
OS_TYPE=$(uname -s)
case "$OS_TYPE" in
    Linux*)
        API_BIN="./api"
        UI_BIN="./server"
        LB_BIN="./lb"
        echo "🐧 System: Linux" ;;
    Darwin*)
        API_BIN="./api_mac"
        UI_BIN="./server_mac"
        LB_BIN="./lb_mac"
        echo "🍎 System: macOS" ;;
    CYGWIN*|MINGW*|MSYS*)
        API_BIN="./api_nt.exe"
        UI_BIN="./server_nt.exe"
        LB_BIN="./lb_nt.exe"
        echo "🪟 System: Windows (Bash)" ;;
    *)
        echo "❌ Unsupported OS: $OS_TYPE"
        exit 1 ;;
esac

echo "🚀 Initializing Odin Infrastructure..."

# --- 1. Start API Server ---
# The LB depends on the API being available
echo "Starting API Server..."
cd "$PROJECT_PATH" || exit 1
$API_BIN --mode=DEV --ports=5001 --port=5000 --tls=true > "$LOG_DIR/api.log" 2>&1 &
API_PID=$!
echo "$API_PID" > "$PID_DIR/.api_pid"
echo "   - API running (PID: $API_PID). Logs: $LOG_DIR/api.log"

# --- 2. Start UI Server ---
echo "Starting UI Server..."
cd "$PROJECT_PATH" || exit 1
$UI_BIN > "$LOG_DIR/ui.log" 2>&1 &
UI_PID=$!
echo "$UI_PID" > "$PID_DIR/.ui_pid"
echo "   - UI running (PID: $UI_PID). Logs: $LOG_DIR/ui.log"

# --- 3. Health Check Wait ---
# We wait a few seconds to ensure ports 5000 and 4000 are actually open
echo "Waiting for services to initialize..."
sleep 3
for i in {1..5}; do
  if curl -s http://localhost:4000 > /dev/null && curl -s http://localhost:5000 > /dev/null; then
    echo "✅ Backends are responsive."
    break
  fi
  echo "   ...still waiting..."
  sleep 2
done

# --- 4. Start Load Balancer ---
# Start this last so it finds the backends active
echo "Starting Load Balancer..."
cd "$PROJECT_PATH" || exit 1
$LB_BIN > "$LOG_DIR/lb.log" 2>&1 &
LB_PID=$!
echo "$LB_PID" > "$PID_DIR/.lb_pid"
echo "   - LB running (PID: $LB_PID). Logs: $LOG_DIR/lb.log"

echo "------------------------------------------------"
echo "✅ All servers started successfully!"
echo "API: http://localhost:5000"
echo "UI:  http://localhost:4000"
echo "LB:  http://localhost:4041 (per .env)"
echo "------------------------------------------------"