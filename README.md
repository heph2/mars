![Logo of the project](assets/gopher-astronaut_dribbble-200x200.png)
# Mars
> Mars is a simple http backend for terraform.

Provide a basic HTTP RESTFul API for managing terraform states. 
Supporting lock, state encryption and multiple backend for storage.

## Usage
For the full documentation, read the mars(1) man page.

## Installing

### Compiling from source
If your system has a supported version of Go, you can easily build from source.
```
go install git.sr.ht/~heph/mars@latest
```
### Container
There's also a docker container available
```
docker run --rm -d --name mars --network="host" -p 8080:8080
```
### Nix
The best way to use mars is through his nix derivation in my NUR repository.
I'd also provide a basic nix module.

## Contributing
Send patches and issue reports to [srht@mrkeebs](mailto:srht@mrkeebs.eu)

## License
`mars` is Copyright @ 2022 heph.

`mars` is licensed under the terms of the GNU Affero General Public License, version3. For more information, see [LICENSE][] of the [GNU website][agpl-3]

[LICENSE]: LICENSE
[agpl-3]: https://www.gnu.org/licenses/agpl-3.0.standalone.html
