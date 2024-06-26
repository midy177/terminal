// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"terminal/ent/migrate"

	"terminal/ent/folders"
	"terminal/ent/hosts"
	"terminal/ent/keys"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Folders is the client for interacting with the Folders builders.
	Folders *FoldersClient
	// Hosts is the client for interacting with the Hosts builders.
	Hosts *HostsClient
	// Keys is the client for interacting with the Keys builders.
	Keys *KeysClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	client := &Client{config: newConfig(opts...)}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Folders = NewFoldersClient(c.config)
	c.Hosts = NewHostsClient(c.config)
	c.Keys = NewKeysClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// newConfig creates a new config for the client.
func newConfig(opts ...Option) config {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	return cfg
}

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:     ctx,
		config:  cfg,
		Folders: NewFoldersClient(cfg),
		Hosts:   NewHostsClient(cfg),
		Keys:    NewKeysClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:     ctx,
		config:  cfg,
		Folders: NewFoldersClient(cfg),
		Hosts:   NewHostsClient(cfg),
		Keys:    NewKeysClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Folders.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Folders.Use(hooks...)
	c.Hosts.Use(hooks...)
	c.Keys.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Folders.Intercept(interceptors...)
	c.Hosts.Intercept(interceptors...)
	c.Keys.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *FoldersMutation:
		return c.Folders.mutate(ctx, m)
	case *HostsMutation:
		return c.Hosts.mutate(ctx, m)
	case *KeysMutation:
		return c.Keys.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// FoldersClient is a client for the Folders schema.
type FoldersClient struct {
	config
}

