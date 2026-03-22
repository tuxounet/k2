# Template Application (template-apply)

## Description

An application file (`k2.apply.yaml`) instantiates a template in a target folder. It declares:
- The reference to the source template (inventory or git)
- Instance-specific variables (which override template parameters)
- Apply-specific lifecycle scripts

## Structure

```yaml
k2:
  metadata:
    id: <unique identifier>
    kind: template-apply
  body:
    template:
      source: inventory|git    # Template source
      params:                   # Resolution parameters
        id: <template-id>      # For source=inventory
        # OR
        repository: <url>      # For source=git
        branch: <branch>
        path: <path>
    vars:                       # Instance variables
      name: "my-component"
      description: "..."
    scripts:                    # Apply-specific scripts
      bootstrap: [...]
      pre: [...]
      post: [...]
      nuke: [...]
```

## Field Details

### `metadata`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | yes | Unique identifier of the apply |
| `kind` | string | yes | Must be `template-apply` |

### `body.template`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `source` | string | yes | `inventory` or `git` |
| `params` | map[string]string | yes | Template resolution parameters |

See [05-template-sources.md](05-template-sources.md) for source details.

### `body.vars`

Instance variables of type `map[string]any`. Supports:
- Strings
- Lists (arrays)
- Nested objects (maps)

These variables **override** the `parameters` defined in the template. The merge mechanism uses `MergeMaps()`: apply variables overwrite template variables for matching keys.

### `body.scripts`

Apply lifecycle scripts. Same structure as template scripts. Apply scripts are executed **after** template scripts at each phase.

## Examples (from samples)

### Simple apply with inventory source

```yaml
k2:
  metadata:
    id: k2.cli.sample.services.product1.component2bis
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: k2.cli.sample.templates.kind1
    vars:
      name: "component2bis"
      description: "Template of type kind1"
```

### Apply with complex variables

```yaml
k2:
  metadata:
    id: k2.cli.sample.services.product1.component2
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: k2.cli.sample.templates.kind1
    vars:
      name: "component2"
      description: "Template of type kind1"
      coll:
        - a
        - b
        - c
      obj:
        a: 1
        b: stc
        c: true
    scripts:
      nuke:
        - echo "Nuke script for component2"
```

### Apply with bootstrap and post scripts

```yaml
k2:
  metadata:
    id: k2.cli.sample.services.product1.component1
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: k2.cli.sample.templates.kind2
    vars:
      name: "component2"
      description: "Template of type kind1"
    scripts:
      bootstrap:
        - echo "boot {{ .name }}"
      post:
        - echo "finjj"
```

### Apply with Git source

```yaml
k2:
  metadata:
    id: k2.cli.sample.services.product2.componentFromGit1
    kind: template-apply
  body:
    template:
      source: git
      params:
        repository: https://github.com/tuxounet/k2.git
        branch: main
        path: samples/templates/fromGit1/k2.template.yaml
    scripts:
      bootstrap:
        - echo "boot"
      post:
        - echo "fin"
    vars:
      name: "component2"
      description: "Template of type kind1"
```

## Apply Behavior

1. The referenced template is resolved (from inventory or Git)
2. Variables are merged: `template.parameters` + `apply.vars` (apply takes priority)
3. The target folder's `.gitignore` is checked to determine if this is a first apply
4. Scripts execute in order: bootstrap (if first apply), pre, file copying, post
5. A `.gitignore` file is generated listing all produced files
6. The `k2.apply.yaml` file is explicitly excluded from `.gitignore` via `!k2.apply.yaml`

## Target Folder

The target folder is the folder containing the `k2.apply.yaml` file. Template files are copied there with Go template expression rendering.
