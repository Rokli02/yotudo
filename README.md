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
- Album tábla autocomplete-hez, egy-egy album egy adott
```sql
CREATE TABLE album (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    UNIQUE(author_id, name)
);

-- For migration-hydration
INSERT OR IGNORE INTO album(...) SELECT album as name, author_id FROM music;
```

## Frontend továbbfejlesztés:
- TODO-k megcsinálása
- formok konroláltá alakítása