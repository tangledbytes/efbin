# efbin - Env to Flag Binary Executor
efbin is an extremely simple and stupid program, all it does is convert environmental variables starting with `EFBIN_` into flags and pass it to the binary pointed by `EFBIN__BIN`. 

## BUT WHY!?
Multiple programs/scripts/makefiles/etc call different programs internally and you may want to tweak the flags that are passed to them without tweaking the code. This is how you can do it using `efbin`:

```bash
# Place efbin in $PATH with name of the binary which efbin is supposed to wrap - for example, docker
# Download efbin to ~/bin/docker where efbin is renamed to docker

export EFBIN__BIN=$(which docker) # Save the path to the actual docker binary
export PATH=~/bin:$PATH # Add ~/bin to the path so that efbin will be invoked whenever `docker` is called `efbin` will be called instead

EFBIN_platform=linux/arm64 docker run -rm busybox uname -a
"Linux 8d11977489ed 5.13.0-41-generic #46-Ubuntu SMP Thu Apr 14 19:28:21 UTC 2022 aarch64 GNU/Linux"

EFBIN_platform=linux/amd64 docker run -rm busybox uname -a
"Linux 8d11977489ed 5.13.0-41-generic #46-Ubuntu SMP Thu Apr 14 19:28:21 UTC 2022 amd64 GNU/Linux"
```

## NOTE
1. `EFBIN__BIN` cannot be empty
2. `EFBIN__BIN` must not refer to itself
3. `EFBIN__*` are for internal use only and will not be parsed