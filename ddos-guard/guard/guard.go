package guard

import (
	"container/list"
	"crypto/rand"
	"math/big"
	"sort"
	"sync/atomic"
	"time"

	"github.com/cornelk/hashmap"
	"github.com/davecgh/go-spew/spew"

	"github.com/yukels/ddos-guard/config"
	"github.com/yukels/util/context"
	"github.com/yukels/util/log"
)

type UserCompareFunc func(i, j *UserStats) bool

type UserStats struct {
	Count      int64
	DurationMs int64
}

func (s *UserStats) Add(ctx context.Context, stats *UserStats) {
	s.Count += stats.Count
	s.DurationMs += stats.DurationMs
}

func (s *UserStats) Sub(ctx context.Context, stats *UserStats) {
	s.Count -= stats.Count
	s.DurationMs -= stats.DurationMs
}

type Bucket struct {
	startTime time.Time
	users     *hashmap.Map[string, *UserStats]
}

func NewBucket(ctx context.Context) *Bucket {
	return &Bucket{
		startTime: time.Now().UTC(),
		users:     hashmap.New[string, *UserStats](),
	}
}

type Guard struct {
	config             *config.GuardConfig
	monitor            *Monitor
	stats              *list.List
	latestBucket       *Bucket
	collector          *Collector
	topUsersByCount    map[string]*UserStats
	topUsersByDuration map[string]*UserStats
	totalUser          map[string]*UserStats
	totalAvg           UserStats
	filterRatio        int64
}

func NewGuard(ctx context.Context, collector *Collector) (*Guard, error) {
	monitor, err := NewMonitor(ctx, collector, &config.Configs.DdosGuardConfig.Monitoring)
	if err != nil {
		return nil, err
	}
	g := &Guard{
		config:             &config.Configs.DdosGuardConfig.Guard,
		monitor:            monitor,
		stats:              list.New(),
		collector:          collector,
		topUsersByCount:    map[string]*UserStats{},
		topUsersByDuration: map[string]*UserStats{},
		totalUser:          map[string]*UserStats{},
		filterRatio:        0,
	}

	return g, nil
}

func (g *Guard) Run(ctx context.Context) error {
	if err := g.monitor.Run(ctx); err != nil {
		return err
	}

	g.addBucket(ctx)
	go g.statsLoop(ctx)
	return nil
}

func (g *Guard) statsLoop(ctx context.Context) {
	// time-shift between monitoring and top-users threads
	time.Sleep(time.Second * 5)
	log.Log(ctx).Info("Guard stats thread is running...")
	waitPeriod := g.config.BucketDuration
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(waitPeriod):
			g.updateStats(ctx)
		}
	}
}

func (g *Guard) addBucket(ctx context.Context) {
	g.latestBucket = NewBucket(ctx)
	g.stats.PushBack(g.latestBucket)
}

func (g *Guard) updateStats(ctx context.Context) {
	rotate := g.stats.Len() >= g.config.BucketsHistory
	oldest, newest := g.stats.Front(), g.stats.Back()
	g.addBucket(ctx)
	g.calcTopUsers(ctx, rotate, oldest.Value.(*Bucket), newest.Value.(*Bucket))
	if rotate {
		g.stats.Remove(oldest)
	}
}

func (g *Guard) calcTopUsers(ctx context.Context, rotate bool, oldest, newest *Bucket) {
	g.calcTotalUserCalls(ctx, rotate, oldest, newest)
	g.topUsersByCount, g.topUsersByDuration = g.nTopUsers(ctx)

	// If we still under ddos attack - increase the "filter ratio" of the incoming requests
	// Otherwise - decrease. ideally back to "filter ratio" = 0
	ratioStep := -g.config.FilterRatioStep
	if g.monitor.InHighUsage(ctx) {
		ratioStep = g.config.FilterRatioStep
	}
	g.updateFilterRatio(ctx, ratioStep)

	log.Log(ctx).Infof(spew.Sprintf("Top users by count [%+v], users by duration [%+v]", g.topUsersByCount, g.topUsersByDuration))
}

func (g *Guard) updateFilterRatio(ctx context.Context, ratioStep int64) {
	filterRatio := g.filterRatio + ratioStep
	if filterRatio > 100 {
		filterRatio = 100
	} else if filterRatio < 0 {
		filterRatio = 0
	}
	g.filterRatio = filterRatio
	log.Log(ctx).Infof("New filterRatio [%d]", g.filterRatio)
}

func (g *Guard) nTopUsers(ctx context.Context) (map[string]*UserStats, map[string]*UserStats) {
	return g.nTopUsersBy(ctx, func(i, j *UserStats) bool { return i.Count > j.Count }),
		g.nTopUsersBy(ctx, func(i, j *UserStats) bool { return i.DurationMs > j.DurationMs })
}

func (g *Guard) nTopUsersBy(ctx context.Context, compare UserCompareFunc) map[string]*UserStats {
	// Find N top elements in list. Use sort - not effective but less lines of code :))
	users := make([]string, 0, len(g.totalUser))
	for user := range g.totalUser {
		users = append(users, user)
	}

	sort.SliceStable(users, func(i, j int) bool {
		return compare(g.totalUser[users[i]], g.totalUser[users[j]])
	})

	topUsers := users
	if len(users) > g.config.TopUserCount {
		topUsers = users[:g.config.TopUserCount]
	}

	topUsersMap := map[string]*UserStats{}
	for _, user := range topUsers {
		topUsersMap[user] = g.totalUser[user]
	}

	return topUsersMap
}

// Update from 'oldest' - remove from total user statistics
// and one before the 'newest' which the first time taken into user statistics
func (g *Guard) calcTotalUserCalls(ctx context.Context, rotate bool, oldest, newest *Bucket) {
	oldestStats := UserStats{}
	if rotate {
		// remove statistics from the "oldest" bucket
		oldest.users.Range(func(user string, calls *UserStats) bool {
			g.totalUser[user].Sub(ctx, calls)
			oldestStats.Add(ctx, calls)
			return true
		})
	}
	g.totalAvg.Sub(ctx, &oldestStats)

	// add statistics from the "newest" bucket
	newestStats := UserStats{}
	newest.users.Range(func(user string, calls *UserStats) bool {
		userStats, ok := g.totalUser[user]
		if !ok {
			userStats = &UserStats{}
			g.totalUser[user] = userStats
		}
		userStats.Add(ctx, calls)
		newestStats.Add(ctx, calls)
		return true
	})
	g.totalAvg.Add(ctx, &newestStats)
}

func (g *Guard) ShouldBlockUser(ctx context.Context, user string) bool {
	if !g.monitor.InHighUsage(ctx) {
		return false
	}
	_, exist := g.topUsersByCount[user]
	if exist {
		return g.toBeOrNotToBe(ctx)
	}
	_, exist = g.topUsersByDuration[user]
	if exist {
		return g.toBeOrNotToBe(ctx)
	}
	return false
}

func (g *Guard) toBeOrNotToBe(ctx context.Context) bool {
	nBig, _ := rand.Int(rand.Reader, big.NewInt(100))
	return nBig.Int64() < g.filterRatio
}

func (g *Guard) RequestCompleteUser(ctx context.Context, user string, duration time.Duration) {
	if user == "" {
		return
	}
	g.callRegister(ctx, user, duration)
}

func (g *Guard) callRegister(ctx context.Context, user string, duration time.Duration) {
	stats := UserStats{}
	userStats, _ := g.latestBucket.users.GetOrInsert(user, &stats)
	atomic.AddInt64(&userStats.Count, 1)
	atomic.AddInt64(&userStats.DurationMs, duration.Milliseconds())
}
