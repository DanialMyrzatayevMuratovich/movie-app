package scripts

import (
	"cinema-booking/config"
	"cinema-booking/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// Глобальные переменные для хранения ID созданных объектов
var (
	userIDs     []primitive.ObjectID
	cinemaIDs   []primitive.ObjectID
	hallIDs     []primitive.ObjectID
	movieIDs    []primitive.ObjectID
	showtimeIDs []primitive.ObjectID
)

func SeedDatabase() {
	ctx := context.Background()

	log.Println("🌱 Starting database seeding...")

	// Очистить существующие данные (опционально)
	clearCollections(ctx)

	// 1. Создать пользователей
	seedUsers(ctx)

	// 2. Создать кинотеатры
	seedCinemas(ctx)

	// 3. Создать залы
	seedHalls(ctx)

	// 4. Создать фильмы
	seedMovies(ctx)

	// 5. Создать сеансы
	seedShowtimes(ctx)

	// 6. Создать тестовые брони
	seedBookings(ctx)

	// 7. Создать тестовые транзакции
	seedTransactions(ctx)

	log.Println("✅ Database seeding completed successfully!")
}

// Очистка коллекций
func clearCollections(ctx context.Context) {
	collections := []string{"users", "cinemas", "halls", "movies", "showtimes", "bookings", "transactions"}

	for _, collName := range collections {
		coll := config.GetCollection(collName)
		if err := coll.Drop(ctx); err != nil {
			log.Printf("⚠️ Warning: Could not drop collection %s: %v", collName, err)
		}
	}

	log.Println("🗑️  Collections cleared")
}

// 1. USERS
func seedUsers(ctx context.Context) {
	collection := config.GetCollection("users")

	// Hash пароль функция
	hashPassword := func(password string) string {
		hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		return string(hash)
	}

	users := []models.User{
		{
			Email:    "admin@cinema.kz",
			Password: hashPassword("admin123"),
			FullName: "Admin User",
			Phone:    "+77001234567",
			Role:     "admin",
			Wallet: models.Wallet{
				Balance:  0,
				Currency: "KZT",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Email:    "manager@cinema.kz",
			Password: hashPassword("manager123"),
			FullName: "Cinema Manager",
			Phone:    "+77001234568",
			Role:     "cinema_manager",
			Wallet: models.Wallet{
				Balance:  0,
				Currency: "KZT",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Email:    "user1@example.com",
			Password: hashPassword("user123"),
			FullName: "Айдар Қасымов",
			Phone:    "+77771234567",
			Role:     "user",
			Wallet: models.Wallet{
				Balance:  5000,
				Currency: "KZT",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Email:    "user2@example.com",
			Password: hashPassword("user123"),
			FullName: "Асель Нұрғалиева",
			Phone:    "+77771234568",
			Role:     "user",
			Wallet: models.Wallet{
				Balance:  10000,
				Currency: "KZT",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Email:    "user3@example.com",
			Password: hashPassword("user123"),
			FullName: "Ерлан Сыдықов",
			Phone:    "+77771234569",
			Role:     "user",
			Wallet: models.Wallet{
				Balance:  3000,
				Currency: "KZT",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for i := range users {
		result, err := collection.InsertOne(ctx, users[i])
		if err != nil {
			log.Printf("❌ Error inserting user: %v", err)
			continue
		}
		userIDs = append(userIDs, result.InsertedID.(primitive.ObjectID))
	}

	log.Printf("✅ Created %d users", len(userIDs))
}

// 2. CINEMAS
func seedCinemas(ctx context.Context) {
	collection := config.GetCollection("cinemas")

	cinemas := []models.Cinema{
		{
			Name:    "Chaplin Cinemas Mega Alma-Ata",
			City:    "Almaty",
			Address: "ул. Розыбакиева, 247А, ТРЦ Mega",
			Location: models.Location{
				Type:        "Point",
				Coordinates: []float64{76.8512, 43.2061}, // [longitude, latitude]
			},
			Facilities:   []string{"IMAX", "4DX", "VIP", "Parking", "Food Court"},
			HallIDs:      []primitive.ObjectID{}, // заполним позже
			Rating:       4.6,
			TotalReviews: 1523,
			Images:       []string{"/uploads/cinemas/chaplin_mega.jpg"},
			CreatedAt:    time.Now(),
		},
		{
			Name:    "Kinopark Sary-Arka",
			City:    "Almaty",
			Address: "пр. Сары-Арка, 10, ТРЦ Sary-Arka",
			Location: models.Location{
				Type:        "Point",
				Coordinates: []float64{76.9286, 43.2425},
			},
			Facilities:   []string{"3D", "VIP", "Parking"},
			HallIDs:      []primitive.ObjectID{},
			Rating:       4.3,
			TotalReviews: 892,
			Images:       []string{"/uploads/cinemas/kinopark_saryarka.jpg"},
			CreatedAt:    time.Now(),
		},
		{
			Name:    "Arman Cinema Dostyk Plaza",
			City:    "Almaty",
			Address: "пр. Достык, 111, Достык Плаза",
			Location: models.Location{
				Type:        "Point",
				Coordinates: []float64{76.9539, 43.2324},
			},
			Facilities:   []string{"3D", "IMAX", "VIP", "Dolby Atmos"},
			HallIDs:      []primitive.ObjectID{},
			Rating:       4.7,
			TotalReviews: 2105,
			Images:       []string{"/uploads/cinemas/arman_dostyk.jpg"},
			CreatedAt:    time.Now(),
		},
	}

	for i := range cinemas {
		result, err := collection.InsertOne(ctx, cinemas[i])
		if err != nil {
			log.Printf("❌ Error inserting cinema: %v", err)
			continue
		}
		cinemaIDs = append(cinemaIDs, result.InsertedID.(primitive.ObjectID))
	}

	log.Printf("✅ Created %d cinemas", len(cinemaIDs))
}

// 3. HALLS
func seedHalls(ctx context.Context) {
	collection := config.GetCollection("halls")

	// Генератор мест для зала
	generateSeats := func(rows int, seatsPerRow int, hallType string) []models.Seat {
		seats := []models.Seat{}
		rowLetters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

		basePrice := 2000.0
		switch hallType {
		case "VIP":
			basePrice = 4000.0
		case "IMAX":
			basePrice = 3500.0
		}

		for r := 0; r < rows; r++ {
			for s := 1; s <= seatsPerRow; s++ {
				seatType := "regular"
				price := basePrice

				// VIP места в середине
				if hallType == "VIP" || (r >= rows/3 && r <= 2*rows/3 && s >= seatsPerRow/4 && s <= 3*seatsPerRow/4) {
					seatType = "vip"
					price = basePrice * 1.5
				}

				// Couple seats в последних рядах
				if r >= rows-2 && s%2 == 1 && s < seatsPerRow {
					seatType = "couple"
					price = basePrice * 1.3
				}

				seats = append(seats, models.Seat{
					Row:    string(rowLetters[r]),
					Number: s,
					Type:   seatType,
					Price:  price,
				})
			}
		}
		return seats
	}

	// Залы для каждого кинотеатра
	halls := []models.Hall{}

	// Chaplin Mega (3 зала)
	if len(cinemaIDs) > 0 {
		halls = append(halls,
			models.Hall{
				CinemaID:   cinemaIDs[0],
				Name:       "IMAX Hall 1",
				HallNumber: 1,
				Capacity:   240,
				Type:       "IMAX",
				Seats:      generateSeats(12, 20, "IMAX"),
			},
			models.Hall{
				CinemaID:   cinemaIDs[0],
				Name:       "Standard Hall 2",
				HallNumber: 2,
				Capacity:   150,
				Type:       "Standard",
				Seats:      generateSeats(10, 15, "Standard"),
			},
			models.Hall{
				CinemaID:   cinemaIDs[0],
				Name:       "VIP Hall 3",
				HallNumber: 3,
				Capacity:   50,
				Type:       "VIP",
				Seats:      generateSeats(5, 10, "VIP"),
			},
		)
	}

	// Kinopark Sary-Arka (2 зала)
	if len(cinemaIDs) > 1 {
		halls = append(halls,
			models.Hall{
				CinemaID:   cinemaIDs[1],
				Name:       "Standard Hall 1",
				HallNumber: 1,
				Capacity:   180,
				Type:       "Standard",
				Seats:      generateSeats(12, 15, "Standard"),
			},
			models.Hall{
				CinemaID:   cinemaIDs[1],
				Name:       "3D Hall 2",
				HallNumber: 2,
				Capacity:   120,
				Type:       "Standard",
				Seats:      generateSeats(10, 12, "Standard"),
			},
		)
	}

	// Arman Dostyk (3 зала)
	if len(cinemaIDs) > 2 {
		halls = append(halls,
			models.Hall{
				CinemaID:   cinemaIDs[2],
				Name:       "IMAX Hall 1",
				HallNumber: 1,
				Capacity:   280,
				Type:       "IMAX",
				Seats:      generateSeats(14, 20, "IMAX"),
			},
			models.Hall{
				CinemaID:   cinemaIDs[2],
				Name:       "VIP Hall 2",
				HallNumber: 2,
				Capacity:   60,
				Type:       "VIP",
				Seats:      generateSeats(6, 10, "VIP"),
			},
			models.Hall{
				CinemaID:   cinemaIDs[2],
				Name:       "Standard Hall 3",
				HallNumber: 3,
				Capacity:   200,
				Type:       "Standard",
				Seats:      generateSeats(12, 17, "Standard"),
			},
		)
	}

	// Вставить все залы
	for i := range halls {
		result, err := collection.InsertOne(ctx, halls[i])
		if err != nil {
			log.Printf("❌ Error inserting hall: %v", err)
			continue
		}
		hallID := result.InsertedID.(primitive.ObjectID)
		hallIDs = append(hallIDs, hallID)

		// Обновить cinema с hallID
		cinemaCol := config.GetCollection("cinemas")
		cinemaCol.UpdateOne(
			ctx,
			primitive.M{"_id": halls[i].CinemaID},
			primitive.M{"$push": primitive.M{"hallIds": hallID}},
		)
	}

	log.Printf("✅ Created %d halls", len(hallIDs))
}

// 4. MOVIES
func seedMovies(ctx context.Context) {
	collection := config.GetCollection("movies")

	movies := []models.Movie{
		{
			Title:          "Dune: Part Three",
			TitleKz:        "Құм: 3-бөлім",
			TitleRu:        "Дюна: Часть третья",
			Description:    "The epic conclusion to Denis Villeneuve's Dune trilogy. Paul Atreides unites with Chani and the Fremen while seeking revenge against the conspirators who destroyed his family.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/d5NXSklXo0qyIYkgV94XAgMIckC.jpg", // ← ДОБАВЛЕНО
			Director:       "Denis Villeneuve",
			Cast:           []string{"Timothée Chalamet", "Zendaya", "Austin Butler", "Florence Pugh"},
			Genres:         []string{"Sci-Fi", "Adventure", "Drama"},
			Duration:       165,
			ReleaseDate:    time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			Rating:         "PG-13",
			IMDBRating:     8.9,
			Language:       []string{"English", "Russian"},
			Subtitles:      []string{"Kazakh", "Russian", "English"},
			IsActive:       true,
			AgeRestriction: 13,
			Reviews:        []models.Review{},
			CreatedAt:      time.Now(),
		},
		{
			Title:          "The Batman 2",
			TitleKz:        "Бэтмен 2",
			TitleRu:        "Бэтмен 2",
			Description:    "Batman continues his war on crime in Gotham City while facing new villains and uncovering deeper conspiracies.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/74xTEgt7R36Fpooo50r9T25onhq.jpg", // ← ДОБАВЛЕНО
			Director:       "Matt Reeves",
			Cast:           []string{"Robert Pattinson", "Zoë Kravitz", "Colin Farrell", "Paul Dano"},
			Genres:         []string{"Action", "Crime", "Thriller"},
			Duration:       155,
			ReleaseDate:    time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
			Rating:         "PG-13",
			IMDBRating:     8.5,
			Language:       []string{"English"},
			Subtitles:      []string{"Kazakh", "Russian"},
			IsActive:       true,
			AgeRestriction: 13,
			Reviews:        []models.Review{},
			CreatedAt:      time.Now(),
		},
		{
			Title:          "Avatar: Fire and Ash",
			TitleKz:        "Аватар: От және Күл",
			TitleRu:        "Аватар: Огонь и Пепел",
			Description:    "The third installment in James Cameron's Avatar saga explores new regions of Pandora and introduces the Ash People.",
			PosterURL:      "https://i.ytimg.com/vi/nb_fFj_0rq8/maxresdefault.jpg", // ← ДОБАВЛЕНО
			Director:       "James Cameron",
			Cast:           []string{"Sam Worthington", "Zoe Saldaña", "Sigourney Weaver", "Kate Winslet"},
			Genres:         []string{"Sci-Fi", "Adventure", "Fantasy"},
			Duration:       190,
			ReleaseDate:    time.Date(2025, 12, 20, 0, 0, 0, 0, time.UTC),
			Rating:         "PG-13",
			IMDBRating:     8.7,
			Language:       []string{"English", "Russian"},
			Subtitles:      []string{"Kazakh", "Russian", "English"},
			IsActive:       true,
			AgeRestriction: 13,
			Reviews:        []models.Review{},
			CreatedAt:      time.Now(),
		},
		{
			Title:          "Mission: Impossible 8",
			TitleKz:        "Мүмкін емес миссия 8",
			TitleRu:        "Миссия невыполнима 8",
			Description:    "Ethan Hunt and his team face their most dangerous mission yet in the thrilling conclusion to the franchise.",
			PosterURL:      "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSIQnVwg6itVzcxg0K1xwhB5VaXkyGWLmO07w&s", // ← ДОБАВЛЕНО
			Director:       "Christopher McQuarrie",
			Cast:           []string{"Tom Cruise", "Hayley Atwell", "Ving Rhames", "Simon Pegg"},
			Genres:         []string{"Action", "Thriller", "Adventure"},
			Duration:       140,
			ReleaseDate:    time.Date(2025, 5, 23, 0, 0, 0, 0, time.UTC),
			Rating:         "PG-13",
			IMDBRating:     8.3,
			Language:       []string{"English"},
			Subtitles:      []string{"Kazakh", "Russian"},
			IsActive:       true,
			AgeRestriction: 13,
			Reviews:        []models.Review{},
			CreatedAt:      time.Now(),
		},
		{
			Title:          "Nosferatu",
			TitleKz:        "Носферату",
			TitleRu:        "Носферату",
			Description:    "Robert Eggers' gothic horror reimagining of the classic vampire tale. A haunting story of obsession between a haunted young woman and the terrifying vampire infatuated with her.",
			PosterURL:      "https://upload.wikimedia.org/wikipedia/ru/b/b7/Nosferatu%2C_eine_Symphonie_des_Grauens.jpg", // ← ДОБАВЛЕНО
			Director:       "Robert Eggers",
			Cast:           []string{"Bill Skarsgård", "Lily-Rose Depp", "Nicholas Hoult", "Willem Dafoe"},
			Genres:         []string{"Horror", "Fantasy", "Drama"},
			Duration:       132,
			ReleaseDate:    time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC),
			Rating:         "R",
			IMDBRating:     8.1,
			Language:       []string{"English"},
			Subtitles:      []string{"Kazakh", "Russian"},
			IsActive:       true,
			AgeRestriction: 18,
			Reviews:        []models.Review{},
			CreatedAt:      time.Now(),
		},
		{
			Title:          "Inside Out 2",
			TitleKz:        "Ішкі әлем 2",
			TitleRu:        "Головоломка 2",
			Description:    "Riley enters her teenage years and new emotions appear in Headquarters, creating chaos as they navigate the complexities of growing up.",
			PosterURL:      "https://image.tmdb.org/t/p/w500/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg", // ← ДОБАВЛЕНО
			Director:       "Kelsey Mann",
			Cast:           []string{"Amy Poehler", "Maya Hawke", "Phyllis Smith", "Lewis Black"},
			Genres:         []string{"Animation", "Adventure", "Comedy", "Family"},
			Duration:       96,
			ReleaseDate:    time.Date(2024, 6, 14, 0, 0, 0, 0, time.UTC),
			Rating:         "PG",
			IMDBRating:     8.6,
			Language:       []string{"English", "Russian", "Kazakh"},
			Subtitles:      []string{"Kazakh", "Russian", "English"},
			IsActive:       true,
			AgeRestriction: 6,
			Reviews:        []models.Review{},
			CreatedAt:      time.Now(),
		},
		{
			Title:          "Gladiator II",
			TitleKz:        "Гладиатор II",
			TitleRu:        "Гладиатор II",
			Description:    "Years after witnessing the death of Maximus at the hands of his uncle, Lucius must enter the Colosseum after the powerful emperors of Rome conquer his home.",
			PosterURL:      "https://upload.wikimedia.org/wikipedia/ru/6/68/%D0%93%D0%BB%D0%B0%D0%B4%D0%B8%D0%B0%D1%82%D0%BE%D1%80_2.jpg", // ← ДОБАВЛЕНО
			Director:       "Ridley Scott",
			Cast:           []string{"Paul Mescal", "Denzel Washington", "Pedro Pascal", "Connie Nielsen"},
			Genres:         []string{"Action", "Adventure", "Drama"},
			Duration:       148,
			ReleaseDate:    time.Date(2024, 11, 22, 0, 0, 0, 0, time.UTC),
			Rating:         "R",
			IMDBRating:     7.9,
			Language:       []string{"English"},
			Subtitles:      []string{"Kazakh", "Russian"},
			IsActive:       true,
			AgeRestriction: 16,
			Reviews:        []models.Review{},
			CreatedAt:      time.Now(),
		},
		{
			Title:       "Wicked",
			TitleKz:     "Зұлым",
			TitleRu:     "Злая",
			Description: "The untold story of the witches of Oz. A vivid reimagining of the classic Wizard of Oz tale, exploring the complex friendship between Elphaba and Glinda.",
			PosterURL:   "https://upload.wikimedia.org/wikipedia/ru/0/00/Wicked_%282024%29_poster.jpg", // ← ДОБАВЛЕНО

			Director:       "Jon M. Chu",
			Cast:           []string{"Cynthia Erivo", "Ariana Grande", "Michelle Yeoh", "Jeff Goldblum"},
			Genres:         []string{"Musical", "Fantasy", "Romance"},
			Duration:       160,
			ReleaseDate:    time.Date(2024, 11, 27, 0, 0, 0, 0, time.UTC),
			Rating:         "PG",
			IMDBRating:     8.4,
			Language:       []string{"English"},
			Subtitles:      []string{"Kazakh", "Russian"},
			IsActive:       true,
			AgeRestriction: 12,
			Reviews:        []models.Review{},
			CreatedAt:      time.Now(),
		},
	}

	// Добавить тестовые отзывы к некоторым фильмам
	if len(userIDs) >= 3 {
		movies[0].Reviews = []models.Review{
			{
				UserID:    userIDs[2],
				Rating:    9.5,
				Comment:   "Шынайы ғылыми-фантастикалық шедевр! Көріністер таңқаларлық!",
				CreatedAt: time.Now().Add(-48 * time.Hour),
			},
			{
				UserID:    userIDs[3],
				Rating:    9.0,
				Comment:   "Лучший научно-фантастический фильм последних лет!",
				CreatedAt: time.Now().Add(-24 * time.Hour),
			},
		}

		movies[5].Reviews = []models.Review{
			{
				UserID:    userIDs[4],
				Rating:    10.0,
				Comment:   "Балаларға және ересектерге арналған керемет мультфильм!",
				CreatedAt: time.Now().Add(-72 * time.Hour),
			},
		}
	}

	for i := range movies {
		result, err := collection.InsertOne(ctx, movies[i])
		if err != nil {
			log.Printf("❌ Error inserting movie: %v", err)
			continue
		}
		movieIDs = append(movieIDs, result.InsertedID.(primitive.ObjectID))
	}

	log.Printf("✅ Created %d movies", len(movieIDs))
}

// 5. SHOWTIMES
func seedShowtimes(ctx context.Context) {
	collection := config.GetCollection("showtimes")

	// Создать сеансы на ближайшие 7 дней
	now := time.Now()
	showtimes := []models.Showtime{}

	// Время начала сеансов
	sessionTimes := []int{10, 13, 16, 19, 22} // часы

	for day := 0; day < 7; day++ {
		currentDate := now.AddDate(0, 0, day)

		// Для каждого фильма
		for movieIdx, movieID := range movieIDs {
			// В разных кинотеатрах
			for cinemaIdx, cinemaID := range cinemaIDs {
				// Не все фильмы во всех кинотеатрах
				if (movieIdx+cinemaIdx)%2 != 0 {
					continue
				}

				// Выбрать зал
				hallIdx := (movieIdx + cinemaIdx) % len(hallIDs)
				if hallIdx >= len(hallIDs) {
					continue
				}
				hallID := hallIDs[hallIdx]

				// Создать 2-3 сеанса в день
				for _, hour := range sessionTimes[:3] {
					startTime := time.Date(
						currentDate.Year(),
						currentDate.Month(),
						currentDate.Day(),
						hour, 0, 0, 0, time.Local,
					)

					// Пропустить прошедшие сеансы
					if startTime.Before(now) {
						continue
					}

					// Формат и цена
					format := "2D"
					basePrice := 2000.0
					switch movieIdx % 3 {
					case 0:
						format = "3D"
						basePrice = 2500.0
					case 1:
						format = "IMAX"
						basePrice = 3500.0
					}

					// Длительность фильма (примерно)
					duration := 120 + movieIdx*10
					endTime := startTime.Add(time.Duration(duration) * time.Minute)

					showtime := models.Showtime{
						MovieID:        movieID,
						CinemaID:       cinemaID,
						HallID:         hallID,
						StartTime:      startTime,
						EndTime:        endTime,
						BasePrice:      basePrice,
						Format:         format,
						Language:       "Russian",
						Subtitles:      "Kazakh",
						AvailableSeats: 150 - (movieIdx * 5), // симуляция занятости
						BookedSeats:    []models.BookedSeat{},
						CreatedAt:      time.Now(),
					}

					showtimes = append(showtimes, showtime)
				}
			}
		}
	}

	// Вставить сеансы
	for i := range showtimes {
		result, err := collection.InsertOne(ctx, showtimes[i])
		if err != nil {
			log.Printf("❌ Error inserting showtime: %v", err)
			continue
		}
		showtimeIDs = append(showtimeIDs, result.InsertedID.(primitive.ObjectID))
	}

	log.Printf("✅ Created %d showtimes", len(showtimeIDs))
}

// 6. BOOKINGS
func seedBookings(ctx context.Context) {
	collection := config.GetCollection("bookings")

	// Создать несколько тестовых броней
	if len(userIDs) < 3 || len(showtimeIDs) < 2 {
		log.Println("⚠️ Not enough data to create bookings")
		return
	}

	bookings := []models.Booking{
		{
			BookingNumber: "BK-20260201-001234",
			UserID:        userIDs[2],
			ShowtimeID:    showtimeIDs[0],
			Seats: []models.BookingSeat{
				{Row: "E", Number: 10, Price: 2000},
				{Row: "E", Number: 11, Price: 2000},
			},
			TotalAmount: 4000,
			Status:      "confirmed",
			Payment: models.Payment{
				Method:        "wallet",
				TransactionID: "TXN-" + time.Now().Format("20060102150405"),
				PaidAt:        time.Now().Add(-2 * time.Hour),
				Status:        "completed",
			},
			QRCode:    "QR-BK-20260201-001234",
			ExpiresAt: time.Now().Add(24 * time.Hour),
			CreatedAt: time.Now().Add(-2 * time.Hour),
			UpdatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			BookingNumber: "BK-20260201-001235",
			UserID:        userIDs[3],
			ShowtimeID:    showtimeIDs[1],
			Seats: []models.BookingSeat{
				{Row: "D", Number: 5, Price: 3000},
				{Row: "D", Number: 6, Price: 3000},
				{Row: "D", Number: 7, Price: 3000},
			},
			TotalAmount: 9000,
			Status:      "pending",
			Payment: models.Payment{
				Method: "card",
				Status: "pending",
			},
			QRCode:    "QR-BK-20260201-001235",
			ExpiresAt: time.Now().Add(15 * time.Minute),
			CreatedAt: time.Now().Add(-5 * time.Minute),
			UpdatedAt: time.Now().Add(-5 * time.Minute),
		},
	}

	for i := range bookings {
		result, err := collection.InsertOne(ctx, bookings[i])
		if err != nil {
			log.Printf("❌ Error inserting booking: %v", err)
			continue
		}
		_ = result.InsertedID.(primitive.ObjectID)
	}

	log.Printf("✅ Created %d bookings", len(bookings))
}

// 7. TRANSACTIONS
func seedTransactions(ctx context.Context) {
	collection := config.GetCollection("transactions")

	if len(userIDs) < 3 {
		log.Println("⚠️ Not enough users to create transactions")
		return
	}

	transactions := []models.Transaction{
		{
			UserID:      userIDs[2],
			Type:        "wallet_topup",
			Amount:      5000,
			Status:      "completed",
			Description: "Пополнение кошелька через карту",
			CreatedAt:   time.Now().Add(-48 * time.Hour),
		},
		{
			UserID:      userIDs[2],
			Type:        "booking",
			Amount:      -4000,
			Status:      "completed",
			Description: "Оплата билетов на 'Dune: Part Three'",
			CreatedAt:   time.Now().Add(-2 * time.Hour),
		},
		{
			UserID:      userIDs[3],
			Type:        "wallet_topup",
			Amount:      10000,
			Status:      "completed",
			Description: "Пополнение кошелька",
			CreatedAt:   time.Now().Add(-72 * time.Hour),
		},
	}

	for i := range transactions {
		_, err := collection.InsertOne(ctx, transactions[i])
		if err != nil {
			log.Printf("❌ Error inserting transaction: %v", err)
			continue
		}
	}

	log.Printf("✅ Created %d transactions", len(transactions))
}