// NewFoldersClient returns a client for the Folders from the given config.
func NewFoldersClient(c config) *FoldersClient {
	return &FoldersClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `folders.Hooks(f(g(h())))`.
func (c *FoldersClient) Use(hooks ...Hook) {
	c.hooks.Folders = append(c.hooks.Folders, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `folders.Intercept(f(g(h())))`.
func (c *FoldersClient) Intercept(interceptors ...Interceptor) {
	c.inters.Folders = append(c.inters.Folders, interceptors...)
}

// Create returns a builder for creating a Folders entity.
func (c *FoldersClient) Create() *FoldersCreate {
	mutation := newFoldersMutation(c.config, OpCreate)
	return &FoldersCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Folders entities.
func (c *FoldersClient) CreateBulk(builders ...*FoldersCreate) *FoldersCreateBulk {
	return &FoldersCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *FoldersClient) MapCreateBulk(slice any, setFunc func(*FoldersCreate, int)) *FoldersCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &FoldersCreateBulk{err: fmt.Errorf("calling to FoldersClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*FoldersCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &FoldersCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Folders.
func (c *FoldersClient) Update() *FoldersUpdate {
	mutation := newFoldersMutation(c.config, OpUpdate)
	return &FoldersUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *FoldersClient) UpdateOne(f *Folders) *FoldersUpdateOne {
	mutation := newFoldersMutation(c.config, OpUpdateOne, withFolders(f))
	return &FoldersUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *FoldersClient) UpdateOneID(id int) *FoldersUpdateOne {
	mutation := newFoldersMutation(c.config, OpUpdateOne, withFoldersID(id))
	return &FoldersUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Folders.
func (c *FoldersClient) Delete() *FoldersDelete {
	mutation := newFoldersMutation(c.config, OpDelete)
	return &FoldersDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *FoldersClient) DeleteOne(f *Folders) *FoldersDeleteOne {
	return c.DeleteOneID(f.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *FoldersClient) DeleteOneID(id int) *FoldersDeleteOne {
	builder := c.Delete().Where(folders.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &FoldersDeleteOne{builder}
}

// Query returns a query builder for Folders.
func (c *FoldersClient) Query() *FoldersQuery {
	return &FoldersQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeFolders},
		inters: c.Interceptors(),
	}
}

// Get returns a Folders entity by its id.
func (c *FoldersClient) Get(ctx context.Context, id int) (*Folders, error) {
	return c.Query().Where(folders.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *FoldersClient) GetX(ctx context.Context, id int) *Folders {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryParent queries the parent edge of a Folders.
func (c *FoldersClient) QueryParent(f *Folders) *FoldersQuery {
	query := (&FoldersClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := f.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(folders.Table, folders.FieldID, id),
			sqlgraph.To(folders.Table, folders.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, folders.ParentTable, folders.ParentColumn),
		)
		fromV = sqlgraph.Neighbors(f.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryChildren queries the children edge of a Folders.
func (c *FoldersClient) QueryChildren(f *Folders) *FoldersQuery {
	query := (&FoldersClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := f.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(folders.Table, folders.FieldID, id),
			sqlgraph.To(folders.Table, folders.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, folders.ChildrenTable, folders.ChildrenColumn),
		)
		fromV = sqlgraph.Neighbors(f.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryHost queries the host edge of a Folders.
func (c *FoldersClient) QueryHost(f *Folders) *HostsQuery {
	query := (&HostsClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := f.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(folders.Table, folders.FieldID, id),
			sqlgraph.To(hosts.Table, hosts.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, folders.HostTable, folders.HostColumn),
		)
		fromV = sqlgraph.Neighbors(f.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *FoldersClient) Hooks() []Hook {
	return c.hooks.Folders
}

// Interceptors returns the client interceptors.
func (c *FoldersClient) Interceptors() []Interceptor {
	return c.inters.Folders
}

func (c *FoldersClient) mutate(ctx context.Context, m *FoldersMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&FoldersCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&FoldersUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&FoldersUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&FoldersDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Folders mutation op: %q", m.Op())
	}
}

// HostsClient is a client for the Hosts schema.
type HostsClient struct {
	config
}

// NewHostsClient returns a client for the Hosts from the given config.
func NewHostsClient(c config) *HostsClient {
	return &HostsClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `hosts.Hooks(f(g(h())))`.
func (c *HostsClient) Use(hooks ...Hook) {
	c.hooks.Hosts = append(c.hooks.Hosts, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `hosts.Intercept(f(g(h())))`.
func (c *HostsClient) Intercept(interceptors ...Interceptor) {
	c.inters.Hosts = append(c.inters.Hosts, interceptors...)
}

// Create returns a builder for creating a Hosts entity.
func (c *HostsClient) Create() *HostsCreate {
	mutation := newHostsMutation(c.config, OpCreate)
	return &HostsCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Hosts entities.
func (c *HostsClient) CreateBulk(builders ...*HostsCreate) *HostsCreateBulk {
	return &HostsCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *HostsClient) MapCreateBulk(slice any, setFunc func(*HostsCreate, int)) *HostsCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &HostsCreateBulk{err: fmt.Errorf("calling to HostsClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*HostsCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &HostsCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Hosts.
func (c *HostsClient) Update() *HostsUpdate {
	mutation := newHostsMutation(c.config, OpUpdate)
	return &HostsUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *HostsClient) UpdateOne(h *Hosts) *HostsUpdateOne {
	mutation := newHostsMutation(c.config, OpUpdateOne, withHosts(h))
	return &HostsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *HostsClient) UpdateOneID(id int) *HostsUpdateOne {
	mutation := newHostsMutation(c.config, OpUpdateOne, withHostsID(id))
	return &HostsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Hosts.
func (c *HostsClient) Delete() *HostsDelete {
	mutation := newHostsMutation(c.config, OpDelete)
	return &HostsDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *HostsClient) DeleteOne(h *Hosts) *HostsDeleteOne {
	return c.DeleteOneID(h.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *HostsClient) DeleteOneID(id int) *HostsDeleteOne {
	builder := c.Delete().Where(hosts.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &HostsDeleteOne{builder}
}

// Query returns a query builder for Hosts.
func (c *HostsClient) Query() *HostsQuery {
	return &HostsQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeHosts},
		inters: c.Interceptors(),
	}
}

// Get returns a Hosts entity by its id.
func (c *HostsClient) Get(ctx context.Context, id int) (*Hosts, error) {
	return c.Query().Where(hosts.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *HostsClient) GetX(ctx context.Context, id int) *Hosts {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryFolder queries the folder edge of a Hosts.
func (c *HostsClient) QueryFolder(h *Hosts) *FoldersQuery {
	query := (&FoldersClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := h.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(hosts.Table, hosts.FieldID, id),
			sqlgraph.To(folders.Table, folders.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, hosts.FolderTable, hosts.FolderColumn),
		)
		fromV = sqlgraph.Neighbors(h.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// QueryKey queries the key edge of a Hosts.
func (c *HostsClient) QueryKey(h *Hosts) *KeysQuery {
	query := (&KeysClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := h.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(hosts.Table, hosts.FieldID, id),
			sqlgraph.To(keys.Table, keys.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, hosts.KeyTable, hosts.KeyColumn),
		)
		fromV = sqlgraph.Neighbors(h.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *HostsClient) Hooks() []Hook {
	return c.hooks.Hosts
}

// Interceptors returns the client interceptors.
func (c *HostsClient) Interceptors() []Interceptor {
	return c.inters.Hosts
}

func (c *HostsClient) mutate(ctx context.Context, m *HostsMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&HostsCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&HostsUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&HostsUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&HostsDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Hosts mutation op: %q", m.Op())
	}
}

// KeysClient is a client for the Keys schema.
type KeysClient struct {
	config
}

// NewKeysClient returns a client for the Keys from the given config.
func NewKeysClient(c config) *KeysClient {
	return &KeysClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `keys.Hooks(f(g(h())))`.
func (c *KeysClient) Use(hooks ...Hook) {
	c.hooks.Keys = append(c.hooks.Keys, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `keys.Intercept(f(g(h())))`.
func (c *KeysClient) Intercept(interceptors ...Interceptor) {
	c.inters.Keys = append(c.inters.Keys, interceptors...)
}

// Create returns a builder for creating a Keys entity.
func (c *KeysClient) Create() *KeysCreate {
	mutation := newKeysMutation(c.config, OpCreate)
	return &KeysCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Keys entities.
func (c *KeysClient) CreateBulk(builders ...*KeysCreate) *KeysCreateBulk {
	return &KeysCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *KeysClient) MapCreateBulk(slice any, setFunc func(*KeysCreate, int)) *KeysCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &KeysCreateBulk{err: fmt.Errorf("calling to KeysClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*KeysCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &KeysCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Keys.
func (c *KeysClient) Update() *KeysUpdate {
	mutation := newKeysMutation(c.config, OpUpdate)
	return &KeysUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *KeysClient) UpdateOne(k *Keys) *KeysUpdateOne {
	mutation := newKeysMutation(c.config, OpUpdateOne, withKeys(k))
	return &KeysUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *KeysClient) UpdateOneID(id int) *KeysUpdateOne {
	mutation := newKeysMutation(c.config, OpUpdateOne, withKeysID(id))
	return &KeysUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Keys.
func (c *KeysClient) Delete() *KeysDelete {
	mutation := newKeysMutation(c.config, OpDelete)
	return &KeysDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *KeysClient) DeleteOne(k *Keys) *KeysDeleteOne {
	return c.DeleteOneID(k.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *KeysClient) DeleteOneID(id int) *KeysDeleteOne {
	builder := c.Delete().Where(keys.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &KeysDeleteOne{builder}
}

// Query returns a query builder for Keys.
func (c *KeysClient) Query() *KeysQuery {
	return &KeysQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeKeys},
		inters: c.Interceptors(),
	}
}

// Get returns a Keys entity by its id.
func (c *KeysClient) Get(ctx context.Context, id int) (*Keys, error) {
	return c.Query().Where(keys.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *KeysClient) GetX(ctx context.Context, id int) *Keys {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryHost queries the host edge of a Keys.
func (c *KeysClient) QueryHost(k *Keys) *HostsQuery {
	query := (&HostsClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := k.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(keys.Table, keys.FieldID, id),
			sqlgraph.To(hosts.Table, hosts.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, keys.HostTable, keys.HostColumn),
		)
		fromV = sqlgraph.Neighbors(k.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *KeysClient) Hooks() []Hook {
	return c.hooks.Keys
}

// Interceptors returns the client interceptors.
func (c *KeysClient) Interceptors() []Interceptor {
	return c.inters.Keys
}

func (c *KeysClient) mutate(ctx context.Context, m *KeysMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&KeysCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&KeysUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&KeysUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&KeysDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Keys mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Folders, Hosts, Keys []ent.Hook
	}
	inters struct {
		Folders, Hosts, Keys []ent.Interceptor
	}
)
