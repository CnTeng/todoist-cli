# todoist-cli

[![Go][go-shield]][go-url]

A powerful CLI client for Todoist.

## ✨ Features

#### Item Management

- Built on [todoist-go-api](https://github.com/CnTeng/todoist-api-go).
- Supports management of `task`, `project`, `section`, `label` and `filter`.

#### Pretty Table

- Built on [table](https://github.com/CnTeng/table), inspired by `taskwarrior`.
- Format text as **Bold**, _Italic_ or ~~Strikethrough~~.
- Automatically wraps lines.
- Colorful text and Nerd Font icons.
- Displays subtasks in a tree structure.

![table](https://github.com/user-attachments/assets/f9ea04c3-1435-45e8-b8d5-84e0b294d780)

#### Autocompletion

- Built on [cobra](https://github.com/spf13/cobra)
- Autocompletion for `task`, `project`, `section`, `color`, and more.
- Supports `bash`, `zsh`, `fish`, and `powershell`.

![cmp](https://github.com/user-attachments/assets/8bcc5ed0-a691-493d-a9d2-2c066c39b7b7)

#### Reorder

- Reorder items interactively, similar to `git rebase -i`.

![reorder](https://github.com/user-attachments/assets/a212e856-5960-46b9-9112-5f432f453f02)

#### Auto sync with Todoist

- Keeps your local data in sync with Todoist using a background daemon.
- Changes made in the Todoist web or mobile app are reflected in the CLI automatically.
- You can also manually sync using `todoist sync`.

## 🛠️ Usage

### Help

```
A CLI for Todoist

Usage:
  todoist [command]

Task commands:
  add         Add task
  close       Close task
  list        List tasks
  modify      Modify task
  move        Move task
  quick-add   Quick add task
  remove      Remove task
  reopen      Reopen task
  reorder     Reorder tasks

Resources commands:
  filter      Filter commands
  label       Label commands
  project     Project commands
  section     Section commands

Additional Commands:
  completion  Generate the autocompletion script for the specified shell
  daemon      Start daemon
  help        Help about any command
  sync        Sync data

Flags:
      --config string   config file (default "/home/user/.config/todoist/config.toml")
  -h, --help            help for todoist

Use "todoist [command] --help" for more information about a command.
```

### Nix

Install with home-manager:

```nix
{
  programs.todoist-cli = {
    enable = true;
    settings = {
      daemon.api_token_file = "/run/secrets/your-todoist-api-token";
    };
  };
}
```

### Other

Install via Go:

```bash
go install github.com/CnTeng/todoist-cli@latest
```

## 🔧 Configuration

### Available Options

`daemon.api_token` or `daemon.api_token_file` is the only required option.

```toml
[daemon]
## Daemon address. "@todo.sock" is the default value.
# Address = "@todo.sock"

## Set your Todoist API token. Either `api_token` or `api_token_file` is required.
# api_token = "your-todoist-api-token"
api_token_file = "/run/secrets/todoist/token"

[icon]
## Icon theme for todoist-cli. Supported values: `nerd` (default), `text`.
# default = "nerd"

## You can set your own icons here or override the default icons.
# none = "  "
# done = " "
# undone = " "
# inbox = " "
# favorite = " "
# indent = "│ "
# last_indent = "└ "
```

### Override Configuration by Environment Variables

| Configuration         | Environment            |
| --------------------- | ---------------------- |
| daemon.address        | TODOIST_ADDRESS        |
| daemon.api_token      | TODOIST_API_TOKEN      |
| daemon.api_token_file | TODOIST_API_TOKEN_FILE |

## ⭐ Credit

- inspired by [todoist](https://github.com/sachaos/todoist)

[go-shield]: https://img.shields.io/github/go-mod/go-version/CnTeng/todoist-cli?style=for-the-badge&logo=go
[go-url]: https://golang.org
