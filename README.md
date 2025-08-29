# **Yo**(_u_)**Tu**(_be_) **Do**(_wnload_)

## Futtatás

A futtatásához még kötelezően jelen kell lennie a `wails` futtatókörnyezetnek, de ez a későbbiekben változhat.

A `make.sh` fájl magába foglalja a development szerver futtatását, production build elkészítését és a tesztek futtatását (többnyire linuxon történő futtatáshoz készült).

## Configuráció

### *data/config.yaml*
```yaml
app:
    # Annak a mappának az elérési útvonala, ahova a zenéket metaadatokkal kitöltve átmásolja/letölti a program
    downloadLocation: "/home/***/Music"
    # A youtube letöltő elérési utvonala
    ytdlLocation: "***"
    # Az ffmpeg elérési útvonalja
    ffmpegLocation: "***"
database:
    # A lokális adatbázis elérési útvonalja
    location: "./data/***"
```

## Backend továbbfejlesztés:
- TODO-k megcsinálása
- Logok lemezre történő kiírása
- Log szintek 'config.yaml'-ben történő megadása
- 'Logger' nevű interface, ami megvalósítja a 'Info', 'InfoF', stb... függvényeket. Ezeket tartalmazó tömböt inicializáljon a 'InitializeLogger' függvény, amiken végigmenve pedig meghívódnak a megfelelő függvények.
- Eltárolni az app adatokat bezáráskor, megnyitáskor vissza betölteni (pl.: ablak méret, pozíció, stb...)

## Frontend továbbfejlesztés:
- TODO-k megcsinálása
- Ha a letöltés mappába mozgatjuk át a zenét, a státuszt állítsa át ("Folyamatban" típusra)
