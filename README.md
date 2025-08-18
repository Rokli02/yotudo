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

## Frontend továbbfejlesztés:
- TODO-k megcsinálása
- Ha a letöltés mappába mozgatjuk át a zenét, a státuszt állítsa át ("Folyamatban" típusra)
