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

 ## Contributing

Send patches and issue reports to [srht@mrkeebs](mailto:srht@mrkeebs.eu)

## License

`mars` is Copyright @ 2022 heph.

`mars` is licensed under the terms of the GNU Affero General Public License, version3. For more information, see [LICENSE][] of the [GNU website][agpl-3]

[LICENSE]: LICENSE
[agpl-3]: https://www.gnu.org/licenses/agpl-3.0.standalone.html
