# edit

## Description

引数として渡されたファイルを編集できるよう、別のプログラムに渡すコマンドラインツールです。  
設定は edit.json ファイルを編集することで変更することができます。

## Usage

```
edit.exe [Options] [Files...]

Options:
  -l int
        line number(Short)
  -line int
        line number
  -noStdin
        do not input from stdin.
  -version
        Print version information and quit.
```

## Install

To install, use `go get`:

```bash
$ go get -d github.com/rohinomiya/edit
```

## Configuration

See "edit.json".

## Todo

+ Error handling
+ Refactoring

## Contribution

1. Fork ([https://github.com/rohinomiya/edit/fork](https://github.com/rohinomiya/edit/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[rohinomiya](https://github.com/rohinomiya)
