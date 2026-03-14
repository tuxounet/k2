# k2

**A declarative, YAML-driven template engine written in Go.**

k2 lets you define reusable templates, organize them in an inventory, and
generate entire project trees with a single command.
It follows a simple three-phase lifecycle: **Render-plan → Render → Unrender**.

> **Status:** Beta — actively developed, feedback welcome.

---

## Highlights

- **Declarative** — describe *what* you want, not *how* to build it.
- **YAML all the way** — inventory, templates, and applies are plain YAML files.
- **Reusable templates** — define once, apply many times with different variables.
- **Git sources** — pull templates straight from remote Git repositories.
- **Nested templates** — templates can contain other templates.
- **Lifecycle scripts** — hook into bootstrap, pre, post and nuke phases.
- **Idempotent** — plan before you render, unrender when you're done.

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

| Command          | Description                                          |
|------------------|------------------------------------------------------|
| **render-plan**  | Preview the execution plan without changing anything |
| **render**       | Generate files from templates                        |
| **unrender**     | Remove all generated files                           |

### Stack commands (`k2 stack <sub-command>`)

Stack commands orchestrate multi-layer deployments. Each sub-command requires a stack name as first argument.

| Sub-command    | Description                                                         |
|----------------|---------------------------------------------------------------------|
| **up**         | Render templates, then start all layers in order                    |
| **down**       | Stop all layers in reverse order                                    |
| **restart**    | `down` then `up`                                                    |
| **build**      | Run `verbs/build.sh` on all (or a specific) layer                   |
| **exec**       | Run any named verb on every layer that has it                       |
| **status**     | Show status table for each layer                                    |
| **logs**       | Stream logs (all layers or a specific one)                          |
| **healthcheck**| Check health via optional hooks                                     |
| **shell**      | Open an interactive shell in a layer                                |
| **urls**       | Display access URLs from `links.env` files                          |
| **run**        | Execute a verb on a specific layer                                  |
| **list**       | List available stacks with descriptions                             |
| **layers**     | List available layers with verb indicators                          |

```bash
# Preview what will happen
k2 render-plan --inventory ./k2.inventory.yaml

# Render templates
k2 render --inventory ./k2.inventory.yaml

# Clean up generated files
k2 unrender --inventory ./k2.inventory.yaml
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