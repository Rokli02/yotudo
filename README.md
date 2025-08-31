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
`Jelenleg nincs fejlesztési terv a backend-hez`

## Frontend továbbfejlesztés:
- TODO-k megcsinálása
- styled(...) komponensek elhagyása -> 'sx' prop használata helyette
- formok konroláltá alakítása
- Zene hozzáadásakor a szerző kiválasztó mező képes legyen szerzőt hozzáadni is