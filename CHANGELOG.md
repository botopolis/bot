# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased](https://github.com/botopolis/bot/compare/v0.5.0...master)

## [v0.5.0](https://github.com/botopolis/bot/compare/v0.4.2...v0.5.0)

### Added

- `help` plugin to show help text for installed plugins
- `Robot.ListPlugins` to list installed plugins.

## [0.4.2](https://github.com/botopolis/bot/compare/v0.4.1...v0.4.2) - 2018-07-10

### Changed

- `Robot.Respond` now responds to usernames when they have an `@` in front of them

## [0.4.1](https://github.com/botopolis/bot/compare/v0.4.0...v0.4.1) - 2018-06-25

### Added

- `Brain.Keys()` returns alphabetized cache keys

## [0.4.0](https://github.com/botopolis/bot/compare/v0.3.0...v0.4.0) - 2018-06-24

### Added

- Mock that satisfies `bot.Logger` interface

### Changed

- `Robot.Logger` is now an interface.

### Removed

- `Robot.Debug()` has been removed. As you now have full control over the logger, you can set your logger's log level without involving `bot.Robot`.
