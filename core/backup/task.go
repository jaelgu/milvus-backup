package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/zilliztech/milvus-backup/core/client/milvus"
	"github.com/zilliztech/milvus-backup/core/meta"
	"github.com/zilliztech/milvus-backup/core/paramtable"
	"github.com/zilliztech/milvus-backup/core/proto/backuppb"
	"github.com/zilliztech/milvus-backup/core/storage"
	"github.com/zilliztech/milvus-backup/core/storage/mpath"
	"github.com/zilliztech/milvus-backup/core/utils"
	"github.com/zilliztech/milvus-backup/internal/log"
	"github.com/zilliztech/milvus-backup/internal/namespace"
)

const (
	_rpcChWarnMessage = "Failed to back up RPC channel position. This won't cause the backup to fail, " +
		"but may lead to inconsistency when reconnecting to CDC for incremental data replication."
)

type Task struct {
	backupID string

	logger *zap.Logger

	milvusStorage storage.Client
	backupStorage storage.Client
	copySem       *semaphore.Weighted

	params *paramtable.BackupParams

	grpc    milvus.Grpc
	restful milvus.Restful

	gcCtrl *gcController

	meta *meta.MetaManager

	request *backuppb.CreateBackupRequest
}

func NewTask(
	backupID string,
	MilvusStorage storage.Client,
	BackupStorage storage.Client,
	request *backuppb.CreateBackupRequest,
	params *paramtable.BackupParams,
	grpc milvus.Grpc,
	restful milvus.Restful,
	meta *meta.MetaManager,
) (*Task, error) {
	logger := log.L().With(zap.String("backup_id", backupID))
	gcCtrl, err := newGCController(request, params)
	if err != nil {
		return nil, fmt.Errorf("backup: new gc controller: %w", err)
	}

	return &Task{
		backupID: backupID,

		logger: logger,

		milvusStorage: MilvusStorage,
		backupStorage: BackupStorage,
		copySem:       semaphore.NewWeighted(params.BackupCfg.BackupCopyDataParallelism),

		params: params,

		grpc:    grpc,
		restful: restful,

		gcCtrl: gcCtrl,

		meta: meta,

		request: request,
	}, nil
}

type collection struct {
	DBName   string
	CollName string
}

func (t *Task) Execute(ctx context.Context) error {
	if err := t.privateExecute(ctx); err != nil {
		return err
	}

	return nil
}

func (t *Task) privateExecute(ctx context.Context) error {
	t.gcCtrl.Pause(ctx)
	defer t.gcCtrl.Resume(ctx)

	dbNames, collections, err := t.listDBAndCollection(ctx)
	if err != nil {
		return fmt.Errorf("backup: list db and collection: %w", err)
	}

	if err := t.runDBTask(ctx, dbNames); err != nil {
		return fmt.Errorf("backup: run db task: %w", err)
	}

	if err := t.runCollTask(ctx, collections); err != nil {
		return fmt.Errorf("backup: run collection task: %w", err)
	}

	if err := t.runRBACTask(ctx); err != nil {
		return fmt.Errorf("backup: run rbac task: %w", err)
	}

	t.runRPCChannelPOSTask(ctx)

	t.logger.Info("backup successfully")
	return nil
}

func (t *Task) listDBAndCollectionFromDBCollections(ctx context.Context, dbCollectionsStr string) ([]string, []collection, error) {
	var dbCollections meta.DbCollections
	if err := json.Unmarshal([]byte(dbCollectionsStr), &dbCollections); err != nil {
		return nil, nil, fmt.Errorf("backup: unmarshal dbCollections: %w", err)
	}

	var dbNames []string
	var collections []collection
	for db, colls := range dbCollections {
		if db == "" {
			db = namespace.DefaultDBName
		}
		dbNames = append(dbNames, db)

		// if collections is empty, list all collections in the database
		if len(colls) == 0 {
			resp, err := t.grpc.ListCollections(ctx, db)
			if err != nil {
				return nil, nil, fmt.Errorf("backup: list collection for db %s: %w", db, err)
			}
			colls = resp.CollectionNames
		}

		for _, coll := range colls {
			ns := namespace.New(db, coll)
			exist, err := t.grpc.HasCollection(ctx, ns.DBName(), ns.CollName())
			if err != nil {
				return nil, nil, fmt.Errorf("backup: check collection %s: %w", ns, err)
			}
			if !exist {
				return nil, nil, fmt.Errorf("backup: collection %s does not exist", ns)
			}
			t.logger.Debug("need backup collection", zap.String("db", db), zap.String("collection", coll))
			collections = append(collections, collection{DBName: db, CollName: coll})
		}
	}

	return dbNames, collections, nil
}

