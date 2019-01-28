# botopolis/help

[![GoDoc](https://godoc.org/github.com/botopolis/bot/help?status.svg)](https://godoc.org/github.com/botopolis/bot/help) [![Build Status](https://circleci.com/gh/botopolis/bot.svg?style=svg)](https://circleci.com/gh/botopolis/bot) [![Test Coverage](https://api.codeclimate.com/v1/badges/b7acc61121363e7405a3/test_coverage)](https://codeclimate.com/github/botopolis/bot/test_coverage)

Output about your botopolis commands.

## Usage

See [example_test.go](./example_test.go) for usage details.

In order to show what commands your plugin implements, simply implement the
[`help.Provider`](https://godoc.org/github.com/botopolis/bot/help#Provider)
interface (`Help() []help.Text`). `help.Plugin` lists the help text for all
installed plugins.

### Example output

```
user: @bot help
bot:
  help - Displays all the help commands that this bot knows about.
  help <query> - Displays all help commands that match <query>.
```

