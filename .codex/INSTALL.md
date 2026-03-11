# Installing tidewave_cli skill for Codex

Enable native skill discovery in Codex with one directory symlink.

## Prerequisites

- Git
- Codex

## Installation

1. Clone the repository:
   ```bash
   git clone <repo-url> ~/.codex/tidewave_cli
   ```

2. Create the skills symlink:
   ```bash
   mkdir -p ~/.codex/skills
   ln -s ~/.codex/tidewave_cli/skills ~/.codex/skills/tidewave
   ```

   If the link already exists and you want to replace it:
   ```bash
   ln -sfn ~/.codex/tidewave_cli/skills ~/.codex/skills/tidewave
   ```

3. Restart Codex to refresh skill discovery.

## Verify

```bash
ls -la ~/.codex/skills/tidewave
```

You should see a symlink pointing to `~/.codex/tidewave_cli/skills`.

## Updating

```bash
cd ~/.codex/tidewave_cli && git pull
```

Skills update through the symlink path.

## Uninstall

```bash
rm ~/.codex/skills/tidewave
```
