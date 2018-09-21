# Gladius Guardian

Watchdog service for managing the various Gladius processes.

## Docs
Check out our pretty slate docs
[here](https://gladiusio.github.io/gladius-guardian-api-docs/#introduction)

## Service Manager Setup

| Action               | Command                    |
| -------------------- | -------------------------- |
| Install service file | `gladius-guardian install` |
| Start service        | `gladius-guardian start`   |
| Stop   service       | `gladius-guardian stop`    |

**Note for macOS users:** The installed version of the Gladius Guardian service
doesn't use this functionality, it uses a custom service file to run this as a
user service rather than a system one.
