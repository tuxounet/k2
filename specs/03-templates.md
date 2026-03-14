# Template Definition

## Description

A template (`k2.template.yaml`) defines a reusable file model. It contains:
- Default parameters (template variables)
- Lifecycle scripts
- Co-located files that will be copied and rendered during apply

## Structure

```yaml
k2:
  metadata:
    id: <unique identifier>
    kind: template
  body:
    name: <template name>
    parameters:              # Default variable values
      name: "default-name"
      description: "default description"
    scripts:                 # Lifecycle scripts (optional)
      bootstrap:
        - echo "initialization"
      pre:
        - echo "before copy"
      post:
        - echo "after copy"
      nuke:
        - echo "destruction"
```

## Field Details

### `metadata`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | string | yes | Unique identifier of the template |
| `kind` | string | yes | Must be `template` |

### `body`

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | yes | Human-readable name of the template |
| `parameters` | map[string]any | no | Default variable values for the template |
| `scripts` | object | no | Lifecycle scripts |

### `body.scripts`

| Phase | Description |
|-------|-------------|
| `bootstrap` | Executed only once, during the first render |
| `pre` | Executed before file copying, on each render |
| `post` | Executed after file copying, on each render |
| `nuke` | Executed during unrender |

See [07-scripts-lifecycle.md](07-scripts-lifecycle.md) for script details.

## Template Files

All files present in the same folder as `k2.template.yaml` (and its subfolders) are considered template files, except for:
- `k2.template.yaml` itself
- `.DS_Store`
- The `.git` folder and its contents

These files are copied to the target folder during render, with Go template expression rendering.

## Examples (from samples)

### Simple template (`kind2`)

```yaml
k2:
  metadata:
    id: k2.cli.sample.templates.kind2
    kind: template
  body:
    name: kind2
    parameters:
      name: kind2
      description: "Template of type kind2"
```

Associated `README.md` file:
```markdown
# {{ .name }}

{{ .description }}
```

### Template with scripts (`kind1`)

```yaml
k2:
  metadata:
    id: k2.cli.sample.templates.kind1
    kind: template
  body:
    name: kind1
    parameters:
      name: kind1
      description: "Template of type kind1"
    scripts:
      bootstrap:
        - echo "template boot of {{ .name }}"
      post:
        - echo "template fin of {{ .name }}"
```

Associated `README.md` file with iteration and object access:
```markdown
# {{ .name }}

{{ .description }}

## colls:
{{ range .coll }}
-  {{ . }}{{ end }}

## map data

{{ .obj.b }}
```

### Template for Git source (`fromGit1`)

```yaml
k2:
  metadata:
    id: k2.cli.sample.templates.fromGit1
    kind: template
  body:
    name: fromGit1
    parameters:
      name: fromGit1
      description: "Template of type fromGit1"
```

### Template with ERB placeholders (`with-placeholder`)

```yaml
k2:
  metadata:
    id: k2.cli.sample.templates.with-placeholder
    kind: template
  body:
    name: with-placeholder
    parameters:
      name: with-placeholder
      description: "Template with-placeholder"
```

Associated `README.md` file using ERB-like syntax:
```markdown
# <%= title %>
## <%= name %>

<%= description %>
```

> **Note**: The `<%= ... %>` syntax is not standard Go template. It is rendered by the rendering engine if it matches a syntax supported by Sprig functions, or can be used as a convention in non-rendered files.
