# Repairs Log

In business, repairs are hard to keep up with. Repairs Log can help. This application has 3 core concepts. Manufacturers, Equipment, and Tickets. Heres how it works. First, you go into your business and register the manufacturers and their associated equipment. From there, your team can start creating tickets. When a ticket is created, it becomes accessible in the admin section of the application. Tickets are private at first, but once the public details of a ticket are filled out, their status becomes viewable by the team. This enables you to communicate the status of a repair to your entire team with ease. Once a ticket is complete, you mark its cost and any repairs notes. From there, the ticket will live under its associated piece of equipment so you can view all the repairs that have occured with a speific piece of equipment overtime.

## Installation

Clone the repo

```bash
git clone https://github.com/phillip-england/skimp-shrimp
```

Then install appropriate Go packages from within the repo

```bash
go mod tidy
```

Install Tailwind (you can opt in for using the binary, I used npm/bun) [Tailwind Installtion](https://tailwindcss.com/docs/installation)

Install using Bun
```bash
bun install
```

Install using npm
```bash
npm install
```

Running Tailwind during development (script found in package.json)
```bash
npm run tailwind
```
## Serving

During developement, I used air to hotreload my builds. For more information on air, check here: [Air](https://github.com/cosmtrek/air) It is as easy as installing a binary. To change air config, check .air.toml

Serving with Air
```bash
air
```

Serving without Air
```bash
go run main.go
```

## Env Variables

This application is intended for a single user. My team primarly speaks spanish, so I use a translation API key for converting responses into english for my R&M director. Here are the following env variables I used:

```env
ADMIN_USERNAME= <some-username>
ADMIN_PASSWORD= <some-password>
ADMIN_SESSION_TOKEN= <some-token>
PORT=8080
BASE_URL= <localhost:8080 for dev || your baseURL during prod>
PUBLIC_SECURITY_TOKEN= <some-token>
GO_ENV=dev
TRANSLATION_API_KEY= <api key https://api-free.deepl.com for translation> OPTIONAL
TRANSLATION_API_URL= <url for translation > OPTIONAL
TRANSLATION_API_RESOURCE= <resource for translation> OPTIONAL
```

## Public Views

This application has a few secure public views. They require a queryparam which includes your public security token from your .env file

This enables you to post a secure QR code in your business for ease of access for team members.

```go
 /app/ticket/public?publicSecurityToken=<your public security token>
 /app/ticket/public/view?publicSecurityToken=<your public security token>
```

## Database

This application uses SQlite during development. Planning to swap it over to postgres for production. Have not built this out yet but this is on my radar.