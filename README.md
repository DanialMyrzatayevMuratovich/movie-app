# Analytics Platform 📊

A comprehensive analytics platform for monitoring bookings and financial indicators. The project includes a powerful Go backend for aggregating data from MongoDB and a modern Vite frontend.

# 🛠 Technology stack
- Backend: Go (Golang), MongoDB (Aggregation Framework), Gin (presumably).

- Frontend: JavaScript, Vite, CSS-in-JS/PostCSS.

- Database: MongoDB.

# 🚀 Key features

- Data aggregation: Automatic calculation of revenue, number of tickets, and reservations.

- Financial statistics: Calculation of average check ($AVG$), minimum and maximum order value.

- Rounding: All financial data is automatically rounded to 2 decimal places at the 

- DB.API level: Standardized JSON responses with error handling.

# 💻 Startup Instructions
1. Prerequisites
- Go installed (version 1.20+)

- Node.js installed (version 18+)

- Access to MongoDB database

2. Backend Configuration
- Go to the directory with the server part:
```
cd backend
```
- Create an .env file (if needed) and specify the connection string to MongoDB:
```
MONGO_URI=mongodb://localhost:27017
PORT=8080
```
- Download dependencies and start the server:
```
go mod tidy
go run cmd/main.go
```
3. Frontend Configuration
- Navigate to the frontend directory:
```
cd frontend
```
- Set dependencies:
```
npm install
```
- Start the project in development mode:
```
npm run dev
```

# 📂 Project structure
```
.
├── backend
│   ├── cmd
│   │   └── main.go
│   ├── config
│   │   ├── config.go
│   │   └── database.go
│   ├── go.mod
│   ├── go.sum
│   ├── handlers
│   │   ├── analytics.go
│   │   ├── auth.go
│   │   ├── bookings.go
│   │   ├── cinemas.go
│   │   ├── movies.go
│   │   └── showtimes.go
│   ├── middleware
│   │   ├── auth.go
│   │   ├── error.go
│   │   └── role.go
│   ├── models
│   │   ├── booking.go
│   │   ├── cinema.go
│   │   ├── hall.go
│   │   ├── movie.go
│   │   ├── showtime.go
│   │   ├── transaction.go
│   │   └── user.go
│   ├── routes
│   │   └── routes.go
│   ├── scripts
│   │   ├── create_indexes.go
│   │   └── seed.go
│   └── utils
│       ├── jwt.go
│       ├── response.go
│       └── validation.go
└── frontend
    ├── README.md
    ├── index.html
    ├── package-lock.json
    ├── package.json
    ├── public
    │   └── vite.svg
    ├── src
    │   ├── App.vue
    │   ├── assets
    │   │   └── vue.svg
    │   ├── components
    │   │   ├── BookingCard.vue
    │   │   ├── MovieCard.vue
    │   │   ├── Navbar.vue
    │   │   └── SeatMap.vue
    │   ├── main.js
    │   ├── router
    │   │   └── index.js
    │   ├── services
    │   │   └── api.js
    │   ├── store
    │   │   └── auth.js
    │   ├── style.css
    │   ├── utils
    │   │   └── formatters.js
    │   └── views
    │       ├── Home.vue
    │       ├── Login.vue
    │       ├── MovieDetails.vue
    │       ├── Profile.vue
    │       └── SeatSelection.vue
    └── vite.config.js
```    

# 📄 License
The project is distributed under the MIT license.

```
graph TD
    %% USER BROWSER Section
    subgraph Browser [USER BROWSER]
        direction TB
        subgraph UI_Components [ ]
            direction LR
            Vue[Vue.js 3<br/>Components]
            Router[Vue Router<br/>(Routing)]
            Pinia[Pinia<br/>(State Mgmt)]
        end
        
        Events{UI Events  clicks, forms}
        
        Axios[Axios HTTP Client<br/>- Adds JWT Token to Headers<br/>- Makes RESTful API Calls]
        
        UI_Components --> Events
        Events --> Axios
    end

    %% Flow to Backend
    Axios -- "HTTP Request (JSON + JWT)<br/>POST /api/bookings" --> Backend

    %% GO BACKEND SERVER Section
    subgraph Backend [GO BACKEND SERVER  Port 8080]
        direction TB
        Gin[Gin HTTP Router]
        
        Middleware[JWT Authentication Middleware<br/>- Validates Token<br/>- Extracts User ID<br/>- Returns 401 if Invalid]
        
        Handlers[Request Handlers<br/>- Validate Input<br/>- Business Logic<br/>- Call Database Operations]
        
        Gin --> Middleware
        Middleware --> Handlers
    end

    %% Flow to Database
    Handlers -- "Database Query (BSON)" --> DB

    %% MONGODB DATABASE Section
    subgraph DB [MONGODB DATABASE  Port 27017]
        direction TB
        subgraph Collections [ ]
            direction LR
            Users[Users]
            Movies[Movies]
            Cinemas[Cinemas]
            Halls[Halls]
        end
        
        subgraph Collections2 [ ]
            direction LR
            Showtimes[Showtimes]
            Bookings[Bookings]
            Files[Files]
        end
        
        Ops[OPERATIONS:<br/>• Find/FindOne (with indexes)<br/>• InsertOne/InsertMany<br/>• UpdateOne/UpdateMany<br/>• DeleteOne/DeleteMany<br/>• Aggregate (for analytics)<br/>• Transactions (ACID compliance)]
        
        Collections --- Collections2
        Collections2 --- Ops
    end

    %% Response Flow
    DB -- "Query Results (BSON Documents)" --> Response
    
    Response[Convert to JSON<br/><br/>HTTP Response<br/>Status: 200 OK<br/>Body: { 'success': true, 'data': {...} }]
    
    Response --> Final[USER BROWSER<br/>Updates UI with Response Data]

    %% Styling for visual accuracy
    style Browser fill:none,stroke:#333,stroke-width:1px
    style Backend fill:none,stroke:#333,stroke-width:1px
    style DB fill:none,stroke:#333,stroke-width:1px
    style UI_Components fill:none,stroke:none
    style Collections fill:none,stroke:none
    style Collections2 fill:none,stroke:none
    style Response fill:none,stroke:none
    style Final fill:none,stroke:none
```
