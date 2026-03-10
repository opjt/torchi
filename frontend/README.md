# Torchi PWA

Frontend application for **Torchi** — a push notification platform via HTTP API. Subscribe and manage push notifications directly from the web without installing a separate app.

> torchi(Backend): [opjt/torchi](https://github.com/opjt/torchi)  
> Live service -> [https://torchi.app](https://torchi.app)

## Tech Stack

| Category     | Technology                  |
| ------------ | --------------------------- |
| Framework    | SvelteKit (Svelte 5)        |
| Build        | Vite                        |
| Styling      | Tailwind CSS 4, DaisyUI     |
| UI Component | shadcn-svelte, Lucide Icons |
| PWA          | vite-plugin-pwa, Workbox    |
| Deploy       | Static Adapter (SPA)        |

## Project Structure

```bash
src/
├── lib/
│   ├── api/            # API call modules
│   ├── assets/         # Static resources
│   ├── client/         # Client-side logic
│   │   └── auth/       # Authentication
│   ├── components/     # Shared components
│   │   ├── lib/        # Project components
│   │   └── ui/         # UI components (button, card, dialog, etc.)
│   └── pkg/            # Internal packages (utilities)
├── routes/
│   ├── app/
│   │   ├── guide/      # Usage guide
│   │   ├── services/   # Endpoint management
│   │   ├── setting/    # Settings
│   │   └── welcome/    # Onboarding
│   ├── privacy/        # Privacy policy
│   └── terms/          # Terms of service
├── app.html            # HTML template
├── app.css             # Global styles
└── service-worker.ts   # Push notification service worker
```

## Getting Started

```bash
# Install dependencies
pnpm install

# Start dev server
pnpm dev

# Production build
pnpm build

# Preview build
pnpm preview
```

## Environment Variables

Copy `.env.sample` to `.env` and fill in the values:

```bash
cp .env.sample .env
```

| Variable                  | Description                       |
| ------------------------- | --------------------------------- |
| `PUBLIC_VAPID_KEY`        | VAPID public key for Web Push API |
| `PUBLIC_GITHUB_CLIENT_ID` | GitHub OAuth app client ID        |

> VAPID key pair can be generated via the backend server (`go run . genkey`).
> GitHub Client ID is obtained by creating a GitHub OAuth App.

## License

MIT License
