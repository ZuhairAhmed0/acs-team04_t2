package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/config"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/database"
	"github.com/ZuhairAhmed0/acs-team04_t2/Backend/app/security"
)

const defaultTestPassword = "Password123!"

func main() {
	clean := flag.Bool("clean", true, "Truncate existing tables before seeding")
	flag.Parse()

	cfg, err := config.Load(".env", "../.env")
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	if err := database.InitDB(cfg.DBDsn); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	if err := Seed(database.DB, *clean); err != nil {
		log.Fatalf("seeding failed: %v", err)
	}
}

// Seed populates the database with a full test dataset covering users, admins,
// stations, branches, wallets, transactions, and cash withdrawals.
func Seed(db *gorm.DB, clean bool) error {
	pwHash, err := security.HashPassword(defaultTestPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	now := time.Now().UTC()

	return db.Transaction(func(tx *gorm.DB) error {
		if clean {
			if err := tx.Exec(
				"TRUNCATE TABLE refresh_tokens, cash_withdrawals, transactions, wallets, branches, stations, users RESTART IDENTITY CASCADE",
			).Error; err != nil {
				return fmt.Errorf("truncate tables: %w", err)
			}
		}

		// 1. Seed Admins & Users (2 admins + 10 regular users = 12 users)
		users := []database.User{
			// Admins
			{Name: "مدير النظام", Email: "admin@fuelms.com", PasswordHash: pwHash, Phone: "+249912300001", Role: "admin", Status: "active", CreatedAt: now.AddDate(0, 0, -30)},
			{Name: "مشرف العمليات", Email: "supervisor@fuelms.com", PasswordHash: pwHash, Phone: "+249912300002", Role: "admin", Status: "active", CreatedAt: now.AddDate(0, 0, -28)},
			// 10 Regular Users (each starts with welcome loyalty points)
			{Name: "أحمد خليل", Email: "ahmed.khalil@example.com", PasswordHash: pwHash, Phone: "+249912345601", Role: "user", Status: "active", LoyaltyPoints: database.LoyaltyWelcomePoints, CreatedAt: now.AddDate(0, 0, -20)},
			{Name: "محمد نور", Email: "mohamed.nour@example.com", PasswordHash: pwHash, Phone: "+249912345602", Role: "user", Status: "active", LoyaltyPoints: database.LoyaltyWelcomePoints, CreatedAt: now.AddDate(0, 0, -19)},
			{Name: "فاطمة آدم", Email: "fatima.adam@example.com", PasswordHash: pwHash, Phone: "+249912345603", Role: "user", Status: "active", LoyaltyPoints: database.LoyaltyWelcomePoints, CreatedAt: now.AddDate(0, 0, -18)},
			{Name: "عمر حسن", Email: "omer.hassan@example.com", PasswordHash: pwHash, Phone: "+249912345604", Role: "user", Status: "active", LoyaltyPoints: database.LoyaltyWelcomePoints, CreatedAt: now.AddDate(0, 0, -17)},
			{Name: "سارة إبراهيم", Email: "sara.ibrahim@example.com", PasswordHash: pwHash, Phone: "+249912345605", Role: "user", Status: "active", LoyaltyPoints: database.LoyaltyWelcomePoints, CreatedAt: now.AddDate(0, 0, -16)},
			{Name: "علي عثمان", Email: "ali.osman@example.com", PasswordHash: pwHash, Phone: "+249912345606", Role: "user", Status: "active", LoyaltyPoints: database.LoyaltyWelcomePoints, CreatedAt: now.AddDate(0, 0, -15)},
			{Name: "منى صديق", Email: "mona.siddig@example.com", PasswordHash: pwHash, Phone: "+249912345607", Role: "user", Status: "active", LoyaltyPoints: database.LoyaltyWelcomePoints, CreatedAt: now.AddDate(0, 0, -14)},
			{Name: "خالد إدريس", Email: "khalid.idris@example.com", PasswordHash: pwHash, Phone: "+249912345608", Role: "user", Status: "active", LoyaltyPoints: database.LoyaltyWelcomePoints, CreatedAt: now.AddDate(0, 0, -13)},
			{Name: "زينب طارق", Email: "zainab.tariq@example.com", PasswordHash: pwHash, Phone: "+249912345609", Role: "user", Status: "active", LoyaltyPoints: database.LoyaltyWelcomePoints, CreatedAt: now.AddDate(0, 0, -12)},
			{Name: "مصطفى بكري", Email: "mustafa.bakri@example.com", PasswordHash: pwHash, Phone: "+249912345610", Role: "user", Status: "suspended", LoyaltyPoints: database.LoyaltyWelcomePoints, CreatedAt: now.AddDate(0, 0, -11)},
		}
		if err := tx.Create(&users).Error; err != nil {
			return fmt.Errorf("insert users: %w", err)
		}

		admin1ID := users[0].ID
		admin2ID := users[1].ID
		regularUsers := users[2:]

		// 2. Seed Stations (8 stations covering active, pending, rejected, revoked)
		reviewTime1 := now.AddDate(0, 0, -25)
		reviewTime2 := now.AddDate(0, 0, -22)
		rejectTime := now.AddDate(0, 0, -5)
		revokeTime := now.AddDate(0, 0, -3)
		rejectionReason := "انتهاء صلاحية السجل التجاري وعدم إمكانية التحقق من الضمان البنكي"

		stations := []database.Station{
			{
				Name: "شركة النيل للبترول", Location: "الخرطوم",
				BankAccountNumber: "100293847501", BankAccountName: "شركة النيل للبترول", BankName: "بنك الخرطوم",
				Email: "al-nile@petroleum.sd", PasswordHash: pwHash, Status: "active",
				ReviewedBy: &admin1ID, ReviewedAt: &reviewTime1, CreatedAt: now.AddDate(0, 0, -26),
			},
			{
				Name: "شركة بشائر للطاقة", Location: "الخرطوم",
				BankAccountNumber: "200481928374", BankAccountName: "شركة بشائر للطاقة", BankName: "بنك أم درمان الوطني",
				Email: "contact@bashair-energy.sd", PasswordHash: pwHash, Status: "active",
				ReviewedBy: &admin1ID, ReviewedAt: &reviewTime1, CreatedAt: now.AddDate(0, 0, -25),
			},
			{
				Name: "بتروتشاينا السودان", Location: "بحري",
				BankAccountNumber: "300192837465", BankAccountName: "بتروتشاينا السودان", BankName: "بنك فيصل الإسلامي",
				Email: "info@petrochina.sd", PasswordHash: pwHash, Status: "active",
				ReviewedBy: &admin2ID, ReviewedAt: &reviewTime2, CreatedAt: now.AddDate(0, 0, -24),
			},
			{
				Name: "النهضة لخدمات الوقود", Location: "أم درمان",
				BankAccountNumber: "100847291039", BankAccountName: "النهضة لخدمات الوقود", BankName: "بنك الخرطوم",
				Email: "support@alnahda-fuel.sd", PasswordHash: pwHash, Status: "active",
				ReviewedBy: &admin1ID, ReviewedAt: &reviewTime2, CreatedAt: now.AddDate(0, 0, -23),
			},
			{
				Name: "بترول البحر الأحمر", Location: "بورتسودان",
				BankAccountNumber: "300994827163", BankAccountName: "بترول البحر الأحمر", BankName: "بنك فيصل الإسلامي",
				Email: "operations@redsea-petro.sd", PasswordHash: pwHash, Status: "active",
				ReviewedBy: &admin2ID, ReviewedAt: &reviewTime2, CreatedAt: now.AddDate(0, 0, -21),
			},
			{
				Name: "الصحراء للنفط والغاز", Location: "الأبيض",
				BankAccountNumber: "200771829304", BankAccountName: "شركة الصحراء للنفط والغاز المحدودة", BankName: "بنك أم درمان الوطني",
				Email: "admin@sahara-oil.sd", PasswordHash: pwHash, Status: "pending",
				CreatedAt: now.AddDate(0, 0, -2),
			},
			{
				Name: "شركة النيل الأزرق للوقود", Location: "ود مدني",
				BankAccountNumber: "100662718293", BankAccountName: "شركة النيل الأزرق للوقود المحدودة", BankName: "بنك الخرطوم",
				Email: "info@bluenile-fuel.sd", PasswordHash: pwHash, Status: "rejected",
				ReviewedBy: &admin1ID, ReviewedAt: &rejectTime, RejectionReason: &rejectionReason,
				CreatedAt: now.AddDate(0, 0, -7),
			},
			{
				Name: "شركة تمازج للطاقة", Location: "كسلا",
				BankAccountNumber: "300882910293", BankAccountName: "شركة تمازج للطاقة", BankName: "بنك فيصل الإسلامي",
				Email: "contact@tamazuj-energy.sd", PasswordHash: pwHash, Status: "revoked",
				ReviewedBy: &admin2ID, ReviewedAt: &revokeTime, CreatedAt: now.AddDate(0, 0, -18),
			},
		}
		if err := tx.Create(&stations).Error; err != nil {
			return fmt.Errorf("insert stations: %w", err)
		}

		// 3. Seed 24 Branches across 7 Sudanese cities and all 3 fuel statuses (full, low, empty)
		branches := []database.Branch{
			// Station 0: شركة النيل للبترول (5 branches)
			{StationID: stations[0].ID, City: "الخرطوم", Address: "شارع أفريقيا (طريق المطار)، شارع 15", IsMain: true, FuelStatus: "full", CongestionPercentage: 65, CreatedAt: now.AddDate(0, 0, -26)},
			{StationID: stations[0].ID, City: "الخرطوم", Address: "شارع الستين، حي الرياض", IsMain: false, FuelStatus: "full", CongestionPercentage: 45, CreatedAt: now.AddDate(0, 0, -20)},
			{StationID: stations[0].ID, City: "أم درمان", Address: "شارع الموردة، مقابل النادي الثقافي", IsMain: false, FuelStatus: "low", CongestionPercentage: 80, CreatedAt: now.AddDate(0, 0, -18)},
			{StationID: stations[0].ID, City: "بحري", Address: "شارع المعونة، وسط بحري", IsMain: false, FuelStatus: "empty", CongestionPercentage: 15, CreatedAt: now.AddDate(0, 0, -15)},
			{StationID: stations[0].ID, City: "ود مدني", Address: "شارع كورنيش النيل، بالقرب من الاستاد", IsMain: false, FuelStatus: "full", CongestionPercentage: 35, CreatedAt: now.AddDate(0, 0, -12)},

			// Station 1: شركة بشائر للطاقة (5 branches)
			{StationID: stations[1].ID, City: "الخرطوم", Address: "المنشية، ناصية شارع النيل", IsMain: true, FuelStatus: "full", CongestionPercentage: 70, CreatedAt: now.AddDate(0, 0, -25)},
			{StationID: stations[1].ID, City: "الخرطوم", Address: "شارع الغابة، المنطقة الصناعية الثقيلة", IsMain: false, FuelStatus: "low", CongestionPercentage: 85, CreatedAt: now.AddDate(0, 0, -19)},
			{StationID: stations[1].ID, City: "أم درمان", Address: "شارع الشنقيطي، الثورة الحارة الخامسة", IsMain: false, FuelStatus: "full", CongestionPercentage: 50, CreatedAt: now.AddDate(0, 0, -16)},
			{StationID: stations[1].ID, City: "بحري", Address: "الصافية، شارع الإنقاذ", IsMain: false, FuelStatus: "full", CongestionPercentage: 40, CreatedAt: now.AddDate(0, 0, -14)},
			{StationID: stations[1].ID, City: "بورتسودان", Address: "طريق ميناء البحر الأحمر، الرصيف الجنوبي", IsMain: false, FuelStatus: "low", CongestionPercentage: 60, CreatedAt: now.AddDate(0, 0, -10)},

			// Station 2: بتروتشاينا السودان (4 branches)
			{StationID: stations[2].ID, City: "بحري", Address: "منطقة كوبر الصناعية، شارع 1", IsMain: true, FuelStatus: "full", CongestionPercentage: 55, CreatedAt: now.AddDate(0, 0, -24)},
			{StationID: stations[2].ID, City: "الخرطوم", Address: "حي الرياض، شارع المشتل", IsMain: false, FuelStatus: "full", CongestionPercentage: 62, CreatedAt: now.AddDate(0, 0, -17)},
			{StationID: stations[2].ID, City: "أم درمان", Address: "شارع الأربعين، تقاطع 4", IsMain: false, FuelStatus: "empty", CongestionPercentage: 20, CreatedAt: now.AddDate(0, 0, -13)},
			{StationID: stations[2].ID, City: "بورتسودان", Address: "الكورنيش الساحلي، حي الثورة", IsMain: false, FuelStatus: "full", CongestionPercentage: 48, CreatedAt: now.AddDate(0, 0, -9)},

			// Station 3: النهضة لخدمات الوقود (4 branches)
			{StationID: stations[3].ID, City: "أم درمان", Address: "تقاطع الطريق الرئيسي، سوق ليبيا", IsMain: true, FuelStatus: "full", CongestionPercentage: 75, CreatedAt: now.AddDate(0, 0, -23)},
			{StationID: stations[3].ID, City: "أم درمان", Address: "شارع الوادي، الثورة الحارة 21", IsMain: false, FuelStatus: "low", CongestionPercentage: 68, CreatedAt: now.AddDate(0, 0, -15)},
			{StationID: stations[3].ID, City: "الخرطوم", Address: "الصحافة زلط، المحور التجاري الرئيسي", IsMain: false, FuelStatus: "full", CongestionPercentage: 52, CreatedAt: now.AddDate(0, 0, -11)},
			{StationID: stations[3].ID, City: "كسلا", Address: "الطريق الدائري، وسط مدينة كسلا", IsMain: false, FuelStatus: "empty", CongestionPercentage: 10, CreatedAt: now.AddDate(0, 0, -8)},

			// Station 4: بترول البحر الأحمر (3 branches)
			{StationID: stations[4].ID, City: "بورتسودان", Address: "شارع المارينا، ميناء بورتسودان", IsMain: true, FuelStatus: "full", CongestionPercentage: 58, CreatedAt: now.AddDate(0, 0, -21)},
			{StationID: stations[4].ID, City: "بورتسودان", Address: "منطقة الجمارك، البوابة 7", IsMain: false, FuelStatus: "full", CongestionPercentage: 44, CreatedAt: now.AddDate(0, 0, -14)},
			{StationID: stations[4].ID, City: "كسلا", Address: "محطة الطريق السريع الحدودي، الكيلو 12", IsMain: false, FuelStatus: "low", CongestionPercentage: 72, CreatedAt: now.AddDate(0, 0, -7)},

			// Station 5: الصحراء للنفط والغاز (Pending - 1 branch)
			{StationID: stations[5].ID, City: "الأبيض", Address: "طريق كردفان السريع، الدائري الأبيض", IsMain: true, FuelStatus: "empty", CongestionPercentage: 0, CreatedAt: now.AddDate(0, 0, -2)},

			// Station 6: شركة النيل الأزرق للوقود (Rejected - 1 branch)
			{StationID: stations[6].ID, City: "ود مدني", Address: "تقاطع السوق الكبير، وسط المدينة", IsMain: true, FuelStatus: "empty", CongestionPercentage: 0, CreatedAt: now.AddDate(0, 0, -7)},

			// Station 7: شركة تمازج للطاقة (Revoked - 1 branch)
			{StationID: stations[7].ID, City: "كسلا", Address: "ضفة نهر القاش، القطاع الشمالي", IsMain: true, FuelStatus: "low", CongestionPercentage: 25, CreatedAt: now.AddDate(0, 0, -18)},
		}
		if err := tx.Create(&branches).Error; err != nil {
			return fmt.Errorf("insert branches: %w", err)
		}

		// 4. Build Wallets & Reconciled Transactions
		userWallets := make(map[uint64]*database.Wallet, len(users))
		stationWallets := make(map[uint64]*database.Wallet, len(stations))

		for _, u := range users {
			userWallets[u.ID] = &database.Wallet{
				OwnerType: "user",
				OwnerID:   u.ID,
				CreatedAt: u.CreatedAt,
			}
		}
		for _, s := range stations {
			stationWallets[s.ID] = &database.Wallet{
				OwnerType: "station",
				OwnerID:   s.ID,
				CreatedAt: s.CreatedAt,
			}
		}

		var txs []database.Transaction
		refSeq := 100001
		nextRef := func(prefix string, ts time.Time) string {
			ref := fmt.Sprintf("%s-%d-%06d", prefix, ts.Unix(), refSeq)
			refSeq++
			return ref
		}

		userType := "user"
		stationType := "station"

		userPoints := make(map[uint64]int64, len(regularUsers))
		for _, u := range regularUsers {
			userPoints[u.ID] = u.LoyaltyPoints
		}

		recordTopup := func(u database.User, amount int64, ts time.Time) {
			w := userWallets[u.ID]
			w.AvailableBalance += amount
			w.TotalBalance += amount
			userPoints[u.ID] += database.LoyaltyActivityPoints
			txs = append(txs, database.Transaction{
				Reference: nextRef("TOP", ts),
				Type:      "topup",
				Amount:    amount,
				ToType:    &userType,
				ToID:      &u.ID,
				ToName:    u.Name,
				Status:    "completed",
				CreatedAt: ts,
			})
		}

		recordAllocation := func(u database.User, s database.Station, amount int64, branch *database.Branch, ts time.Time) {
			uw := userWallets[u.ID]
			uw.AvailableBalance -= amount
			uw.AllocatedBalance += amount

			sw := stationWallets[s.ID]
			sw.AvailableBalance += amount
			sw.TotalBalance += amount

			var bID *uint64
			var bName *string
			if branch != nil {
				id := branch.ID
				name := fmt.Sprintf("%s - %s", branch.City, branch.Address)
				bID = &id
				bName = &name
			}

			txs = append(txs, database.Transaction{
				Reference:  nextRef("TRF", ts),
				Type:       "allocation",
				Amount:     amount,
				FromType:   &userType,
				FromID:     &u.ID,
				FromName:   u.Name,
				ToType:     &stationType,
				ToID:       &s.ID,
				ToName:     s.Name,
				BranchID:   bID,
				BranchName: bName,
				Status:     "completed",
				CreatedAt:  ts,
			})
		}

		recordStationTransfer := func(fromS, toS database.Station, amount int64, ts time.Time) {
			fw := stationWallets[fromS.ID]
			fw.AvailableBalance -= amount
			fw.TotalBalance -= amount

			tw := stationWallets[toS.ID]
			tw.AvailableBalance += amount
			tw.TotalBalance += amount

			txs = append(txs, database.Transaction{
				Reference: nextRef("STF", ts),
				Type:      "station_transfer",
				Amount:    amount,
				FromType:  &stationType,
				FromID:    &fromS.ID,
				FromName:  fromS.Name,
				ToType:    &stationType,
				ToID:      &toS.ID,
				ToName:    toS.Name,
				Status:    "completed",
				CreatedAt: ts,
			})
		}

		recordWithdrawalTx := func(s database.Station, amount int64, ts time.Time) {
			sw := stationWallets[s.ID]
			sw.AvailableBalance -= amount
			sw.TotalBalance -= amount

			txs = append(txs, database.Transaction{
				Reference: nextRef("WTH", ts),
				Type:      "withdrawal",
				Amount:    amount,
				FromType:  &stationType,
				FromID:    &s.ID,
				FromName:  s.Name,
				Status:    "completed",
				CreatedAt: ts,
			})
		}

		// Top-ups for all 10 regular users
		topupAmounts := []int64{150000, 280000, 75000, 210000, 95000, 40000, 120000, 50000, 300000, 10000}
		for i, u := range regularUsers {
			recordTopup(u, topupAmounts[i], now.AddDate(0, 0, -10).Add(time.Duration(i)*time.Hour))
		}

		// Allocations from users to active stations
		recordAllocation(regularUsers[0], stations[0], 35000, &branches[0], now.AddDate(0, 0, -9))
		recordAllocation(regularUsers[0], stations[1], 25000, nil, now.AddDate(0, 0, -8))
		recordAllocation(regularUsers[1], stations[0], 60000, &branches[1], now.AddDate(0, 0, -8).Add(3*time.Hour))
		recordAllocation(regularUsers[1], stations[2], 40000, nil, now.AddDate(0, 0, -7))
		recordAllocation(regularUsers[2], stations[1], 30000, &branches[5], now.AddDate(0, 0, -6))
		recordAllocation(regularUsers[3], stations[3], 35000, nil, now.AddDate(0, 0, -6).Add(4*time.Hour))
		recordAllocation(regularUsers[3], stations[4], 25000, &branches[18], now.AddDate(0, 0, -5))
		recordAllocation(regularUsers[4], stations[2], 30000, nil, now.AddDate(0, 0, -4))
		recordAllocation(regularUsers[5], stations[0], 15000, &branches[2], now.AddDate(0, 0, -4).Add(2*time.Hour))
		recordAllocation(regularUsers[6], stations[1], 40000, nil, now.AddDate(0, 0, -3))
		recordAllocation(regularUsers[7], stations[4], 15000, nil, now.AddDate(0, 0, -2))
		recordAllocation(regularUsers[8], stations[0], 50000, &branches[4], now.AddDate(0, 0, -1))
		recordAllocation(regularUsers[8], stations[1], 30000, nil, now.Add(-12*time.Hour))

		// Station-to-station transfers
		recordStationTransfer(stations[0], stations[1], 15000, now.AddDate(0, 0, -5).Add(5*time.Hour))
		recordStationTransfer(stations[1], stations[2], 10000, now.AddDate(0, 0, -3).Add(6*time.Hour))
		recordStationTransfer(stations[2], stations[3], 5000, now.AddDate(0, 0, -1).Add(4*time.Hour))

		// Completed cash withdrawal transactions
		wth1Time := now.AddDate(0, 0, -5)
		wth2Time := now.AddDate(0, 0, -3)
		recordWithdrawalTx(stations[0], 25000, wth1Time)
		recordWithdrawalTx(stations[1], 15000, wth2Time)

		// Persist all wallets
		allWallets := make([]database.Wallet, 0, len(users)+len(stations))
		for _, u := range users {
			allWallets = append(allWallets, *userWallets[u.ID])
		}
		for _, u := range regularUsers {
			if err := tx.Model(&database.User{}).Where("id = ?", u.ID).Update("loyalty_points", userPoints[u.ID]).Error; err != nil {
				return fmt.Errorf("update user loyalty points: %w", err)
			}
		}
		for _, s := range stations {
			allWallets = append(allWallets, *stationWallets[s.ID])
		}
		if err := tx.Create(&allWallets).Error; err != nil {
			return fmt.Errorf("insert wallets: %w", err)
		}

		// Persist all transactions
		if err := tx.Create(&txs).Error; err != nil {
			return fmt.Errorf("insert transactions: %w", err)
		}

		// 5. Seed Cash Withdrawals (covering completed, approved, pending, rejected)
		approvedTime := now.AddDate(0, 0, -1)
		rejectedWthTime := now.AddDate(0, 0, -2)
		wthRejectReason := "المبلغ المطلوب يتجاوز الحد اليومي للسحب النقدي للمحطة"
		extRef1 := "BOK-20261109-441820"
		extRef2 := "ONB-20261111-903154"

		withdrawals := []database.CashWithdrawal{
			{
				StationID: stations[0].ID, Amount: 25000, Status: "completed", ExternalReference: &extRef1,
				ReviewedBy: &admin1ID, ReviewedAt: &wth1Time, CreatedAt: wth1Time.Add(-4 * time.Hour),
			},
			{
				StationID: stations[1].ID, Amount: 15000, Status: "completed", ExternalReference: &extRef2,
				ReviewedBy: &admin1ID, ReviewedAt: &wth2Time, CreatedAt: wth2Time.Add(-3 * time.Hour),
			},
			{
				StationID: stations[0].ID, Amount: 20000, Status: "approved",
				ReviewedBy: &admin2ID, ReviewedAt: &approvedTime, CreatedAt: approvedTime.Add(-5 * time.Hour),
			},
			{
				StationID: stations[1].ID, Amount: 10000, Status: "pending",
				CreatedAt: now.Add(-6 * time.Hour),
			},
			{
				StationID: stations[2].ID, Amount: 15000, Status: "pending",
				CreatedAt: now.Add(-2 * time.Hour),
			},
			{
				StationID: stations[3].ID, Amount: 50000, Status: "rejected",
				ReviewedBy: &admin1ID, ReviewedAt: &rejectedWthTime, RejectionReason: &wthRejectReason,
				CreatedAt: rejectedWthTime.Add(-6 * time.Hour),
			},
		}
		if err := tx.Create(&withdrawals).Error; err != nil {
			return fmt.Errorf("insert cash withdrawals: %w", err)
		}

		fmt.Printf("Seeded successfully!\n")
		fmt.Printf("  - Users: %d (2 admins, %d regular users)\n", len(users), len(regularUsers))
		fmt.Printf("  - Stations: %d (5 active, 1 pending, 1 rejected, 1 revoked)\n", len(stations))
		fmt.Printf("  - Branches: %d\n", len(branches))
		fmt.Printf("  - Wallets: %d\n", len(allWallets))
		fmt.Printf("  - Transactions: %d\n", len(txs))
		fmt.Printf("  - Cash Withdrawals: %d\n", len(withdrawals))
		fmt.Printf("\nDefault password for all accounts: %s\n", defaultTestPassword)
		fmt.Printf("  Admin login:   admin@fuelms.com\n")
		fmt.Printf("  User login:    ahmed.khalil@example.com\n")
		fmt.Printf("  Station login: al-nile@petroleum.sd\n")
		return nil
	})
}
