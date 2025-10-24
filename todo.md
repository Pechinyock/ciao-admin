# TODO

- [ ] Log multi writter. configure options for server.config in order to add\remove log output. for now stdout and file would be enough. by default file is disable, stdout is enabled
- [ ] Implement router main goal is use stdlib if things will go too dificult swap to gin
- [ ] A Route struct. ApiRoutes struct, should have a version, version has to be passed on build. Major version is a part of uri like '/api/v1' minor version is simply inside of struct.
- [ ] Api routes could be enabled\disable by the server config
- [ ] Development Routes. Also enabled\disbled by the server config. Method to make life easier like generate postman colletion, test result view etc
- [ ] Form routes. Usage only for UI, no versoin