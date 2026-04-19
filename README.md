portwatch

Lightweight daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a configuration file:

```bash
portwatch --config config.yaml
```

Example `config.yaml`:

```yaml
interval: 30s
alert:
  method: log
  output: /var/log/portwatch.log
whitelist:
  - 22
  - 80
  - 443
```

portwatch will scan open ports at the defined interval and alert whenever a port outside the whitelist is opened or an expected port disappears.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to config file |
| `--interval` | `30s` | Scan interval |
| `--verbose` | `false` | Enable verbose logging |

### Example Alert Output

```
[ALERT] 2024-01-15T10:23:01Z - Unexpected port opened: 8080 (TCP)
[ALERT] 2024-01-15T10:23:01Z - Expected port closed: 443 (TCP)
```

## License

MIT © 2024 yourusername