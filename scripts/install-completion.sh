#!/bin/bash
# Install shell completion for pk

set -e

# Detect shell
SHELL_NAME=$(basename "$SHELL")

case "$SHELL_NAME" in
  zsh)
    echo "Installing zsh completion..."
    # Check if Homebrew is available
    if command -v brew &> /dev/null; then
      COMPLETION_DIR="$(brew --prefix)/share/zsh/site-functions"
    else
      # Fallback to first directory in fpath
      COMPLETION_DIR="${fpath[1]}"
      if [ -z "$COMPLETION_DIR" ]; then
        COMPLETION_DIR="$HOME/.zsh/completions"
        mkdir -p "$COMPLETION_DIR"
        echo "Created $COMPLETION_DIR"
        echo "Add this to your ~/.zshrc:"
        echo "  fpath=($COMPLETION_DIR \$fpath)"
      fi
    fi

    pk completion zsh > "$COMPLETION_DIR/_pk"
    echo "✓ Zsh completion installed to $COMPLETION_DIR/_pk"
    echo ""
    echo "Reload your shell:"
    echo "  exec zsh"
    ;;

  bash)
    echo "Installing bash completion..."
    COMPLETION_DIR="$HOME/.bash_completion.d"
    mkdir -p "$COMPLETION_DIR"

    pk completion bash > "$COMPLETION_DIR/pk"
    echo "✓ Bash completion installed to $COMPLETION_DIR/pk"
    echo ""
    echo "Add this to your ~/.bashrc if not already present:"
    echo "  for f in ~/.bash_completion.d/*; do source \"\$f\"; done"
    echo ""
    echo "Then reload:"
    echo "  source ~/.bashrc"
    ;;

  fish)
    echo "Installing fish completion..."
    COMPLETION_DIR="$HOME/.config/fish/completions"
    mkdir -p "$COMPLETION_DIR"

    pk completion fish > "$COMPLETION_DIR/pk.fish"
    echo "✓ Fish completion installed to $COMPLETION_DIR/pk.fish"
    echo ""
    echo "Completions will be available in new shells automatically"
    ;;

  *)
    echo "Unsupported shell: $SHELL_NAME"
    echo ""
    echo "Manual installation:"
    echo "  pk completion zsh   # For zsh"
    echo "  pk completion bash  # For bash"
    echo "  pk completion fish  # For fish"
    exit 1
    ;;
esac

echo ""
echo "Test completion by typing:"
echo "  pk session <TAB>"
echo "  pk list <TAB>"
