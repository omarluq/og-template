#!/usr/bin/env bash
set -euo pipefail

OLD_MODULE="github.com/username/myapp"
OLD_BINARY="myapp"
OLD_PREFIX="MYAPP"

gum style --bold --foreground 212 "og-template init"
echo ""

MODULE=$(gum input --prompt "Module path: " --placeholder "$OLD_MODULE" --value "$OLD_MODULE")
BINARY=$(gum input --prompt "Binary name: " --placeholder "$OLD_BINARY" --value "$OLD_BINARY")
PREFIX=$(gum input --prompt "Env prefix:  " --placeholder "$OLD_PREFIX" --value "$(echo "$BINARY" | tr '[:lower:]-' '[:upper:]_')")

echo ""
gum style --faint "Module:  $OLD_MODULE → $MODULE"
gum style --faint "Binary:  $OLD_BINARY → $BINARY"
gum style --faint "Prefix:  $OLD_PREFIX → $PREFIX"
echo ""

if ! gum confirm "Apply changes?"; then
  gum style --foreground 196 "Aborted."
  exit 1
fi

echo ""

# Module path
gum spin --title "Updating module paths..." -- \
  bash -c "find . -type f \( -name '*.go' -o -name '*.yml' -o -name '*.yaml' -o -name '*.mod' -o -name 'Taskfile*' -o -name '*.md' \) \
    -not -path './.git/*' -not -path './scripts/*' \
    -exec sed -i 's|$OLD_MODULE|$MODULE|g' {} +"

# Binary name
gum spin --title "Updating binary name..." -- \
  bash -c "find . -type f \( -name '*.go' -o -name '*.yml' -o -name '*.yaml' -o -name 'Taskfile*' -o -name '*.md' \) \
    -not -path './.git/*' -not -path './scripts/*' \
    -exec sed -i 's|$OLD_BINARY|$BINARY|g' {} +"

# Rename cmd directory
if [ -d "cmd/$OLD_BINARY" ] && [ "$OLD_BINARY" != "$BINARY" ]; then
  mv "cmd/$OLD_BINARY" "cmd/$BINARY"
fi

# Env prefix
gum spin --title "Updating env prefix..." -- \
  bash -c "find . -type f \( -name '*.go' -o -name '*.md' \) \
    -not -path './.git/*' -not -path './scripts/*' \
    -exec sed -i 's|$OLD_PREFIX|$PREFIX|g' {} +"

# Remove gum from mise
sed -i '/# Interactive prompts for init/d' .mise.toml
sed -i '/charmbracelet\/gum/d' .mise.toml
# Remove trailing blank lines from mise
sed -i -e :a -e '/^\n*$/{$d;N;ba}' .mise.toml

# Remove init task from Taskfile
sed -i '/^  init:$/,/^$/d' Taskfile.yml

# Tidy deps
gum spin --title "Running go mod tidy..." -- go mod tidy

# Self-destruct: remove scripts dir
rm -rf scripts

echo ""
gum style --bold --foreground 76 "Done!"
echo ""
echo "Next steps:"
echo "  mise exec -- task ci"
echo "  git add -A && git commit -m 'feat: initialize $BINARY'"
