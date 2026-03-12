package bootstrap

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

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

	for i, r := range rows {

		if i == 0 {
			continue
		}

		row := Row{
			UID:          r[1],
			Name:         r[2],
			ShortName:    r[3],
			NHFRID:       r[4],
			SubcountyUID: r[5],
			Subcounty:    r[6],
			DistrictUID:  r[9],
			District:     r[10],
			Region:       r[12],
			Level:        r[13],
			Ownership:    r[14],
			Status:       r[15],
		}

		// -------------------------
		// insert district
		// -------------------------

		_, err := db.Exec(ctx, `
		INSERT INTO ref_districts (code,name,region)
		VALUES ($1,$2,$3)
		ON CONFLICT (name) DO NOTHING
		`, row.DistrictUID, row.District, row.Region)

		if err != nil {
			return err
		}

		// -------------------------
		// fetch district id
		// -------------------------

		var districtID string

		err = db.QueryRow(ctx, `
		SELECT id FROM ref_districts WHERE name=$1
		`, row.District).Scan(&districtID)

		if err != nil {
			return err
		}

		// -------------------------
		// insert subcounty
		// -------------------------

		_, err = db.Exec(ctx, `
		INSERT INTO ref_subcounties (code,name,district_id)
		VALUES ($1,$2,$3)
		ON CONFLICT (district_id,name) DO NOTHING
		`,
			row.SubcountyUID,
			row.Subcounty,
			districtID,
		)

		if err != nil {
			return err
		}

		// -------------------------
		// fetch subcounty id
		// -------------------------

		var subcountyID string

		err = db.QueryRow(ctx, `
		SELECT id FROM ref_subcounties
		WHERE name=$1 AND district_id=$2
		`, row.Subcounty, districtID).Scan(&subcountyID)

		if err != nil {
			return err
		}

		// -------------------------
		// facility level
		// -------------------------

		var levelID *string

		err = db.QueryRow(ctx, `
		SELECT id FROM ref_facility_levels WHERE name=$1
		`, row.Level).Scan(&levelID)

		if err != nil {
			levelID = nil
		}

		// -------------------------
		// insert facility
		// -------------------------

		_, err = db.Exec(ctx, `
		INSERT INTO ref_facilities
		(code, name, short_name, nhfr_id, district_id, subcounty_id, level_id, ownership, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (code) DO UPDATE SET
			name = EXCLUDED.name,
			short_name = EXCLUDED.short_name,
			nhfr_id = EXCLUDED.nhfr_id,
			district_id = EXCLUDED.district_id,
			subcounty_id = EXCLUDED.subcounty_id,
			level_id = EXCLUDED.level_id,
			ownership = EXCLUDED.ownership,
			is_active = EXCLUDED.is_active,
			updated_at = now()
		`,
			row.UID,
			row.Name,
			row.ShortName,
			normalizeNHFR(row.NHFRID),
			districtID,
			subcountyID,
			levelID,
			row.Ownership,
			row.Status == "Functional",
		)

	}

	fmt.Println("facility seeding completed")

	return nil
}

func normalizeNHFR(v string) string {
	v = strings.TrimSpace(v)

	if v == "" {
		return ""
	}

	if strings.Contains(v, "E+") {
		f, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return fmt.Sprintf("%.0f", f)
		}
	}

	return v
}
