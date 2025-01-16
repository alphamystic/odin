#!/bin/bash

# Installer for Odin Project

# Set repository details
REPO_URL="https://github.com/alphamystic/odin.git"
BRANCH="update"

# Get the user's home directory
HOME_DIR="$HOME"

# Set the target directory
TARGET_DIR="$HOME_DIR/.odin"

# Check if the .odin directory already exists
if [ -d "$TARGET_DIR" ]; then
    echo "Directory $TARGET_DIR already exists. Removing it to re-clone."
    rm -rf "$TARGET_DIR"
fi

# Clone the update branch into the .odin directory
echo "Cloning the 'update' branch of $REPO_URL into $TARGET_DIR..."
git clone --branch $BRANCH $REPO_URL $TARGET_DIR

# Check if cloning was successful
if [ $? -ne 0 ]; then
    echo "Failed to clone the repository. Please check your internet connection and repository URL."
    exit 1
fi

# Navigate to the .odin directory
cd "$TARGET_DIR"

# Check if odin.go exists
if [ ! -f "odin.go" ]; then
    echo "File odin.go not found in $TARGET_DIR. Exiting."
    exit 1
fi

# Run the odin.go file
echo "Running odin.go..."
go run odin.go

# Check if the Go script ran successfully
if [ $? -ne 0 ]; then
    echo "Failed to run odin.go. Please ensure Go is installed and properly configured on your system."
    exit 1
fi

echo "Odin project setup and execution completed successfully."
