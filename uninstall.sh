#!/usr/bin/env sh
set -eu

INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BIN_NAME="gocraft"
TARGET="$INSTALL_DIR/$BIN_NAME"

remove_binary() {
  if [ ! -e "$TARGET" ]; then
    printf '%s\n' "gocraft is not installed at $TARGET"
    return
  fi

  if [ -w "$INSTALL_DIR" ]; then
    rm -f "$TARGET"
    return
  fi

  if command -v sudo >/dev/null 2>&1; then
    sudo rm -f "$TARGET"
    return
  fi

  printf '%s\n' "error: cannot remove $TARGET and sudo is not available" >&2
  exit 1
}

remove_binary

if [ -e "$TARGET" ]; then
  printf '%s\n' "error: failed to uninstall gocraft from $TARGET" >&2
  exit 1
fi

printf '%s\n' "gocraft uninstalled from $TARGET"
