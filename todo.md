# TODO

- [X] Log multi writter. configure options for server.config in order to add\remove log output. for now stdout and file would be enough. by default file is disable to enable provide log path in config file
- [ ] Every day new log file if file log is enabled. Come up with ability to rotate log. Rotation should be configurable with options like: delete, move to(filePath), compress() etc
- [X] Implement router main goal is use stdlib if things will go too dificult swap to gin
- [X] A Route struct. ApiRoutes struct
- [X] Api routes could be enabled\disable by the server config
- [ ] Development Routes. Also enabled\disbled by the server config. Method to make life easier like generate postman colletion, test result view etc
- [X] Form routes. Usage only for UI, no versoin
- [ ] there's several 'loadConfig' func, but could be one using genercs. Inmplement single func for all configs
- [ ] there's several func for register endpoints. could be one...
- [ ] if runs via https add redirect middleware for http requests