package bootstrap

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Row struct {
	UID          string
	Name         string
	ShortName    string
	NHFRID       string
	SubcountyUID string
	Subcounty    string
	DistrictUID  string
	District     string
	Region       string
	Level        string
	Ownership    string
	Status       string
}

func SeedFacilities(ctx context.Context, db *pgxpool.Pool, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1

	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	facilityLevelMap, err := loadFacilityLevelMap(ctx, db)
	if err != nil {
		return err
	}

	for i, r := range rows {
		if i == 0 {
			continue
		}

		if len(r) < 16 {
			return fmt.Errorf("row %d has insufficient columns: got %d", i+1, len(r))
		}

		row := Row{
			UID:          strings.TrimSpace(r[1]),
			Name:         strings.TrimSpace(r[2]),
			ShortName:    strings.TrimSpace(r[3]),
			NHFRID:       strings.TrimSpace(r[4]),
			SubcountyUID: strings.TrimSpace(r[5]),
			Subcounty:    strings.TrimSpace(r[6]),
			DistrictUID:  strings.TrimSpace(r[9]),
			District:     strings.TrimSpace(r[10]),
			Region:       strings.TrimSpace(r[12]),
			Level:        strings.TrimSpace(r[13]),
			Ownership:    strings.TrimSpace(r[14]),
			Status:       strings.TrimSpace(r[15]),
		}

		_, err := db.Exec(ctx, `
			INSERT INTO ref_districts (code, name, region)
			VALUES ($1, $2, $3)
			ON CONFLICT (name) DO NOTHING
		`, row.DistrictUID, row.District, row.Region)
		if err != nil {
			return fmt.Errorf("insert district %q: %w", row.District, err)
		}

		var districtID string
		err = db.QueryRow(ctx, `
			SELECT id
			FROM ref_districts
			WHERE name = $1
		`, row.District).Scan(&districtID)
		if err != nil {
			return fmt.Errorf("fetch district id for %q: %w", row.District, err)
		}

		_, err = db.Exec(ctx, `
			INSERT INTO ref_subcounties (code, name, district_id)
			VALUES ($1, $2, $3)
			ON CONFLICT (district_id, name) DO NOTHING
		`, row.SubcountyUID, row.Subcounty, districtID)
		if err != nil {
			return fmt.Errorf("insert subcounty %q: %w", row.Subcounty, err)
		}

		var subcountyID string
		err = db.QueryRow(ctx, `
			SELECT id
			FROM ref_subcounties
			WHERE name = $1 AND district_id = $2
		`, row.Subcounty, districtID).Scan(&subcountyID)
		if err != nil {
			return fmt.Errorf("fetch subcounty id for %q: %w", row.Subcounty, err)
		}

		var levelID any
		levelCode := normalizeFacilityLevelCode(row.Level)
		if levelCode != "" {
			if id, ok := facilityLevelMap[levelCode]; ok {
				levelID = id
			} else {
				levelID = nil
				fmt.Printf("warning: level code %q mapped from %q not found in ref_facility_levels for facility %q\n", levelCode, row.Level, row.Name)
			}
		} else {
			levelID = nil
			fmt.Printf("warning: unmapped facility level %q for facility %q\n", row.Level, row.Name)
		}

		nhfrID := normalizeNHFR(row.NHFRID)
		isActive := strings.EqualFold(row.Status, "Functional")

		if err := upsertFacility(ctx, db, row, nhfrID, districtID, subcountyID, levelID, isActive); err != nil {
			return fmt.Errorf("upsert facility %q: %w", row.Name, err)
		}
	}

	fmt.Println("facility seeding completed")
	return nil
}

func upsertFacility(
	ctx context.Context,
	db *pgxpool.Pool,
	row Row,
	nhfrID string,
	districtID string,
	subcountyID string,
	levelID any,
	isActive bool,
) error {
	tag, err := db.Exec(ctx, `
		UPDATE ref_facilities
		SET
			name = $2,
			short_name = $3,
			nhfr_id = $4,
			district_id = $5,
			subcounty_id = $6,
			level_id = $7,
			ownership = $8,
			is_active = $9,
			updated_at = now()
		WHERE code = $1
	`,
		row.UID,
		row.Name,
		row.ShortName,
		nhfrID,
		districtID,
		subcountyID,
		levelID,
		row.Ownership,
		isActive,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() > 0 {
		return nil
	}

	tag, err = db.Exec(ctx, `
		UPDATE ref_facilities
		SET
			code = $1,
			short_name = $3,
			nhfr_id = $4,
			subcounty_id = $6,
			level_id = $7,
			ownership = $8,
			is_active = $9,
			updated_at = now()
		WHERE name = $2
		  AND district_id = $5
	`,
		row.UID,
		row.Name,
		row.ShortName,
		nhfrID,
		districtID,
		subcountyID,
		levelID,
		row.Ownership,
		isActive,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() > 0 {
		return nil
	}

	_, err = db.Exec(ctx, `
		INSERT INTO ref_facilities
		(code, name, short_name, nhfr_id, district_id, subcounty_id, level_id, ownership, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`,
		row.UID,
		row.Name,
		row.ShortName,
		nhfrID,
		districtID,
		subcountyID,
		levelID,
		row.Ownership,
		isActive,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			return fmt.Errorf("%s (%s)", pgErr.Message, pgErr.ConstraintName)
		}
		return err
	}

	return nil
}

func normalizeNHFR(v string) string {
	v = strings.TrimSpace(v)

	if v == "" {
		return ""
	}

	if strings.Contains(strings.ToUpper(v), "E+") {
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return fmt.Sprintf("%.0f", f)
		}
	}

	return v
}

func loadFacilityLevelMap(ctx context.Context, db *pgxpool.Pool) (map[string]string, error) {
	rows, err := db.Query(ctx, `
		SELECT id, code
		FROM ref_facility_levels
		WHERE is_active = true
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	levelMap := make(map[string]string)
	for rows.Next() {
		var id, code string
		if err := rows.Scan(&id, &code); err != nil {
			return nil, err
		}
		levelMap[strings.ToUpper(strings.TrimSpace(code))] = id
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return levelMap, nil
}

func normalizeFacilityLevelCode(src string) string {
	v := strings.ToUpper(strings.TrimSpace(src))

	switch v {
	case "HC II", "HCII":
		return "HCII"
	case "HC III", "HCIII":
		return "HCIII"
	case "HC IV", "HCIV":
		return "HCIV"
	case "GENERAL HOSPITAL", "HOSPITAL":
		return "HOSPITAL"
	case "RRH":
		return "RRH"
	case "NRH":
		return "NRH"
	case "BCDP":
		return "BCDP"
	case "CLINIC":
		return "CLINIC"
	case "DRUG SHOP", "DRUGSHOP":
		return "DRUGSHOP"
	case "NBB":
		return "NBB"
	case "RBB":
		return "RBB"
	default:
		return ""
	}
}
