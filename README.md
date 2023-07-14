# pi-hole-manager

Pi-Hole Local DNS / DHCP management API

This project provides a RESTful API to manage some of the DHCP/DNS resources on a pi-hole
installation, like:

- Manage static DHCP entries
- Manage local DNS entries
- Manage CNAME aliases

We are not interested on managing the Adlists, Whitelists, Query logs or any other primary
functionality of the pi-hole or the FTL.

## TO-DO

### MVP

- [x] Add authentication / authorization control over the routes
- [x] Improve logging with severity classified messages
- [x] Move routes to a /api/v1 prefix
- [x] Make the system configurable (maybe use Viper?)
- [x] Make the code release-ready (~~set gin properly~~ and remove pretty JSON methods)
- [x] Create the systemd service files
- [ ] Setup a CI pipeline
- [x] Create .deb deployable package (target at least to armv6)
    - [x] .rpm, .apk, ArchLinux packages and tarball archives for all the main architectures (BONUS)
- [x] Create a OpenAI/Swagger documentation
- [ ] Add unit tests

### Phase-2

- [ ] Manage local DNS entries
- [ ] Manage CNAME alias
