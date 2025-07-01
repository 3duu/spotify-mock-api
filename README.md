Spotify Mock API – Go Backend
A fully featured, lightweight mock backend that emulates core Spotify API endpoints for mobile and web apps.
Built in Go (Golang), using Gin and SQLite.
Ideal for local/mobile development, hackathons, and building Spotify clones without dealing with real API keys or licensing.

Features
🎵 REST endpoints for tracks, albums, playlists, artists, podcasts, user profiles, newsletters, recent plays, recommendations, and more

💾 Data is loaded from JSON on startup and persisted in local SQLite (app.db)

▶️ Serves real MP3 files and cover art images from /media folder

🔑 Simulated OAuth/auth endpoints: /login, /me, etc.

📨 Newsletter/"TGIF" cards for home feed

📋 Recent played items, recommendations, etc.

🧩 Compatible with mobile and web frontends (React Native/Expo, React.js, etc.)

Requirements
Go (1.18+)

Git

(For development: sqlite3 CLI for inspecting data)

Quickstart
1. Clone the repository
bash
Copiar
Editar
git clone https://github.com/3duu/spotify-mock-api.git
cd spotify-mock-api
2. Download Go dependencies
bash
Copiar
Editar
go mod download
3. Prepare your media and data
Place all MP3 files and cover art images in data/media/

Place your seed data (tracks, albums, playlists, users, newsletters, etc) in data/defaults.json

4. Run the backend server
bash
Copiar
Editar
go run main.go
The server will start at http://localhost:8080

All API endpoints are available under / (see Endpoints below)

Media files are served from /media (e.g. http://localhost:8080/media/track1.mp3)

Endpoints
A selection of implemented endpoints:

Method	Endpoint	Description
GET	/tracks/:id	Get metadata and audio URL for a track
GET	/albums/:id	Get album details and track list
GET	/artists/:id	Get artist info and their top tracks
GET	/playlists	List all playlists
GET	/playlists/:id	Playlist details and tracks
GET	/users/:id/recent-playlists	Recent playlists for a user
GET	/search	Search tracks, artists, albums, playlists
GET	/newsletters	Newsletter/TGIF home cards
GET	/podcasts/:id	Podcast details and episodes
GET	/library	User’s saved playlists/albums/podcasts
POST	/login	Simulate user login, returns mock token
GET	/me	Get profile of the current user
...	...	...

For the full list, see the source code or generate docs from comments.

Data Model & Structure
SQLite (app.db) stores all user data, tracks, albums, playlists, etc.

data/defaults.json is loaded at startup for initial seed data (can be edited to add more songs/playlists/newsletters)

data/media/ contains all referenced audio and image files

Editing & Adding Data
Add tracks/albums/etc: Edit data/defaults.json and re-run go run main.go

Add media: Place new MP3s or cover art in data/media/ (reference the filenames in your JSON)

Wipe/reseed: Delete app.db and restart server for a clean seed

For Development With Mobile App
Ensure backend is running (go run main.go)

Set the API base URL in your mobile app (React Native) to match your machine’s local network IP (e.g., http://192.168.0.10:8080)

The backend is CORS-enabled and safe for use with mobile, emulator, or web

Project Structure
bash
Copiar
Editar
spotify-mock-api/
│
├── cmd/                # CLI entrypoint (main.go)
├── internal/
│   ├── handlers/       # API endpoint handlers
│   ├── models/         # DB models
│   ├── data/           # Seed and static data
│   └── ...
├── data/
│   ├── defaults.json   # Seed data (songs, albums, users, newsletters, etc.)
│   └── media/          # MP3s and images
├── app.db              # SQLite database (auto-created)
└── README.md
Troubleshooting
"file not found": Make sure referenced images/MP3s are present in data/media/

Cannot connect from mobile device: Double check API URL matches your computer’s local IP and firewall allows connections

Seeding errors: Delete app.db and restart, verify your JSON formatting

Contributing
Pull requests, issues, and suggestions are welcome!

License
MIT © 3duu
