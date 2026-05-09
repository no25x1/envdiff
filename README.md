# envdiff

> Diff and reconcile `.env` files across environments with secret masking.

---

## Installation

```bash
go install github.com/youruser/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/youruser/envdiff.git && cd envdiff && go build -o envdiff .
```

---

## Usage

Compare two `.env` files and mask sensitive values:

```bash
envdiff .env.development .env.production
```

**Example output:**

```
~ DATABASE_URL   dev=postgres://localhost/app  prod=postgres://****@prod-db/app
+ NEW_FEATURE_FLAG                             prod=true
- DEBUG                dev=true
```

### Flags

| Flag | Description |
|------|-------------|
| `--mask` | Mask secret values in output (default: true) |
| `--export` | Export reconciled `.env` to a file |
| `--format` | Output format: `text`, `json`, `table` |

### Reconcile

Merge differences and write a reconciled output file:

```bash
envdiff --export .env.reconciled .env.staging .env.production
```

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © youruser