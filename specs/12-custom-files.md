# Custom Files

## Description

An apply folder can contain additional files that do not come from the template. These custom files coexist with the generated files and are preserved during render and unrender operations.

## Mechanism

### During Render

Files already present in the target folder that are not part of the template are neither modified nor deleted. Render only affects:
- Files copied from the template (overwritten on each render)
- The `.gitignore` (regenerated)

### During Unrender

Only files listed in the `.gitignore` are deleted. Custom files, not being referenced in the `.gitignore`, are preserved.

### Protection in `.gitignore`

The `k2.apply.yaml` file is explicitly protected via the `!k2.apply.yaml` entry in the generated `.gitignore`.

## Example (from samples)

### Structure

```
services/product2/withCustomFiles/
├── k2.apply.yaml              ← Apply using the with-placeholder template
└── files/
    └── MMM.md                 ← Custom file (not in the template)
```

### Apply File

```yaml
k2:
  metadata:
    id: k2.cli.sample.services.product2.withCustomFiles
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: k2.cli.sample.templates.with-placeholder
    vars:
      name: "component2"
      description: "Template of type with-placeholder"
```

### Custom File

The file `files/MMM.md` contains content specific to this component:

```
avec du contneu
```

### Behavior

1. **render**: template files from `with-placeholder` (e.g., `README.md`) are copied and rendered. The `files/` folder and `MMM.md` are not touched.
2. **unrender**: only files listed in the `.gitignore` (e.g., `README.md`) are deleted. `files/MMM.md` is preserved.

## Use Cases

Custom files allow you to:
- Add configuration specific to a component
- Store additional resources
- Maintain manual files alongside generated files
- Enrich a template with local non-generic data
