<h1 align="center">iCalHub</h1>

iCalHub is a web application that allows users to find and subscribe to public calendars.

The application is built using the gin framework for Go.

## 🗓 Calendars

<details open>
<summary>🏖️ Holidays</summary>

- China
```shell
/holidays/china
```
</details>

<details open>
<summary>🍿 Movies</summary>


- Upcoming releases - IMDb
```shell
/movies/imdb/:region
```

The region is optional. If not provided, the region will be set to China.

The region code can be found [🔗here](https://www.imdb.com/calendar/).

- Upcoming Movies - Douban
```shell
/movies/douban
````
</details>

<details open>
<summary>📡 Astronomy</summary>

- Date and Time of the Moon Phase｜Hong Kong Observatory(HKO)
```
/astronomy/moon/:year
```

The year is optional. If not provided, the current year will be used.

</details>

<details open>
<summary>🎮 Games</summary>

- Upcoming Releases - Steam
```
/games/steam/:type/:language
```

the `type` must be one of the following: `all`, `popular`;

the `language` is optional, and its default value is `zh_CH`, currently only supports `zh_CH` and `en_US`.


</details>