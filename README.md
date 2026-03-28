# littleLabor

A household chores management application for families, built with Go.

## Features (Planned)
- Task/chore management
- Mobile app for iOS (for kids to view/complete chores)
- Web dashboard for family management
- Local network deployment on Raspberry Pi

## Project Structure
```
littleLabor/
├── cmd/                # Command-line applications
│   └── server/         # Main API server
├── internal/           # Internal packages
│   ├── models/         # Data models
│   ├── handlers/       # HTTP handlers
│   └── db/             # Database logic
├── web/                # Frontend web application
├── mobile/             # Mobile app (iOS/Android)
└── docs/               # Documentation
```

## Getting Started

### Prerequisites
- Go 1.21+
- (Frontend/Mobile - TBD)

### Build
```bash
go build -o bin/littleLabor ./cmd/server
```

### Run Locally
```bash
./bin/littleLabor
```

The API will be available at `http://localhost:8080`

## Development

See [DEVELOPMENT.md](./docs/DEVELOPMENT.md) for detailed development setup.

## License

MIT
