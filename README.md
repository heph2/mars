# Debris
Debris is a simple http backend for terraform.

# Description
debris(1) is a http backend for terraform, supporting multiple
backend for storing the statefile, including sqlite3 and plain
filesystems. Supporting lock.

# Rationale
The most common http backend for terraform is the one built-in to
gitlab, which it's okay, but i'm not a big fan of gitlab and i do not
want myself tighted to a forge.

# Installing
You can easily compile yourself from this git repository, or
use the docker container provided. Or better, if you're a nix user
you can use the module available on my NUR
