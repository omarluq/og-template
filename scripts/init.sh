#!/usr/bin/env bash
set -euo pipefail

OLD_MODULE="github.com/username/myapp"
OLD_BINARY="myapp"
OLD_PREFIX="MYAPP"

gum style --bold --foreground 212 "og-template init"
echo ""

echo "Module path (e.g. github.com/yourname/yourproject):"
MODULE=$(gum input --placeholder "github.com/yourname/yourproject")
gum style --faint "  → $MODULE"

echo ""
echo "Binary name (e.g. yourproject):"
BINARY=$(gum input --placeholder "yourproject")
gum style --faint "  → $BINARY"

echo ""
echo "Environment variable prefix (e.g. YOURPROJECT):"
PREFIX=$(gum input --placeholder "YOURPROJECT")
gum style --faint "  → $PREFIX"

echo ""
echo "Summary:"
gum style --faint "  Module:  $OLD_MODULE → $MODULE"
gum style --faint "  Binary:  $OLD_BINARY → $BINARY"
gum style --faint "  Prefix:  $OLD_PREFIX → $PREFIX"
gum style --faint "  Cmd dir: cmd/$OLD_BINARY → cmd/$BINARY"
echo ""

if ! gum confirm "Apply changes?"; then
  gum style --foreground 196 "Aborted."
  exit 0
fi

echo ""

EXCLUDE="-not -path './.git/*' -not -path './scripts/*' -not -path './.gocache/*' -not -path './.gomodcache/*' -not -path './.tmp/*' -not -path './bin/*' -not -path './.agents/*'"

echo "Updating module paths..."
eval "find . -type f \( -name '*.go' -o -name '*.yml' -o -name '*.yaml' -o -name '*.mod' -o -name 'Taskfile*' -o -name '*.md' \) ${EXCLUDE}" \
  | xargs sed -i "s|${OLD_MODULE}|${MODULE}|g"

echo "Updating binary name..."
eval "find . -type f \( -name '*.go' -o -name '*.yml' -o -name '*.yaml' -o -name 'Taskfile*' -o -name '*.md' \) ${EXCLUDE}" \
  | xargs sed -i "s|${OLD_BINARY}|${BINARY}|g"

if [ -d "cmd/${OLD_BINARY}" ] && [ "${OLD_BINARY}" != "${BINARY}" ]; then
  echo "Renaming cmd/${OLD_BINARY} → cmd/${BINARY}..."
  mv "cmd/${OLD_BINARY}" "cmd/${BINARY}"
fi

echo "Updating env prefix..."
eval "find . -type f \( -name '*.go' -o -name '*.md' \) ${EXCLUDE}" \
  | xargs sed -i "s|${OLD_PREFIX}|${PREFIX}|g"

echo "Cleaning up init scaffolding..."
sed -i '/# Interactive prompts for init/d' .mise.toml
sed -i '/charmbracelet\/gum/d' .mise.toml
sed -i -e :a -e '/^\n*$/{$d;N;ba}' .mise.toml
sed -i '/^  init:$/,/^$/d' Taskfile.yml

echo "Running go mod tidy..."
go mod tidy

rm -rf scripts

echo ""
gum style --bold --foreground 76 "Done!"
echo ""
echo "Next steps:"
echo "  mise exec -- task ci"
echo "  git add -A && git commit -m 'feat: initialize ${BINARY}'"
