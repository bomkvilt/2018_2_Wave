package platform

import (
	"wave/apps/platform/models"

	"github.com/jackc/pgx"
)

const appFields = "aid, name, image, link, url, about, installations, category"

func scanApp(row *pgx.Row, a *models.App) (err error) {
	return row.Scan(&a.AID, &a.Name, &a.Image, &a.Link, &a.URL, &a.About, &a.Installations, &a.Category)
}

func scanApps(rows *pgx.Rows, apps *models.Apps) (err error) {
	for rows.Next() {
		a := models.App{}
		err = rows.Scan(&a.AID, &a.Name, &a.Image, &a.Link, &a.URL, &a.About, &a.Installations, &a.Category)
		if err != nil {
			return err
		}
		*apps = append(*apps, a)
	}
	return nil
}

// --------------------------|

func (db db) GetApps() (apps models.Apps, err error) {
	rows, err := db.Query(`
		SELECT ` + appFields + `
			FROM apps
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	err = scanApps(rows, &apps)
	return apps, nil
}

func (db db) GetAppCategories() (ctgs models.Categories, err error) {
	rows, err := db.Query(`
		SELECT DISTINCT category
		FROM apps
		ORDER BY category DESC
	`)
	if err != nil {
		return ctgs, err
	}

	defer rows.Close()
	for rows.Next() {
		ctg := ""
		err := rows.Scan(&ctg)
		if err != nil {
			return ctgs, err
		}
		ctgs = append(ctgs, ctg)
	}
	return ctgs, nil
}

func (db db) GetUserApps(uid int64) (apps models.Apps, err error) {
	rows, err := db.Query(`
		SELECT `+appFields+`
			FROM apps
			JOIN userapps USING(aid)
			WHERE uid=$1
	`, uid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	err = scanApps(rows, &apps)
	return apps, nil
}

func (db db) GetCategoryApps(ctg string) (apps models.Apps, err error) {
	rows, err := db.Query(`
		SELECT `+appFields+`
			FROM apps
			JOIN userapps USING(aid)
			WHERE category=$1
			ORDER BY installs DESC
	`, ctg)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	err = scanApps(rows, &apps)
	return apps, nil
}

func (db db) GetApp(name string) (app models.App, err error) {
	err = scanApp(db.QueryRow(`
		SELECT `+appFields+`
			FROM apps
			WHERE name=$1
	`, name), &app)
	return app, err
}

func (db db) AddMyApp(uid int64, name string) (err error) {
	_, err = db.Exec(`
		INSERT INTO userapps(uid, aid)
		VALUES ($1, (SELECT aid FROM apps WHERE name=$2))
		ON CONFLICT DO NOTHING
	`, uid, name)
	return err
}
