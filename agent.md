# Project Information

This project is named `hcli`,
built with the `Golang` programming language.
It is used to quickly create `Hugo` posts and optimize post content via command line.

After the code update, the README.md should be updated.

The project module was `github.io/uberate/hcli`

# Feature List

This section describes all features of `hcli` and provides helpful tools for MCP servers.

## Main behavior

1. `hcli config` provides a configuration command to display configuration and debug commands for the config file.
    1. `hcli config demo` to show the demo config file. 
2. `hcli gen` provides commands to create files quickly:
    1. `hcli gen posts -n|--template-name template-name file-name`: creates a new post using a specified template (
       defined in the hcli config).
        1. `-n|--template-name` defines the template name for the post to be created.
           Template files define the following information:
            1. Name: the template name
            2. Categories: the post categories
            3. NeedDir: whether the post needs a new directory (pictures are only effective when NeedDir is true)
            4. Tags: the template tags
            5. Template: the post template (hcli will render the template)
    2. `hcli gen pic` creates an article poster with a specified file.
    3. `hcli optimize posts` optimizes file content. Hcli optimizes posts with internal prompts, and users can override this behavior through configuration fields in templates.
## MCP Server
1. `hcli mcp start` starts the MCP server.

## AI Case
Hcli provides an AI switcher. Please try to abstract the AI framework.

## Hcli Configuration

Hcli has a base configuration.

Hcli configuration example:

```yaml
Templates:
  - Name: sa # Hcli gen posts -nsa xxx
    Categories: [ "Read" ]
    Tags: [ "short_reads" ]
    Templates: |
      +++
      date = '{{.createAt}}'
      title = '{{.title}}'
      categories = {{.categories}}
      tags = {{.tags}}
      +++
      Some contents...
AI:
  Provider: "volc" 
```

# Code Format

All agents must follow this format in this project.

## Go Code Format

This section describes the Go code format.

All the file must has the unit test file.

### Bootstrap files

1. The `hcli` command is built using the `github.com/spf13/cobra` project.

### Behaviors

1. The `hcli` project supports standard output/error handling.

## Markdown Code Format

1. A newline character must follow every punctuation mark.

### Docs

`README.md` describes the basic information of the project. It should only include these
sections, all other section should be dropped.
1. A short description of the project.
2. How to build and basic usage of main behaviors.
3. Feature plans and status.

All detailed documentation is defined in the `${PROJECT_HOME}/docs/` directory.

The documentation provides multiple language support.

## Project Layout

Project layout describes how to manage project files.

1. All commands are defined in the path: `${PROJECT_HOME}/cmd/cli/root/`.
2. All program bootstrap files must be in: `${PROJECT_HOME}/cmd/`, and each bootstrap is defined in a unique directory.
3. Library code is defined in `${PROJECT_HOME}/pkg/`, and any module that does not belong to the project's main behavior
   is considered a library.
4. All module doc should define in the mod.go file. All the `Interface` define should in 'type.go'.

# Project tools:
1. used `make test` to test the project.
2. used `make releases` to create all platform bins.

# Notices

- `${PROJECT_HOME}` represents the project root path.