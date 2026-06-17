#!/usr/bin/env sh
set -eu

REPO_URL="${GOCRAFT_REPO_URL:-https://github.com/Mountok/gocraft.git}"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BIN_NAME="gocraft"

need() {
  if ! command -v "$1" >/dev/null 2>&1; then
    printf '%s\n' "error: $1 is required" >&2
    exit 1
  fi
}

install_binary() {
  src="$1"
  dst="$2"

  mkdir -p "$INSTALL_DIR" 2>/dev/null || true
  if [ -w "$INSTALL_DIR" ]; then
    cp "$src" "$dst"
    chmod 755 "$dst"
    return
  fi

  if command -v sudo >/dev/null 2>&1; then
    sudo mkdir -p "$INSTALL_DIR"
    sudo cp "$src" "$dst"
    sudo chmod 755 "$dst"
    return
  fi

  printf '%s\n' "error: cannot write to $INSTALL_DIR and sudo is not available" >&2
  printf '%s\n' "try: INSTALL_DIR=\$HOME/.local/bin sh install.sh" >&2
  exit 1
}

need git
need go

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT INT TERM

git clone --depth 1 "$REPO_URL" "$tmp_dir/gocraft" >/dev/null 2>&1
mkdir -p "$tmp_dir/build"
(cd "$tmp_dir/gocraft" && git fetch --tags --force >/dev/null 2>&1 || true)
version="dev"
for tag in $(cd "$tmp_dir/gocraft" && git tag --sort=-v:refname); do
  version="$tag"
  break
done
(cd "$tmp_dir/gocraft" && go build -ldflags "-X github.com/Mountok/gocraft/internal/version.Version=$version" -o "$tmp_dir/build/$BIN_NAME" ./cmd/gocraft)

install_binary "$tmp_dir/build/$BIN_NAME" "$INSTALL_DIR/$BIN_NAME"

printf '%s\n' "gocraft installed to $INSTALL_DIR/$BIN_NAME"
printf '%s\n' "run: gocraft help"
