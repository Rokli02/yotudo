# **Yo**(_u_)**Tu**(_be_) **Do**(_wnload_)

## Futtatás

A futtatásához még kötelezően jelen kell lennie a `wails` futtatókörnyezetnek, de ez a későbbiekben változhat.

A `make.sh` fájl magába foglalja a development szerver futtatását, production build elkészítését és a tesztek futtatását (többnyire linuxon történő futtatáshoz készült).

## Configuráció

Szükséges a megléte két másik programnak a működéshez:
- <u>*__ffmpeg__*</u>: A zene fájlok feldolgozásához szükséges (https://ffmpeg.org/download.html)
- <u>*__yt-dlp__*</u>: A zenék letöltéséhez szükséges (https://github.com/yt-dlp/yt-dlp?tab=readme-ov-file#installation)

### *data/config.yaml*
```yaml
app:
    # Annak a mappának az elérési útvonala, ahova a zenéket metaadatokkal kitöltve átmásolja/letölti a program
    downloadLocation: "/home/***/Music"
    # A youtube letöltő elérési utvonala
    ytdlLocation: "***"
    # Az ffmpeg elérési útvonalja
    ffmpegLocation: "***"
logger:
    # Milyen szintű futási információkat osszon meg velünk a program
    level: info
    # Ezeket az információkat milyen formában és hol közölje velünk - a konzol alapú egyedül fejlesztés közben indokolt
    types:
        - console
        - file
```

## Backend továbbfejlesztés:
- Album tábla autocomplete-hez, egy-egy album egy adott előadóhoz kötődjön
```sql
CREATE TABLE album (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    UNIQUE(author_id, name)
);
```
- Abort logika implementálása a felesleges lekérdezések elkerüléséhez
- TODO-k megcsinálása

## Frontend továbbfejlesztés:
- TODO-k megcsinálása
- formok konroláltá alakítása
- Szerver hosztolás, amin lehet látni, belehallgatni a zenékbe és letölteni azokat