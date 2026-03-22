# Stacks

## Description

Stacks allow k2 to orchestrate multi-service deployments by defining ordered layers of services. Each stack is a YAML file that declares a sequence of layers, each with its own configuration. All layer operations are delegated to shell script verbs (`verbs/*.sh`).

This feature is inspired by the [k-anissa](https://github.com/tuxounet/k-anissa) project and brings all its capabilities natively into k2.

## Stack YAML Structure

Stack files are located in a `stacks/` directory at the project root (same directory as `k2.inventory.yaml`).

```yaml
version: v0

stack:
  description: "Human-readable description of this stack"
  extends:                           # Optional: inherit layers/env from parent stacks
    - base-stack.yaml                # Paths relative to stacks/ directory

  env:                       # Global environment variables for all layers
    TZ: Europe/Paris
    MY_VAR: my-value

  layers:
    - layer: layers/2.locals-subscriptions   # Path relative to project root
      plan: ollama-embedded                  # Plan directory name
      env:                                   # Per-layer env overrides
        OLLAMA_PORT: "11434"

    - layer: layers/3.proxy
      plan: llm-proxy
      env:
        LITELLM_PORT: "4000"
```

## Stack Extends

A stack can extend one or more parent stacks using the `extends` field (a list of stack file names). Each parent's layers are prepended in order before the child's own layers, and parent environment variables are inherited (later parents and the child take precedence).

```yaml
version: v0

stack:
  description: "Extended stack"
  extends:                          # List of parent stacks (paths relative to stacks/)
    - base-infra.yaml
    - base-services.yaml
  env:
    EXTRA_VAR: value                # Merged with parent env

  layers:
    - layer: layers/additional
      plan: my-extra-service        # Appended after all parent layers
```

### Resolution rules

- **Layers**: parent layers are prepended in list order, then child layers.
- **Env**: parents merged in order (later wins), child env overrides all.
- **Chaining**: each parent can itself extend other stacks.
- **Circular detection**: circular extends chains are detected and rejected.

## Layer Types

k2 automatically detects the layer type by looking for verb scripts:

| Type | Detection | Execution |
|------|-----------|-----------||
| `shell` | Presence of `verbs/up.sh` | `bash verbs/up.sh` / `bash verbs/down.sh` |
| `unknown` | No `verbs/up.sh` | Skipped |

## Layer Directory Structure

```
layers/<N.category>/<plan>/
├── k2.apply.yaml       # Template application (optional, for k2 rendering)
├── verbs/               # Shell script verbs
│   ├── up.sh            # Start the service (required)
│   ├── down.sh          # Stop the service
│   ├── build.sh         # Build the service (optional)
│   ├── status.sh        # Output status (UP, DOWN, DEGRADED...)
│   ├── logs.sh          # Stream logs
│   └── <custom>.sh      # Custom verbs
├── links.env            # Access URLs (label=url per line)
├── defaults.env         # Default environment variables
└── secrets/             # Credentials (should be git-ignored)
```

All operations are delegated to verb scripts. For example, a Docker Compose based layer would have a `verbs/up.sh` that calls `docker compose up -d` internally.

## Environment Variable Cascade

Variables are resolved in order of priority (last wins):

1. `.env` file at project root
2. `stack.env` (global stack variables)
3. `stack.layers[n].env` (per-layer overrides)
4. `defaults.env` in the plan directory

## Hooks

Hooks can be defined in `k2.apply.yaml` under `k2.body.hooks`:

```yaml
k2:
  body:
    hooks:
      pre_start: "echo starting..."
      post_start: "echo started!"
      pre_stop: "echo stopping..."
      post_stop: "echo stopped!"
      pre_build: "echo building..."
      post_build: "echo built!"
      healthcheck: "curl -f http://localhost:8080/health"
```

## CLI Commands

All stack commands are under `k2 stack` (alias: `k2 s`):

### `k2 stack up <name>`

Render k2 templates, then start all layers in order.

```bash
k2 stack up my-stack
k2 stack --debug up my-stack
k2 stack --inventory ./k2.inventory.yaml up my-stack
```

### `k2 stack down <name>`

Stop all layers in reverse order.

### `k2 stack build <name> [layer]`

Render k2 templates, then run `verbs/build.sh` on each layer that has one. Layers without a `build.sh` are skipped.

With an optional layer argument, only that layer is built.

```bash
k2 stack build my-stack
k2 stack build my-stack my-layer
```

Hooks `pre_build` and `post_build` are called before and after `build.sh` if defined in `k2.apply.yaml`.

### `k2 stack restart <name>`

Execute `down` then `up`.

### `k2 stack status <name>`

Display status table of each layer. Status is determined by `verbs/status.sh` if available.

### `k2 stack logs <name> [layer]`

Show logs. Without a layer argument, shows all logs in parallel. With a layer name, shows only that layer's logs.

### `k2 stack healthcheck <name>`

Check health of each layer using optional healthcheck hooks defined in `k2.apply.yaml`.

### `k2 stack urls <name>`

Display all access URLs defined in `links.env` files.

### `k2 stack run <name> <layer> [verb] [args...]`

Execute a specific verb (shell script) on a layer. Without a verb, lists available verbs.

```bash
k2 stack run my-stack my-layer              # list verbs
k2 stack run my-stack my-layer custom-verb   # execute verb
```

### `k2 stack exec <name> <verb> [args...]`

Run a verb on **every layer** of a stack that has it. Layers without the verb script are silently skipped. This differs from `run`, which targets a specific layer.

```bash
k2 stack exec my-stack migrate          # run migrate.sh on all layers that have it
k2 stack exec my-stack deploy -- --env prod  # pass extra args
```

Hooks `pre_<verb>` and `post_<verb>` are called when defined in `k2.apply.yaml`.

### `k2 stack list`

List all available stacks with descriptions.

### `k2 stack layers`

List all available layers with verb indicators:
- `◆` Has verb scripts
- `○` Unknown/empty

## Common Flags

| Flag | Description |
|------|-------------|
| `--inventory <path>` | Path to inventory file (determines project root) |
| `--debug` | Enable debug mode with verbose output |

## Compatibility

The `render` command has an `apply` alias and `unrender` has a `destroy` alias for compatibility with projects using those names.
