# buf

A tiny command-line buffer for your shell.

```sh
# Save output:
$ printf 'apple\nbanana\ncherry\n' | buf
apple
banana
cherry

# And pipe it to other commands:
$ buf | wc -l
3

# The buffer stays around until replaced:
$ buf
apple
banana
cherry
```

## Install

```sh
go install github.com/oosawy/buf@latest
```

## Usage

### `command | buf`

Save the output of a command to the buffer:

```sh
$ printf 'apple\nbanana\ncherry\n' | buf
```

When stdin is not a tty, `buf` saves it to the buffer and passes it through unchanged.

### `buf`

Print the contents of the buffer:

```sh
$ buf
apple
banana
cherry
```

When stdin is a tty, `buf` prints the saved buffer instead.

### `buf -- command`

Transform it and replace the buffer:

```sh
$ buf -- sort
apple
banana
cherry

$ buf -- tail -n 1
cherry
```

With `--`, `buf` feeds the saved buffer to the command and saves its output.

`buf -- command` is equivalent to `buf | command | buf`.

## Notes

- `command | buf` saves the input and passes it through unchanged.
- Each shell has its own buffer; saving replaces the previous one.
- The buffer lives in the system temp directory and is cleared on reboot.