func (t *Task) listDBAndCollectionFromCollectionNames(ctx context.Context, collectionNames []string) ([]string, []collection, error) {
	dbNames := make(map[string]struct{})
	var collections []collection
	for _, collectionName := range collectionNames {
		dbName := namespace.DefaultDBName
		if strings.Contains(collectionName, ".") {
			splits := strings.Split(collectionName, ".")
			dbName = splits[0]
			collectionName = splits[1]
		}

		dbNames[dbName] = struct{}{}
		exist, err := t.grpc.HasCollection(ctx, dbName, collectionName)
		if err != nil {
			return nil, nil, fmt.Errorf("backup: check collection %s.%s: %w", dbName, collectionName, err)
		}
		if !exist {
			return nil, nil, fmt.Errorf("backup: collection %s.%s does not exist", dbName, collectionName)
		}
		collections = append(collections, collection{DBName: dbName, CollName: collectionName})
	}

	return maps.Keys(dbNames), collections, nil
}

func (t *Task) listDBAndCollectionFromAPI(ctx context.Context) ([]string, []collection, error) {
	var dbNames []string
	var collections []collection
	// compatible to milvus under v2.2.8 without database support
	if t.grpc.HasFeature(milvus.MultiDatabase) {
		t.logger.Info("the milvus server support multi database")
		dbs, err := t.grpc.ListDatabases(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("backup: list databases: %w", err)
		}
		for _, db := range dbs {
			t.logger.Debug("need backup db", zap.String("db", db))
			dbNames = append(dbNames, db)
			resp, err := t.grpc.ListCollections(ctx, db)
			if err != nil {
				return nil, nil, fmt.Errorf("backup: list collections for db %s: %w", db, err)
			}
			for _, coll := range resp.CollectionNames {
				t.logger.Debug("need backup collection", zap.String("db", db), zap.String("collection", coll))
				collections = append(collections, collection{DBName: db, CollName: coll})
			}
		}
	} else {
		t.logger.Info("the milvus server does not support multi database")
		dbNames = append(dbNames, namespace.DefaultDBName)
		resp, err := t.grpc.ListCollections(ctx, namespace.DefaultDBName)
		if err != nil {
			return nil, nil, fmt.Errorf("backup: list collections for db %s: %w", namespace.DefaultDBName, err)
		}
		for _, coll := range resp.CollectionNames {
			t.logger.Debug("need backup collection", zap.String("db", namespace.DefaultDBName), zap.String("collection", coll))
			collections = append(collections, collection{DBName: namespace.DefaultDBName, CollName: coll})
		}
	}

	return dbNames, collections, nil
}

