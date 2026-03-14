# General Architecture

## Overview

k2 is a declarative template engine written in Go. It generates files from templates using a declarative YAML configuration. The system follows a three-phase cycle: **render-plan**, **render**, **unrender**.

## Core Concepts

### YAML Entities

All k2 configuration files share a common root structure:

```yaml
k2:
  metadata:
    id: <unique identifier>
    kind: <entity type>
  body:
    # content specific to the kind
```

The supported entity types (`kind`) are:

| Kind | File | Description |
|------|---------|-------------|
| `inventory` | `k2.inventory.yaml` | Entry point: declares folders and global variables |
| `template` | `k2.template.yaml` | Definition of a reusable template |
| `template-apply` | `k2.apply.yaml` | Template application instance |

### Metadata

Each entity has metadata:

| Field | Description |
|-------|-------------|
| `id` | Unique identifier of the entity (convention: dot notation) |
| `kind` | Entity type (`inventory`, `template`, `template-apply`) |
| `version` | Version (optional) |
| `path` | Absolute file path (automatically populated) |
| `folder` | Parent folder of the file (automatically populated) |

## Software Layers

```
┌─────────────────────────────────────┐
│         CLI (cmds/)                 │  ← User entry point
│   render-plan │ render │ unrender            │
├─────────────────────────────────────┤
│         Inventory (stores/)         │  ← Workflow orchestration
├─────────────────────────────────────┤
│         ActionPlan (stores/)        │  ← Execution plan with tasks
├──────────┬──────────┬───────────────┤
│ FileStore│Templating│TemplatingStore│  ← YAML reading, resolution, copying
├──────────┴──────────┴───────────────┤
│         Libs (libs/)                │  ← Rendering, commands, logging
├─────────────────────────────────────┤
│         Types (types/)              │  ← Data structures
└─────────────────────────────────────┘
```

### FileStore (`stores/files.go`)

Responsible for discovering and parsing YAML files:
- **Scan**: recursively traverses a directory using glob patterns
- **Unmarshalling**: deserializes YAML files into the corresponding Go types
- Automatically populates `path` and `folder` in metadata

### Inventory (`stores/inventory.go`)

Orchestrates the overall workflow:
1. Loads the inventory file
2. Scans folders looking for templates and applies
3. Creates an `ActionPlan` with the tasks to execute

### ActionPlan (`stores/actions.go`)

Represents an execution plan containing:
- An ordered list of **tasks** (`ActionTask`)
- A registry of **entities** (loaded templates and applies)
- The template **references** to resolve

Task types:

| Type | Description |
|------|-------------|
| `local-resolve` | Resolution of a template from the local inventory |
| `git-resolve` | Resolution of a template from a Git repository |
| `apply` | Application of a template to a component |

### TemplatingStore (`stores/templating.*.go`)

Manages template resolution and application:
- **Inventory resolution**: finds a template by its ID in the inventory
- **Git resolution**: clones/updates a repository and extracts the template
- **Render**: copies template files, performs rendering, creates the `.gitignore`
- **Unrender**: deletes files listed in the `.gitignore`

### Libs (`libs/`)

Shared utilities:
- **rendering.go**: Go template rendering with Sprig functions
- **commands.go**: shell command execution with template rendering
- **logging.go**: formatted output (stdout/stderr)
- **constants.go**: constants (`.refs` for the Git cache)

## Typical Project Structure

```
my-project/
├── k2.inventory.yaml           ← Inventory (entry point)
├── templates/
│   ├── my-template/
│   │   ├── k2.template.yaml   ← Template definition
│   │   ├── README.md          ← Files to copy (with rendering)
│   │   └── ...
├── services/
│   ├── my-service/
│   │   ├── k2.apply.yaml      ← Template application
│   │   ├── .gitignore          ← Automatically generated
│   │   └── README.md           ← Generated from the template
└── .refs/                       ← Git template cache (auto)
```
