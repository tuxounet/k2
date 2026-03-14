# Stacks Listing

## Description

The `k2 stacks` command lists all available stacks defined in the project. The stacks directory path is read from the inventory file (`k2.inventory.yaml`), under `body.folders.stacks`.

This command provides a quick overview of all stacks without requiring a stack name argument.

## Inventory Configuration

The stacks folder must be declared in the inventory:

```yaml
k2:
  metadata:
    id: my-project
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.*.yaml
      applies:
        - services/**/k2.apply.yaml
      stacks: stacks          # <-- stacks directory (relative to inventory)
    vars:
      title: my project
```

### `body.folders.stacks`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `stacks` | string | no | Relative path to the directory containing stack YAML files |

If the `stacks` field is omitted or empty, the `k2 stacks` command returns an error indicating that no stacks folder is configured.

## CLI Usage

```bash
k2 stacks [--inventory <path>]
```

### Aliases

- `k2 ls`

### Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--inventory` | string | `./k2.inventory.yaml` | Path to the inventory file |

## Behavior

1. Loads the inventory file (default: `./k2.inventory.yaml`)
2. Reads the `body.folders.stacks` field to determine the stacks directory
3. Scans the stacks directory for `.yaml` and `.yml` files
4. For each stack file found:
   - Extracts the stack name (filename without extension)
   - Parses the YAML to read `stack.description` and count `stack.layers`
5. Displays the list with name, description, and layer count

### Output Format

```
  Stacks disponibles :
  ────────────────────────────────────────
  dev  — Development stack (2 layers)
  prod  — Production stack (1 layers)
  minimal  (0 layers)
```

### Edge Cases

| Scenario | Behavior |
|----------|----------|
| No `stacks` field in inventory | Returns error: "no stacks folder defined in inventory" |
| Stacks directory does not exist | Returns error: "cannot read stacks directory" |
| Empty stacks directory | Displays "Aucune stack trouvée" |
| Non-YAML files in stacks dir | Ignored |
| Subdirectories in stacks dir | Ignored |
| Malformed YAML stack file | Name displayed without description |

## Difference with `k2 stack list`

| | `k2 stacks` | `k2 stack list` |
|---|---|---|
| Stacks directory | From inventory (`body.folders.stacks`) | Hardcoded `stacks/` relative to `--inventory` dir |
| Layer count | Displayed | Not displayed |
| Requires inventory | Yes | No |

## Example

```bash
# List stacks using default inventory
k2 stacks

# List stacks using a specific inventory
k2 stacks --inventory ./path/to/k2.inventory.yaml
```
