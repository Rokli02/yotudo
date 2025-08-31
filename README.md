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
logger:
    level: info
    types:
        - console
        - file
```

## Backend továbbfejlesztés:
- Thumbnail automatikus letöltése előre meghatározott URL alapján, nem pedig a letöltő általi beszerzés -> így lehetőség van már az első mentéskor kép megjelenítésére

## Frontend továbbfejlesztés:
- TODO-k megcsinálása
- styled(...) komponensek elhagyása -> 'sx' prop használata helyette
- formok konroláltá alakítása
- Ha a letöltés mappába mozgatjuk át a zenét, a státuszt állítsa át ("Folyamatban" típusra)
- picFileSelector componens:
    - Ha VAN kiválasztott kép, akkor mutasson egy default kép ikon-t és a kép nevét, valamint egy 'X'-et, amivel ki tudjuk szedni
    - Ha NINCS, akkor legyen egy switch komponens ami tartalmazza a következő elemeket és az alapján változzon a komponens működése:
        - none: Ennyi, nincs semmi
        - thumbnail: Továbbra sincs semmi, gyakorlatilag úgy veszi, hogy az eddigi checkbox ki lett pipálva
        - web: Input mező, ami egy URL-t vár egy képről
        - local: Gomb, ami egy file selector-t nyit meg 
