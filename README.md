# k2

**A declarative, YAML-driven template engine written in Go.**

k2 lets you define reusable templates, organize them in an inventory, and
generate entire project trees with a single command.
It follows a simple three-phase lifecycle: **Plan → Apply → Destroy**.

> **Status:** Beta — actively developed, feedback welcome.

---

## Highlights

- **Declarative** — describe *what* you want, not *how* to build it.
- **YAML all the way** — inventory, templates, and applies are plain YAML files.
- **Reusable templates** — define once, apply many times with different variables.
- **Git sources** — pull templates straight from remote Git repositories.
- **Nested templates** — templates can contain other templates.
- **Lifecycle scripts** — hook into bootstrap, pre, post and nuke phases.
- **Idempotent** — plan before you apply, destroy when you're done.

---

## Quick Start

### Install

```bash
go install github.com/tuxounet/k2@latest
```

### Verify

```bash
k2 help
```

### Uninstall

```bash
rm "$(go env GOPATH)/bin/k2"
```

---

## Usage

k2 exposes three commands. Each accepts an optional `--inventory` flag
(defaults to `./k2.inventory.yaml`).

| Command     | Description                                          |
|-------------|------------------------------------------------------|
| **plan**    | Preview the execution plan without changing anything |
| **apply**   | Generate files from templates                        |
| **destroy** | Remove all generated files                           |

```bash
# Preview what will happen
k2 plan --inventory ./k2.inventory.yaml

# Apply templates
k2 apply --inventory ./k2.inventory.yaml

# Clean up generated files
k2 destroy --inventory ./k2.inventory.yaml
```

---

## Development

```bash
# Clone & enter the repo
git clone https://github.com/tuxounet/k2.git && cd k2

# See all available Make targets
make help

# Run tests
make test

# Build the binary
make build
```

---

## Documentation

Full technical specifications live in the [specs/](specs/) folder — covering
architecture, inventory format, template definitions, variables & rendering,
lifecycle scripts, and more.

---

## License

See [LICENSE](LICENSE) for details.

## Author

**Krux** — [github.com/tuxounet](https://github.com/tuxounet)