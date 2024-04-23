# tmpmod

tmpmod is a tool for temporary use of modified Go modules.

## Usage

First, submit a pull request to the upstream repository for modification.

So, you have the modified branch. Now use `tmpmod`.

### 1. Create a renamed branch ( `tmpmod switch` )

Create a branch with the module renamed to the specified ( `--as` ) name.

Temporarily use a patched module by pushing the created branch and importing it.

``` console
# /src/github.com/supercool/greatmodule (fix-something)> tmpmod switch --as github.com/k1low/greatmodule
Switching to renamed-github.com/k1low/greatmodule-by-tmpmod...
Renaming module to github.com/k1low/greatmodule...
Committed

Usage: push renamed-github.com/k1low/greatmodule-by-tmpmod and use `go get github.com/k1low/greatmodule@f120013f64dca79ae9da1978ac6a54d780bb98e7`
# /src/github.com/supercool/greatmodule (renamed-github.com/k1low/greatmodule-by-tmpmod)> git push k1low renamed-github.com/k1low/greatmodule-by-tmpmod
```

### 2. Retrieve source code from a branch and rename it ( `tmpmod get` )

Retrieve the source code from the specified branch and rename it.

``` console
# /src/github.com/k1LoW/myproject (main)> tmpmod get github.com/supercool/greatmodule@fix-something
Getting github.com/supercool/greatmodule@fix-something...
Renaming module to github.com/k1LoW/myproject/tmpmod/github.com/supercool/greatmodule...

Usage: use `github.com/k1LoW/myproject/tmpmod/github.com/supercool/greatmodule`
# /src/github.com/k1LoW/myproject (main)>
```

#### `--rename-all`

If you want to rename also the module path in the importing source codes, use `--rename-all`.

#### Revert renamed module ( `tmpmod revert` )

``` console
# /src/github.com/k1LoW/myproject (main)> tmpmod revert tmpmod/github.com/supercool/greatmodule
Getting github.com/supercool/greatmodule@fix-something...
Reverting module from github.com/k1LoW/myproject/tmpmod/github.com/supercool/greatmodule to github.com/supercool/greatmodule...

Reverted
# /src/github.com/k1LoW/myproject (main)>
```

## Install

**homebrew tap:**

```console
$ brew install k1LoW/tap/tmpmod
```

**manually:**

Download binary from [releases page](https://github.com/k1LoW/tmpmod/releases)

**go install:**

```console
$ go install github.com/k1LoW/tmpmod@latest
```

**deb:**

``` console
$ export TMPMOD_VERSION=X.X.X
$ curl -o tmpmod.deb -L https://github.com/k1LoW/tmpmod/releases/download/v$TMPMOD_VERSION/tmpmod_$TMPMOD_VERSION-1_amd64.deb
$ dpkg -i tmpmod.deb
```

**RPM:**

``` console
$ export TMPMOD_VERSION=X.X.X
$ yum install https://github.com/k1LoW/tmpmod/releases/download/v$TMPMOD_VERSION/tmpmod_$TMPMOD_VERSION-1_amd64.rpm
```

**apk:**

``` console
$ export TMPMOD_VERSION=X.X.X
$ curl -o tmpmod.apk -L https://github.com/k1LoW/tmpmod/releases/download/v$TMPMOD_VERSION/tmpmod_$TMPMOD_VERSION-1_amd64.apk
$ apk add tmpmod.apk
```

## Why not use `replace`?

Because modules that use `replace` can't `go install`.
