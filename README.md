# TibiaData API assets

JSON assets that are almost static for the TibiaData API.

This repo contains tooling that generates the assets json files used by the [tibiadata-api-go](https://github.com/TibiaData/tibiadata-api-go) image.

## What's inside

The generation will be generating new json files on a scheduled interval, providing updates JSON data files used by the container.

There are some details missing from the upstream servers called on by [tibiadata-api-go](https://github.com/TibiaData/tibiadata-api-go) application, so therefore this data provides information.

Example of data that is missing is a list of valid worlds that can have houses/guildhalls or mapping for every houseid for us to provide the town and housetype.

## General information

Tibia is a registered trademark of [CipSoft GmbH](https://www.cipsoft.com/en/). Tibia and all products related to Tibia are copyright by [CipSoft GmbH](https://www.cipsoft.com/en/).

## Credits

- Authors: [Tobias Lindberg](https://github.com/tobiasehlert) – [List of contributors](https://github.com/TibiaData/tibiadata-api-assets/graphs/contributors)
- Distributed under [MIT License](LICENSE)
