package service

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/Notailab/Notailab/internal/dao"
)

type StatsService struct {
	dao *dao.StatsDAO
}

type StatsTrendPoint struct {
	Label string `json:"label"`
	Value int    `json:"value"`
}

type StatsProjectVelocity struct {
	ProjectID uint   `json:"project_id"`
	Title     string `json:"title"`
	Progress  int    `json:"progress"`
	FileCount int    `json:"file_count"`
}

type StatsOverview struct {
	TotalProjects        int64                  `json:"total_projects"`
	TotalFiles           int64                  `json:"total_files"`
	TotalWords           int                    `json:"total_words"`
	ActiveProjects       int64                  `json:"active_projects"`
	AIConversations      int64                  `json:"ai_conversations"`
	RecentProjectUpdates int64                  `json:"recent_project_updates"`
	RecentFileUpdates    int64                  `json:"recent_file_updates"`
	MonthlyActivity      []StatsTrendPoint      `json:"monthly_activity"`
	ProjectVelocity      []StatsProjectVelocity `json:"project_velocity"`
	Insights             []string               `json:"insights"`
	LastUpdated          time.Time              `json:"last_updated"`
}

func NewStatsService(statsDAO *dao.StatsDAO) *StatsService {
	return &StatsService{dao: statsDAO}
}

var latinWordPattern = regexp.MustCompile(`[A-Za-z0-9]+(?:'[A-Za-z0-9]+)?`)

func countWords(text string) int {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return 0
	}

	latinWords := len(latinWordPattern.FindAllString(trimmed, -1))
	cjkChars := 0
	for _, r := range trimmed {
		if unicode.Is(unicode.Han, r) {
			cjkChars++
		}
	}
	return latinWords + cjkChars
}

func (s *StatsService) GetOverview(userID uint) (*StatsOverview, error) {
	projects, err := s.dao.ListProjectsByUserID(userID)
	if err != nil {
		return nil, err
	}

	files, err := s.dao.ListFilesByUserID(userID)
	if err != nil {
		return nil, err
	}

	totalProjects, err := s.dao.CountProjectsByUserID(userID)
	if err != nil {
		return nil, err
	}

	totalFiles, err := s.dao.CountFilesByUserID(userID)
	if err != nil {
		return nil, err
	}

	aiConversations, err := s.dao.CountAgentConversationsByUserID(userID)
	if err != nil {
		return nil, err
	}

	recentWindow := time.Now().AddDate(0, 0, -30)
	recentProjectUpdates, err := s.dao.CountRecentProjectsByUserID(userID, recentWindow)
	if err != nil {
		return nil, err
	}

	recentFileUpdates, err := s.dao.CountRecentFilesByUserID(userID, recentWindow)
	if err != nil {
		return nil, err
	}

	monthlyActivity, err := s.dao.MonthlyFileActivityByUserID(userID, 6)
	if err != nil {
		return nil, err
	}

	projectVelocityRows, err := s.dao.TopProjectsByFileCount(userID, 4)
	if err != nil {
		return nil, err
	}

	totalWords := 0
	activeProjects := int64(0)
	for _, project := range projects {
		if project.Progress < 100 {
			activeProjects++
		}
	}

	for _, file := range files {
		totalWords += countWords(file.Content)
	}

	trend := make([]StatsTrendPoint, 0, len(monthlyActivity))
	for _, item := range monthlyActivity {
		trend = append(trend, StatsTrendPoint{Label: item.Month, Value: item.Count})
	}

	velocity := make([]StatsProjectVelocity, 0, len(projectVelocityRows))
	for _, item := range projectVelocityRows {
		velocity = append(velocity, StatsProjectVelocity{
			ProjectID: item.ProjectID,
			Title:     item.Title,
			Progress:  item.Progress,
			FileCount: item.FileCount,
		})
	}

	sort.SliceStable(velocity, func(i, j int) bool {
		if velocity[i].FileCount == velocity[j].FileCount {
			return velocity[i].Progress > velocity[j].Progress
		}
		return velocity[i].FileCount > velocity[j].FileCount
	})

	insights := buildInsights(totalProjects, totalFiles, totalWords, activeProjects, aiConversations, recentFileUpdates, velocity)

	return &StatsOverview{
		TotalProjects:        totalProjects,
		TotalFiles:           totalFiles,
		TotalWords:           totalWords,
		ActiveProjects:       activeProjects,
		AIConversations:      aiConversations,
		RecentProjectUpdates: recentProjectUpdates,
		RecentFileUpdates:    recentFileUpdates,
		MonthlyActivity:      trend,
		ProjectVelocity:      velocity,
		Insights:             insights,
		LastUpdated:          time.Now(),
	}, nil
}

func buildInsights(totalProjects, totalFiles int64, totalWords int, activeProjects, aiConversations, recentFileUpdates int64, velocity []StatsProjectVelocity) []string {
	if totalProjects == 0 {
		return []string{"先创建一个项目，Stats 面板就会开始积累你的学习轨迹。"}
	}

	insights := make([]string, 0, 4)

	if recentFileUpdates > 0 {
		insights = append(insights, "最近 30 天你还有持续更新，适合把这些内容整理成阶段总结。")
	} else {
		insights = append(insights, "最近 30 天更新较少，建议先从一个最常用项目继续推进。")
	}

	if len(velocity) > 0 {
		insights = append(insights, "当前最活跃的项目是「"+velocity[0].Title+"」，可以优先复盘和补充结构化内容。")
	}

	if activeProjects > 0 {
		insights = append(insights, "你当前有 "+strconv.FormatInt(activeProjects, 10)+" 个未完成项目，建议优先推进正在进行中的内容。")
	}

	if aiConversations > 0 {
		insights = append(insights, "你已经和 AI 进行过多次协作，可以进一步提炼成提示词模板。")
	}

	if totalWords > 0 && totalFiles > 0 {
		average := float64(totalWords) / float64(totalFiles)
		insights = append(insights, "平均每个文件约 "+strings.TrimRight(strings.TrimRight(strconv.FormatFloat(average, 'f', 1, 64), "0"), ".")+" 个词，内容密度还不错。")
	}

	return insights
}