// listDBAndCollection lists the database and collection names to be backed up.
// 1. Use dbCollection in the request if it is not empty.
// 2. If dbCollection is empty, use collectionNames in the request.
// 3. If collectionNames is empty, list all databases and collections.
func (t *Task) listDBAndCollection(ctx context.Context) ([]string, []collection, error) {
	dbCollectionsStr := utils.GetDBCollections(t.request.GetDbCollections())
	t.logger.Debug("get dbCollections from request", zap.String("dbCollections", dbCollectionsStr))

	// 1. dbCollections
	if dbCollectionsStr != "" {
		t.logger.Info("read need backup db and collection from dbCollection", zap.String("dbCollections", dbCollectionsStr))
		dbNames, collections, err := t.listDBAndCollectionFromDBCollections(ctx, dbCollectionsStr)
		if err != nil {
			return nil, nil, fmt.Errorf("backup: list db and collection from dbCollections: %w", err)
		}
		t.logger.Info("list db and collection from dbCollections done", zap.Strings("dbNames", dbNames), zap.Any("collections", collections))
		return dbNames, collections, nil
	}

	// 2. collectionNames
	if len(t.request.GetCollectionNames()) > 0 {
		t.logger.Info("read need backup db and collection from collectionNames", zap.Strings("collectionNames", t.request.GetCollectionNames()))
		dbNames, collections, err := t.listDBAndCollectionFromCollectionNames(ctx, t.request.GetCollectionNames())
		if err != nil {
			return nil, nil, fmt.Errorf("backup: list db and collection from collectionNames: %w", err)
		}
		t.logger.Info("list db and collection from collectionNames done", zap.Strings("dbNames", dbNames), zap.Any("collections", collections))
		return dbNames, collections, nil
	}

	// 3. list all databases and collections
	t.logger.Info("read need backup db and collection from API")
	dbNames, collections, err := t.listDBAndCollectionFromAPI(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("backup: list db and collection from API: %w", err)
	}
	t.logger.Info("list db and collection from API done", zap.Strings("dbNames", dbNames), zap.Any("collections", collections))

	return dbNames, collections, nil
}

func (t *Task) runDBTask(ctx context.Context, dbNames []string) error {
	for _, dbName := range dbNames {
		dbTask := NewDatabaseTask(t.backupID, dbName, t.grpc, t.meta)
		if err := dbTask.Execute(ctx); err != nil {
			return fmt.Errorf("backup: execute db task %s: %w", dbName, err)
		}
	}

	t.logger.Info("backup db done")

	return nil
}

func (t *Task) newCollectionTaskOpt() CollectionOpt {
	backupDir := mpath.BackupDir(t.params.MinioCfg.BackupRootPath, t.request.GetBackupName())

	crossStorage := t.params.MinioCfg.CrossStorage
	if t.backupStorage.Config().Provider != t.milvusStorage.Config().Provider {
		crossStorage = true
	}

	return CollectionOpt{
		BackupID:       t.request.GetRequestId(),
		MetaOnly:       t.request.GetMetaOnly(),
		SkipFlush:      t.request.GetForce(),
		MilvusStorage:  t.milvusStorage,
		MilvusRootPath: t.params.MinioCfg.RootPath,
		CrossStorage:   crossStorage,
		BackupStorage:  t.backupStorage,
		CopySem:        t.copySem,
		BackupDir:      backupDir,
		Meta:           t.meta,
		Grpc:           t.grpc,
		Restful:        t.restful,
	}
}

func (t *Task) runCollTask(ctx context.Context, collections []collection) error {
	t.logger.Info("start backup all collections", zap.Any("collections", collections))
	opt := t.newCollectionTaskOpt()

	g, subCtx := errgroup.WithContext(ctx)
	g.SetLimit(t.params.BackupCfg.BackupCollectionParallelism)
	for _, coll := range collections {
		g.Go(func() error {
			task := NewCollectionTask(coll.DBName, coll.CollName, opt)
			if err := task.Execute(subCtx); err != nil {
				return fmt.Errorf("create: backup collection %s.%s failed, err: %w", coll.DBName, coll.CollName, err)
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("backup: wait worker pool: %w", err)
	}

	t.logger.Info("backup all collections successfully")

	return nil
}

func (t *Task) runRBACTask(ctx context.Context) error {
	if !t.request.GetRbac() {
		t.logger.Info("skip backup rbac")
		return nil
	}

	rt := NewRBACTask(t.backupID, t.meta, t.grpc)
	if err := rt.Execute(ctx); err != nil {
		return fmt.Errorf("backup: execute rbac task: %w", err)
	}

	return nil
}

func (t *Task) runRPCChannelPOSTask(ctx context.Context) {
	t.logger.Info("start backup rpc channel pos")
	rpcPosTask := NewRPCChannelPOSTask(t.backupID, t.params.MilvusCfg.RPCChanelName, t.grpc, t.meta)
	if err := rpcPosTask.Execute(ctx); err != nil {
		t.logger.Warn(_rpcChWarnMessage, zap.Error(err))
		return
	}
	t.logger.Info("backup rpc channel pos done")
}
