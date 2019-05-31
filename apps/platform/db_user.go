package platform

import "wave/apps/platform/models"

// --------------------------|

func (db db) CreatePrifile(p models.UserProfile) (err error) {
	_, err = db.Exec(`
		INSERT INTO users(uid, name, avatar)
		VALUES ($1, $2, $3)
	`, p.UID, p.Username, p.Avatar)
	return err
}

func (db db) GetProfile(uid int64) (p models.UserProfile, err error) {
	if err := db.QueryRow(`
		SELECT uid, name, avatar
			FROM users
			WHERE uid=$1
	`, uid).Scan(&p.UID, &p.Username, &p.Avatar); err != nil {
		return p, err
	}
	return p, nil
}
