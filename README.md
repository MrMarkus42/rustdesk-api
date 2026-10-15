**Deutsch** · [English](README.en.md)

# RustDesk API

In Go geschriebene API für selbst gehostete [RustDesk](https://rustdesk.com)-Server, mit Admin-Oberfläche. Sie bildet die Schnittstellen nach, die der RustDesk-Client für Anmeldung, Adressbuch und Geräteverwaltung nutzt.

Dieser Fork baut auf [lejianwen/rustdesk-api](https://github.com/lejianwen/rustdesk-api) (v2.7) auf, der nicht mehr gepflegt wird. Er wird zusammen mit [MrMarkus42/rustdesk-server](https://github.com/MrMarkus42/rustdesk-server) als s6-Image ausgeliefert.

## Funktionen

- **Client-API:** Anmeldung (Passwort, OIDC, GitHub, Google, LDAP), Adressbücher mit Tags und Freigaben, Gruppen, Online-Status, Verbindungs- und Dateiprotokolle
- **Web-Client** (eigene Entwicklung): Geräte direkt in der Admin-Oberfläche im Browser fernsteuern, Ende-zu-Ende-verschlüsselt, ohne zusätzliche Ports
- **Admin-Oberfläche** (`/_admin/`) auf **Deutsch** und Englisch (weitere: Französisch, Spanisch, Russisch, Koreanisch, Chinesisch): Benutzer, Geräte, Adressbücher, Gruppen, OAuth, Anmeldeprotokoll, Server-Befehle

### Änderungen gegenüber lejianwen

- **Sicherheit:** Captcha und IP-Sperre lassen sich nicht mehr per gefälschtem `X-Forwarded-For` umgehen, zufällige Access-Tokens, keine Verzeichnisauflistung, aktualisierte Abhängigkeiten (govulncheck und npm audit ohne Befund)
- **Login:** Clients werden nicht mehr nach 7 Tagen abgemeldet, wenn `token-expire` höher gesetzt ist
- **Admin-Oberfläche:** deutsche Übersetzung, englischer Standard statt Chinesisch, Formulare können Felder wieder leeren, OIDC mit Amazon Cognito
- **Lieferkette:** die Admin-Oberfläche wird aus einem geprüften, fest angegebenen Stand von [MrMarkus42/rustdesk-api-web](https://github.com/MrMarkus42/rustdesk-api-web) gebaut

## Docker-Image

```
ghcr.io/mrmarkus42/rustdesk-api:<version>
```

Für den Betrieb zusammen mit hbbs/hbbr empfiehlt sich das kombinierte Image `ghcr.io/mrmarkus42/rustdesk-server-s6` (siehe [rustdesk-server](https://github.com/MrMarkus42/rustdesk-server)).

## Wichtige Einstellungen

Alle Einstellungen aus `conf/config.yaml` lassen sich per Umgebungsvariable setzen: Präfix `RUSTDESK_API_`, Abschnitt und Schlüssel mit `_` verbunden, Bindestriche als `_` (Beispiel: `gin.trust-proxy` → `RUSTDESK_API_GIN_TRUST_PROXY`).

| Variable | Bedeutung |
|---|---|
| `RUSTDESK_API_RUSTDESK_ID_SERVER` | ID-Server, den die API an Clients weitergibt |
| `RUSTDESK_API_RUSTDESK_RELAY_SERVER` | Relay-Server |
| `RUSTDESK_API_RUSTDESK_API_SERVER` | Öffentliche Adresse der API (z. B. `https://rustdesk.example.org`) |
| `RUSTDESK_API_RUSTDESK_KEY_FILE` | Öffentlicher Server-Schlüssel (Standard `/data/id_ed25519.pub`) |
| `RUSTDESK_API_JWT_KEY` | Signatur der Login-Tokens, nötig für `MUST_LOGIN`. **Nicht nachträglich ändern** |
| `RUSTDESK_API_APP_TOKEN_EXPIRE` | Gültigkeit eines Logins (Standard `168h`) |
| `RUSTDESK_API_GIN_TRUST_PROXY` | Adresse/Bereich des Reverse Proxys (kommagetrennt, CIDR möglich) |
| `RUSTDESK_API_APP_BAN_THRESHOLD` | IP-Sperre nach n Fehlversuchen (`0` = aus). Hinter einem Proxy nur mit `TRUST_PROXY` |
| `RUSTDESK_API_APP_REGISTER` | Selbstregistrierung (`false` empfohlen) |
| `RUSTDESK_API_APP_WEB_CLIENT` | Web-Client in der Admin-Oberfläche (`1` = an, Standard; `0` = aus) |
| `RUSTDESK_API_RUSTDESK_WS_ID_TARGET` / `…_WS_RELAY_TARGET` | WebSocket-Ports von hbbs/hbbr, zu denen die API den Web-Client weiterleitet (Standard `ws://127.0.0.1:21118` / `ws://127.0.0.1:21119`, passend für das s6-Image) |
| `RUSTDESK_API_LANG` | Sprache der API-Meldungen ohne Browser-Angabe (Standard `en`) |

### Web-Client

Der Browser verbindet sich nur mit der API (`/webclient/ws/id` und `/webclient/ws/relay`); die API leitet intern an hbbs/hbbr weiter. Jede WebSocket-Verbindung braucht ein Einmal-Ticket (30 Sekunden gültig), das nur angemeldete Benutzer erhalten. Hinter einem Reverse Proxy muss dieser WebSocket-Verbindungen durchlassen (beim Nginx Proxy Manager: „Websockets Support“). Der Browser prüft die Identität des Geräts mit dem Server-Schlüssel und verschlüsselt die Sitzung Ende-zu-Ende; ohne Server-Schlüssel wird keine Verbindung aufgebaut.

Admin-Passwort zurücksetzen: `./apimain reset-admin-pwd <neues-passwort>`

## Selbst bauen

```bash
go build -o apimain ./cmd/apimain.go
```

Die Admin-Oberfläche wird separat aus [rustdesk-api-web](https://github.com/MrMarkus42/rustdesk-api-web) gebaut (`npm ci && npm run build`) und nach `resources/admin/` kopiert. Die Workflows unter `.github/workflows/` erledigen das automatisch.

## Herkunft und Lizenz

- Ursprüngliches Projekt: [lejianwen/rustdesk-api](https://github.com/lejianwen/rustdesk-api)
- Lizenz: [MIT](LICENSE)
