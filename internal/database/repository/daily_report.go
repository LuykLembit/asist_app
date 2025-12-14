package repository

import (
	"database/sql"
	"fmt"
	"time"

	"telemonitor/internal/database"
)

// DailyReportRepository handles daily_reports operations
type DailyReportRepository struct {
	db *database.DB
}

// NewDailyReportRepository creates a new DailyReportRepository
func NewDailyReportRepository(db *database.DB) *DailyReportRepository {
	return &DailyReportRepository{db: db}
}

// Create inserts a new daily report
func (r *DailyReportRepository) Create(report *database.DailyReport) error {
	query := `
		INSERT INTO daily_reports (chat_id, report_date, summary, full_json)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (chat_id, report_date) DO UPDATE
		SET summary = EXCLUDED.summary, full_json = EXCLUDED.full_json, created_at = NOW()
		RETURNING id, created_at
	`
	
	err := r.db.QueryRow(query,
		report.ChatID,
		report.ReportDate,
		report.Summary,
		report.FullJSON,
	).Scan(&report.ID, &report.CreatedAt)
	
	if err != nil {
		return fmt.Errorf("failed to create daily report: %w", err)
	}
	
	return nil
}

// GetByID retrieves a report by ID
func (r *DailyReportRepository) GetByID(id int) (*database.DailyReport, error) {
	query := `
		SELECT id, chat_id, report_date, summary, full_json, created_at
		FROM daily_reports
		WHERE id = $1
	`
	
	report := &database.DailyReport{}
	err := r.db.QueryRow(query, id).Scan(
		&report.ID,
		&report.ChatID,
		&report.ReportDate,
		&report.Summary,
		&report.FullJSON,
		&report.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get report: %w", err)
	}
	
	return report, nil
}

// GetByChatAndDate retrieves a report for a specific chat and date
func (r *DailyReportRepository) GetByChatAndDate(chatID int64, date time.Time) (*database.DailyReport, error) {
	query := `
		SELECT id, chat_id, report_date, summary, full_json, created_at
		FROM daily_reports
		WHERE chat_id = $1 AND report_date = $2
	`
	
	report := &database.DailyReport{}
	err := r.db.QueryRow(query, chatID, date.Format("2006-01-02")).Scan(
		&report.ID,
		&report.ChatID,
		&report.ReportDate,
		&report.Summary,
		&report.FullJSON,
		&report.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get report by chat and date: %w", err)
	}
	
	return report, nil
}

// GetByChatID retrieves all reports for a chat
func (r *DailyReportRepository) GetByChatID(chatID int64) ([]*database.DailyReport, error) {
	query := `
		SELECT id, chat_id, report_date, summary, full_json, created_at
		FROM daily_reports
		WHERE chat_id = $1
		ORDER BY report_date DESC
	`
	
	rows, err := r.db.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reports by chat: %w", err)
	}
	defer rows.Close()
	
	var reports []*database.DailyReport
	for rows.Next() {
		report := &database.DailyReport{}
		if err := rows.Scan(
			&report.ID,
			&report.ChatID,
			&report.ReportDate,
			&report.Summary,
			&report.FullJSON,
			&report.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan report: %w", err)
		}
		reports = append(reports, report)
	}
	
	return reports, nil
}

// Search performs a simple text search on report summaries
func (r *DailyReportRepository) Search(query string) ([]*database.DailyReport, error) {
	sqlQuery := `
		SELECT id, chat_id, report_date, summary, full_json, created_at
		FROM daily_reports
		WHERE summary ILIKE $1
		ORDER BY report_date DESC
		LIMIT 50
	`
	
	rows, err := r.db.Query(sqlQuery, "%"+query+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to search reports: %w", err)
	}
	defer rows.Close()
	
	var reports []*database.DailyReport
	for rows.Next() {
		report := &database.DailyReport{}
		if err := rows.Scan(
			&report.ID,
			&report.ChatID,
			&report.ReportDate,
			&report.Summary,
			&report.FullJSON,
			&report.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan report: %w", err)
		}
		reports = append(reports, report)
	}
	
	return reports, nil
}

// GetLatest retrieves the most recent reports
func (r *DailyReportRepository) GetLatest(limit int) ([]*database.DailyReport, error) {
	query := `
		SELECT id, chat_id, report_date, summary, full_json, created_at
		FROM daily_reports
		ORDER BY report_date DESC
		LIMIT $1
	`
	
	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest reports: %w", err)
	}
	defer rows.Close()
	
	var reports []*database.DailyReport
	for rows.Next() {
		report := &database.DailyReport{}
		if err := rows.Scan(
			&report.ID,
			&report.ChatID,
			&report.ReportDate,
			&report.Summary,
			&report.FullJSON,
			&report.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan report: %w", err)
		}
		reports = append(reports, report)
	}
	
	return reports, nil
}
