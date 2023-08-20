<p align="center">
  <a href="https://go.dev"><img src="https://img.shields.io/badge/go-white?logo=go&style=for-the-badge"></a> 
  <a href="https://github.com/vareversat/chabo-api/actions"><img src="https://img.shields.io/github/actions/workflow/status/vareversat/chabo-api/dev.yaml?logo=github&style=for-the-badge"></a>
  <a href="https://github.com/vareversat/chabo-api/releases"><img src="https://img.shields.io/github/v/tag/vareversat/chabo-api?label=version&logo=git&logoColor=white&style=for-the-badge"></a>
  <a href="https://codecov.io/gh/vareversat/chabo-api/"><img src="https://img.shields.io/codecov/c/github/vareversat/chabo-api?logo=codecov&style=for-the-badge&token=97YDVRS0X4"></a>
</p>

# Chabo API

**The place to get the Chaban Delmas event schedules !**
This REST API is entended to improve the already existing Open Data API provided by Bordeaux Métropole available [here](https://opendata.bordeaux-metropole.fr/explore/dataset/previsions_pont_chaban/information/).

## Overview

This API allows you to get the schedules you want by filtering over the type of closing, the closing date or even the boat maneuver. You can also specify the timezone you want on you queries

## Installation / Run

On you computer, you'll need to download these softwares :

- Docker
- Go v1.20 (only if you want to run the code without Docker)

And n run these commands

```bash
git clone https://github.com/vareversat/chabo-api.git
cd chabo-api
docker compose build && docker compose run
```

Ta-dam ! The Swagger is running on <http://localhost:8080/v1/swagger/index.html>
