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
- Képek tárolása táblában ---> a zenék hivatkoznak rá, így nem kell mindig új fájlt létrehozni, ha netán album borítónak felel meg
```sql
CREATE TABLE image (
    id INTEGER PRIMARY KEY,
    path TEXT UNIQUE NOT NULL,
    referedCount INTEGER DEFAULT 1,
);
```
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
`Jelenleg nincs fejlesztési terv a backend-hez`

## Frontend továbbfejlesztés:
- TODO-k megcsinálása
- styled(...) komponensek elhagyása -> 'sx' prop használata helyette
- formok konroláltá alakítása
- Zene hozzáadásakor a közreműködő kiválasztó mező képes legyen szerzőt létrehozni