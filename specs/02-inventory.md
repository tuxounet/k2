# Inventory File

## Description

The inventory file (`k2.inventory.yaml`) is the entry point for k2. It declares:
- Glob patterns to locate templates and applies
- Global variables shared across all components
- Folders to ignore

## Structure

```yaml
k2:
  metadata:
    id: <unique identifier>
    kind: inventory
  body:
    folders:
      ignore: []            # File patterns to ignore
      templates:            # Glob patterns to find templates
        - templates/**/k2.*.yaml
        - templates/**/k2.*.yml
      applies:              # Glob patterns to find applies
        - services/**/k2.apply.yaml
        - services/**/k2.apply.yml
      stacks: stacks        # Directory containing stack YAML files
    vars:                   # Global variables (string key/value)
      title: my project
```

## Field Details

### `metadata`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | yes | Unique identifier of the inventory |
| `kind` | string | yes | Must be `inventory` |

### `body.folders`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `ignore` | string[] | no | Glob patterns of files to ignore during scanning |
| `templates` | string[] | yes | Glob patterns to locate `k2.template.yaml` files |
| `applies` | string[] | yes | Glob patterns to locate `k2.apply.yaml` files |
| `stacks` | string | no | Relative path to the directory containing stack YAML files |

Glob patterns use the `gobwas/glob` library and support:
- `*`: any character except `/`
- `**`: any path (recursive)
- `{a,b}`: alternatives

### `body.vars`

Global variables of type `map[string]string`. These variables are available throughout the entire project.

## Example (from samples)

```yaml
k2:
  metadata:
    id: k2.cli.sample.inventory
    kind: inventory
  body:
    folders:
      ignore: []
      templates:
        - templates/**/k2.*.yaml
        - templates/**/k2.*.yml
      applies:
        - services/**/k2.apply.yaml
        - services/**/k2.apply.yml
      stacks: stacks
    vars:
      title: k2 cli samples and tests
```

## Behavior

1. Pattern paths are relative to the folder containing the inventory file
2. The `FileStore` recursively traverses the folder and filters files according to the glob patterns
3. Each file found is deserialized and its `path` and `folder` are automatically populated
4. If no inventory file is specified via CLI, k2 looks for `./k2.inventory.yaml` by default
