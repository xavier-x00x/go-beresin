package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"go-beresin/internal/domain"
	"go-beresin/internal/repository"
	"go-beresin/pkg/database"
)

func main() {
	// Load .env file
	_ = godotenv.Load()

	ctx := context.Background()

	// Initialize database pool
	pool, err := database.InitPool(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer pool.Close()

	// Initialize repository
	queries := repository.New(pool)

	log.Println("[INFO] Clearing existing data (clean seed)...")
	// Clean DB using raw exec
	_, err = pool.Exec(ctx, `
		TRUNCATE service_packages, service_listings, categories, talent_profiles, users CASCADE;
	`)
	if err != nil {
		log.Fatalf("Failed to truncate tables: %v", err)
	}
	log.Println("[SUCCESS] Database tables truncated.")

	log.Println("[INFO] Seeding Main Categories...")
	categoriesData := []struct {
		Name      string
		Slug      string
		SortOrder int32
	}{
		{"Rumah & Properti", "rumah-properti", 1},
		{"Event & Entertainment", "event-entertainment", 2},
		{"Kesehatan & Konseling", "kesehatan-konseling", 3},
		{"Pendidikan", "pendidikan", 4},
		{"Kecantikan", "kecantikan", 5},
	}

	catMap := make(map[string]pgtype.UUID)

	for _, c := range categoriesData {
		cat, err := queries.CreateCategory(ctx, repository.CreateCategoryParams{
			ID:        domain.NewUUIDV7(),
			Name:      c.Name,
			Slug:      c.Slug,
			SortOrder: pgtype.Int4{Int32: c.SortOrder, Valid: true},
			IsActive:  pgtype.Bool{Bool: true, Valid: true},
		})
		if err != nil {
			log.Fatalf("Failed to create category %s: %v", c.Name, err)
		}
		catMap[c.Slug] = cat.ID
		log.Printf("[SUCCESS] Category created: %s (%s)", cat.Name, cat.Slug)
	}

	log.Println("[INFO] Seeding Sub-Categories...")
	subCategoriesData := []struct {
		Name      string
		Slug      string
		Parent    string
		SortOrder int32
	}{
		// Rumah & Properti
		{"Cleaning Service", "cleaning-service", "rumah-properti", 1},
		{"Kelistingan & AC", "kelistrikan-ac", "rumah-properti", 2},
		{"Renovasi Rumah", "renovasi-rumah", "rumah-properti", 3},
		// Event & Entertainment
		{"Fotografi & Videografi", "fotografi-videografi", "event-entertainment", 1},
		{"Katering Makanan", "katering-makanan", "event-entertainment", 2},
		{"Dekorasi Panggung", "dekorasi-panggung", "event-entertainment", 3},
		// Kesehatan & Konseling
		{"Fisioterapi Rumah", "fisioterapi-rumah", "kesehatan-konseling", 1},
		{"Psikologi & Konseling", "psikologi-konseling", "kesehatan-konseling", 2},
		// Pendidikan
		{"Les Privat Matematika", "les-privat-matematika", "pendidikan", 1},
		// Kecantikan
		{"Makeup Artist (MUA)", "makeup-artist-mua", "kecantikan", 1},
	}

	subCatMap := make(map[string]pgtype.UUID)

	for _, sc := range subCategoriesData {
		parentID, exists := catMap[sc.Parent]
		if !exists {
			log.Fatalf("Parent category %s not found for %s", sc.Parent, sc.Name)
		}

		subCat, err := queries.CreateCategory(ctx, repository.CreateCategoryParams{
			ID:        domain.NewUUIDV7(),
			Name:      sc.Name,
			Slug:      sc.Slug,
			ParentID:  parentID,
			SortOrder: pgtype.Int4{Int32: sc.SortOrder, Valid: true},
			IsActive:  pgtype.Bool{Bool: true, Valid: true},
		})
		if err != nil {
			log.Fatalf("Failed to create sub-category %s: %v", sc.Name, err)
		}
		subCatMap[sc.Slug] = subCat.ID
		log.Printf("[SUCCESS] Sub-category created: %s -> %s", sc.Name, sc.Parent)
	}

	log.Println("[INFO] Seeding Dummy Users...")
	usersData := []struct {
		FullName string
		Email    string
		Phone    string
		City     string
		Role     string
	}{
		{"Budi Santoso", "budi@gmail.com", "081234567890", "Jakarta", "user"},
		{"Siti Aminah", "siti@gmail.com", "081234567891", "Tangerang", "user"},
		{"Andi Wijaya", "andi@gmail.com", "081234567892", "Bandung", "user"},
	}

	for _, u := range usersData {
		_, err := queries.CreateUser(ctx, repository.CreateUserParams{
			ID:           domain.NewUUIDV7(),
			FullName:     u.FullName,
			Email:        pgtype.Text{String: u.Email, Valid: true},
			Phone:        u.Phone,
			PasswordHash: pgtype.Text{String: "$2a$12$L7R2QfU8K1i4tG/f1j3n1e8E0KqgR17gR0t7.9zD1bO8G4Q7c1X/e", Valid: true}, // dummy hash
			Role:         u.Role,
			City:         pgtype.Text{String: u.City, Valid: true},
			IsVerified:   pgtype.Bool{Bool: true, Valid: true},
			IsActive:     pgtype.Bool{Bool: true, Valid: true},
		})
		if err != nil {
			log.Fatalf("Failed to create user %s: %v", u.FullName, err)
		}
		log.Printf("[SUCCESS] User created: %s (%s)", u.FullName, u.Role)
	}

	log.Println("[INFO] Seeding Dummy Talents & Profiles (with PostGIS locations)...")
	talentsData := []struct {
		FullName  string
		Email     string
		Phone     string
		City      string
		Tagline   string
		Bio       string
		ExpYears  int32
		RadiusKm  int32
		Longitude float64
		Latitude  float64
		Listing   struct {
			Title       string
			SubCatSlug  string
			Description string
			Package     struct {
				Name        string
				Price       int64
				DurationHrs int32
				Description string
			}
		}
	}{
		{
			FullName:  "Joko Susilo",
			Email:     "joko@gmail.com",
			Phone:     "081234567893",
			City:      "Jakarta",
			Tagline:   "Spesialis Kelistrikan & Pendingin Ruangan Panggilan",
			Bio:       "Berpengalaman memperbaiki berbagai AC rumah, kantor, dan instalasi kelistrikan rumah tangga sejak 2015. Cepat, amanah, dan bergaransi.",
			ExpYears:  9,
			RadiusKm:  25,
			Longitude: 106.8456, // Bundaran HI, Jakarta Pusat
			Latitude:  -6.2088,
			Listing: struct {
				Title       string
				SubCatSlug  string
				Description string
				Package     struct {
					Name        string
					Price       int64
					DurationHrs int32
					Description string
				}
			}{
				Title:      "Jasa Service AC & Kelistrikan Rumah Panggilan",
				SubCatSlug: "kelistrikan-ac",
				Description: "Kami melayani service cuci AC, tambah freon, perbaikan AC bocor atau tidak dingin, serta pemasangan/perbaikan instalasi listrik rumah Anda secara aman.",
				Package: struct {
					Name        string
					Price       int64
					DurationHrs int32
					Description string
				}{
					Name:        "Service Cuci AC Standard",
					Price:       150000,
					DurationHrs: 2,
					Description: "Cuci AC indoor & outdoor, cek tekanan freon, pembersihan filter, dan cek kelistrikan standar AC.",
				},
			},
		},
		{
			FullName:  "Dewi Lestari",
			Email:     "dewi@gmail.com",
			Phone:     "081234567894",
			City:      "Tangerang",
			Tagline:   "Professional Makeup Artist (MUA) Wisuda, Wedding & Event",
			Bio:       "Lulusan sekolah kecantikan ternama. Melayani makeup panggilan untuk wisuda, lamaran, prewedding, wedding, dan photoshoot. Menggunakan kosmetik premium.",
			ExpYears:  5,
			RadiusKm:  15,
			Longitude: 106.6319, // Tangerang Kota
			Latitude:  -6.1783,
			Listing: struct {
				Title       string
				SubCatSlug  string
				Description string
				Package     struct {
					Name        string
					Price       int64
					DurationHrs int32
					Description string
				}
			}{
				Title:      "Jasa MUA & Hairdo/Hijab Panggilan Professional",
				SubCatSlug: "makeup-artist-mua",
				Description: "Dapatkan riasan wajah flawless, glowing, dan tahan lama untuk hari spesial Anda. Layanan sudah termasuk hairdo atau hijab styling di lokasi Anda.",
				Package: struct {
					Name        string
					Price       int64
					DurationHrs int32
					Description string
				}{
					Name:        "Makeup & Hairdo Wisuda / Lamaran",
					Price:       500000,
					DurationHrs: 3,
					Description: "Flawless makeup dengan kosmetik premium, hairdo/hijab styling, dan fake eyelashes standar.",
				},
			},
		},
		{
			FullName:  "Rian Hidayat",
			Email:     "rian@gmail.com",
			Phone:     "081234567895",
			City:      "Bandung",
			Tagline:   "Fotografer & Videografer Event Kreatif",
			Bio:       "Menyediakan jasa dokumentasi foto dan video cinematic untuk pernikahan, event corporate, ulang tahun, dan commercial product. Memiliki tim solid dan gear profesional.",
			ExpYears:  7,
			RadiusKm:  30,
			Longitude: 107.6191, // Gedung Sate, Bandung
			Latitude:  -6.9175,
			Listing: struct {
				Title       string
				SubCatSlug  string
				Description string
				Package     struct {
					Name        string
					Price       int64
					DurationHrs int32
					Description string
				}
			}{
				Title:      "Jasa Dokumentasi Foto & Video Cinematic Pernikahan / Event",
				SubCatSlug: "fotografi-videografi",
				Description: "Abadikan momen spesial Anda dengan kualitas gambar cinematic yang dramatis dan memukau. Tim profesional yang berpengalaman siap membantu konsep Anda.",
				Package: struct {
					Name        string
					Price       int64
					DurationHrs int32
					Description string
				}{
					Name:        "Dokumentasi Event Half-Day",
					Price:       1500000,
					DurationHrs: 5,
					Description: "1 Fotografer, 1 Videografer, 50 foto editan, 1 menit video teaser cinematic media sosial, seluruh file mentahan dikirim via Google Drive.",
				},
			},
		},
	}

	for _, t := range talentsData {
		// 1. Create User
		user, err := queries.CreateUser(ctx, repository.CreateUserParams{
			ID:           domain.NewUUIDV7(),
			FullName:     t.FullName,
			Email:        pgtype.Text{String: t.Email, Valid: true},
			Phone:        t.Phone,
			PasswordHash: pgtype.Text{String: "$2a$12$L7R2QfU8K1i4tG/f1j3n1e8E0KqgR17gR0t7.9zD1bO8G4Q7c1X/e", Valid: true},
			Role:         "talent",
			City:         pgtype.Text{String: t.City, Valid: true},
			IsVerified:   pgtype.Bool{Bool: true, Valid: true},
			IsActive:     pgtype.Bool{Bool: true, Valid: true},
		})
		if err != nil {
			log.Fatalf("Failed to create talent user %s: %v", t.FullName, err)
		}

		// 2. Create Talent Profile with PostGIS coordinates
		tp, err := queries.CreateTalentProfile(ctx, repository.CreateTalentProfileParams{
			ID:               domain.NewUUIDV7(),
			UserID:           user.ID,
			Bio:              pgtype.Text{String: t.Bio, Valid: true},
			Tagline:          pgtype.Text{String: t.Tagline, Valid: true},
			YearsExperience:  pgtype.Int4{Int32: t.ExpYears, Valid: true},
			ServiceRadiusKm:  pgtype.Int4{Int32: t.RadiusKm, Valid: true},
			IsKycVerified:    pgtype.Bool{Bool: true, Valid: true},
			SubscriptionTier: pgtype.Text{String: "pro", Valid: true},
			AverageRating:    pgtype.Numeric{Valid: true, Int: nil}, // nil is handled or let default
			Longitude:        t.Longitude,
			Latitude:         t.Latitude,
		})
		if err != nil {
			log.Fatalf("Failed to create talent profile for %s: %v", t.FullName, err)
		}

		log.Printf("[SUCCESS] Talent Profile created: %s (Coord: %f, %f)", t.FullName, t.Longitude, t.Latitude)

		// 3. Create Service Listing
		subCatID, exists := subCatMap[t.Listing.SubCatSlug]
		if !exists {
			log.Fatalf("Sub-category %s not found for listing", t.Listing.SubCatSlug)
		}

		listing, err := queries.CreateServiceListing(ctx, repository.CreateServiceListingParams{
			ID:          domain.NewUUIDV7(),
			TalentID:    tp.ID,
			CategoryID:  subCatID,
			Title:       t.Listing.Title,
			Description: t.Listing.Description,
			Tags:        []string{"jasa", "professional", t.Listing.SubCatSlug},
			Status:      pgtype.Text{String: "published", Valid: true},
		})
		if err != nil {
			log.Fatalf("Failed to create service listing for %s: %v", t.FullName, err)
		}

		log.Printf("[SUCCESS] Service Listing created: '%s'", listing.Title)

		// 4. Create Service Package
		_, err = queries.CreateServicePackage(ctx, repository.CreateServicePackageParams{
			ID:            domain.NewUUIDV7(),
			ListingID:     listing.ID,
			Name:          t.Listing.Package.Name,
			Description:   pgtype.Text{String: t.Listing.Package.Description, Valid: true},
			PriceAmount:   t.Listing.Package.Price,
			PriceType:     pgtype.Text{String: "fixed", Valid: true},
			DurationHours: pgtype.Int4{Int32: t.Listing.Package.DurationHrs, Valid: true},
			Inclusions:    []string{"Tenaga Kerja Ahli", "Alat & Transportasi"},
			MaxRevisions:  pgtype.Int4{Int32: 3, Valid: true},
			SortOrder:     pgtype.Int4{Int32: 1, Valid: true},
		})
		if err != nil {
			log.Fatalf("Failed to create service package for listing %s: %v", listing.Title, err)
		}

		log.Printf("[SUCCESS] Service Package created: '%s' (Rp %d)", t.Listing.Package.Name, t.Listing.Package.Price)
	}

	log.Println("\n[SUCCESS] Seeding database completed successfully!")
}
