# Nested Templates

## Description

A template can itself contain a `k2.apply.yaml` file in its subfolders. When the parent template is applied, these files are discovered and included in the execution plan. This allows composing templates from other templates.

## Mechanism

### Discovery

`k2.apply.yaml` files within a template are discovered by the inventory scan, since the glob patterns recursively traverse all folders.

For example, with the pattern `templates/**/k2.*.yaml` and the structure:

```
templates/
└── kind2/
    ├── k2.template.yaml
    ├── README.md
    └── sub/
        └── k2.apply.yaml     ← Nested apply
```

The file `sub/k2.apply.yaml` is found by the pattern `templates/**/k2.*.yaml` and included as a standard apply. It will be treated as a standalone apply during scanning, resolving its own template.

### Resolution

The nested apply works exactly like a top-level apply:
- It references a template (here `kind1`)
- Its variables are independent
- Its scripts are executed normally

## Example (from samples)

### Structure

```
templates/kind2/
├── k2.template.yaml           ← Template kind2
├── README.md
└── sub/
    └── k2.apply.yaml          ← Nested apply that uses kind1
```

### Parent template (`kind2`)

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

### Nested apply (`kind2/sub`)

```yaml
k2:
  metadata:
    id: k2.cli.sample.templates.kind2.sub
    kind: template-apply
  body:
    template:
      source: inventory
      params:
        id: k2.cli.sample.templates.kind1
    vars:
      name: "sub-kind2"
      description: "Sub Template of type kind2"
```

### Result

When a component uses the `kind2` template, the `sub/` subfolder is copied with the `k2.apply.yaml` file. The inventory scan then discovers this apply and includes it in the execution plan, applying the `kind1` template in the subfolder.

## Implications

- Nested templates increase **composability**: a complex template can be assembled from basic building blocks
- **Deduplication** applies: if the same base template is used multiple times (directly and via nesting), its resolution is performed only once
- **Scripts** at each level are executed independently
