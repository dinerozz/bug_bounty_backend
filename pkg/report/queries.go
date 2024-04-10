package report

var ReportsTableQuery = `SELECT 
    c.name, r.title, r.status 
FROM 
    reports r 
    LEFT JOIN categories c ON r.category_id = c.id 
    JOIN users u ON r.author_id = u.id 
    JOIN team_members tm ON u.id = tm.user_id 
WHERE 
    tm.team_id IN (SELECT team_id FROM team_members WHERE user_id = $1) order by r.id`

var DetailsQuery = `SELECT 
    rr.reviewer_id, u.username, rr.review_text, 
    r.id, u2.username, r.title, r.description, r.status, c.name 
FROM 
    report_reviews rr 
JOIN reports r ON rr.report_id = r.id 
	LEFT JOIN users u ON rr.reviewer_id = u.id 
	LEFT JOIN categories c ON r.category_id = c.id 
	JOIN users u2 ON r.author_id = u2.id 
	JOIN team_members tm ON u2.id = tm.user_id 
	JOIN team_members tm2 ON tm.team_id = tm2.team_id AND tm2.user_id = $2 
WHERE rr.report_id = $1`

var ReportsQuery = `SELECT r.id, u.username, c.name, r.title, r.status, r.description 
				FROM reports r
				LEFT JOIN users u on r.author_id = u.id
				LEFT JOIN categories c on r.category_id = c.id
				JOIN users u2 ON r.author_id = u2.id 
				JOIN team_members tm ON u2.id = tm.user_id 
				JOIN team_members tm2 ON tm.team_id = tm2.team_id AND tm2.user_id = $2 
            	WHERE r.id = $1`
