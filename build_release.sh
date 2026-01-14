#!/bin/bash

# Configuration
APP_NAME="relix"
VERSION="0.1.0"
DIST_DIR="dist"

# Platforms to build for
PLATFORMS=("darwin/amd64" "darwin/arm64" "linux/amd64" "linux/arm64")

# Clean dist directory
rm -rf $DIST_DIR
mkdir -p $DIST_DIR

echo "Building Relix v$VERSION..."

for PLATFORM in "${PLATFORMS[@]}"; do
    OS="${PLATFORM%%/*}"
    ARCH="${PLATFORM#*/}"
    OUTPUT_NAME="$APP_NAME"
    
    if [ "$OS" == "windows" ]; then
        OUTPUT_NAME="$APP_NAME.exe"
    fi

    echo "Building for $OS/$ARCH..."
    
    # Build
    env GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w" -o "$DIST_DIR/$OS-$ARCH/$OUTPUT_NAME" .
    
    if [ $? -ne 0 ]; then
        echo "Error building for $OS/$ARCH"
        exit 1
    fi

    # Archive
    cd $DIST_DIR/$OS-$ARCH
    ARCHIVE_NAME="${APP_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
    tar -czf "../$ARCHIVE_NAME" "$OUTPUT_NAME"
    cd - > /dev/null
    
    echo "Created $DIST_DIR/$ARCHIVE_NAME"
done

# Clean up uncompressed binaries
rm -rf "$DIST_DIR/darwin-amd64" "$DIST_DIR/darwin-arm64" "$DIST_DIR/linux-amd64" "$DIST_DIR/linux-arm64"

echo "Build complete! Artifacts are in $DIST_DIR/"
ls -lh $DIST_DIR
