# Gitignore and Dropbox
Ignore directory, files in Dropbox based on `.gitignore`.

## Installation
Build from source or download a binary from the Release

## Usage
Approach 1: Saving the binary in the root of the Dropbox directory. It will scan all of the `.gitignore` files recursively.

Approach 2: Save the binary anywhere you like and call it with the options

## Options
| Flag           | Note
| -------------- | ---------------------------------------------------------- |
| -h, --help     | display help information                                   |
| -p, --path     | path to your Dropbox or a folder which you want to scan    |
| -d, --dry-run  | dry run                                                    |
| -u, --unignore | unignore                                                   |

## Note
- Tested on Windows, Linux
- Cannot test on macOS
- Does not support negative pattern (!)