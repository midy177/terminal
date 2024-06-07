// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"terminal/ent/folders"
	"terminal/ent/hosts"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// FoldersCreate is the builder for creating a Folders entity.
type FoldersCreate struct {
	config
	mutation *FoldersMutation
	hooks    []Hook
}

// SetLabel sets the "label" field.
func (fc *FoldersCreate) SetLabel(s string) *FoldersCreate {
	fc.mutation.SetLabel(s)
	return fc
}

// SetParentID sets the "parent_id" field.
func (fc *FoldersCreate) SetParentID(i int) *FoldersCreate {
	fc.mutation.SetParentID(i)
	return fc
}

// SetNillableParentID sets the "parent_id" field if the given value is not nil.
func (fc *FoldersCreate) SetNillableParentID(i *int) *FoldersCreate {
	if i != nil {
		fc.SetParentID(*i)
	}
	return fc
}

// SetParent sets the "parent" edge to the Folders entity.
func (fc *FoldersCreate) SetParent(f *Folders) *FoldersCreate {
	return fc.SetParentID(f.ID)
}

// AddChildIDs adds the "children" edge to the Folders entity by IDs.
func (fc *FoldersCreate) AddChildIDs(ids ...int) *FoldersCreate {
	fc.mutation.AddChildIDs(ids...)
	return fc
}

// AddChildren adds the "children" edges to the Folders entity.
func (fc *FoldersCreate) AddChildren(f ...*Folders) *FoldersCreate {
	ids := make([]int, len(f))
	for i := range f {
		ids[i] = f[i].ID
	}
	return fc.AddChildIDs(ids...)
}

// SetHostID sets the "host" edge to the Hosts entity by ID.
func (fc *FoldersCreate) SetHostID(id int) *FoldersCreate {
	fc.mutation.SetHostID(id)
	return fc
}

// SetNillableHostID sets the "host" edge to the Hosts entity by ID if the given value is not nil.
func (fc *FoldersCreate) SetNillableHostID(id *int) *FoldersCreate {
	if id != nil {
		fc = fc.SetHostID(*id)
	}
	return fc
}

// SetHost sets the "host" edge to the Hosts entity.
func (fc *FoldersCreate) SetHost(h *Hosts) *FoldersCreate {
	return fc.SetHostID(h.ID)
}

// Mutation returns the FoldersMutation object of the builder.
func (fc *FoldersCreate) Mutation() *FoldersMutation {
	return fc.mutation
}

// Save creates the Folders in the database.
func (fc *FoldersCreate) Save(ctx context.Context) (*Folders, error) {
	return withHooks(ctx, fc.sqlSave, fc.mutation, fc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (fc *FoldersCreate) SaveX(ctx context.Context) *Folders {
	v, err := fc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fc *FoldersCreate) Exec(ctx context.Context) error {
	_, err := fc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fc *FoldersCreate) ExecX(ctx context.Context) {
	if err := fc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (fc *FoldersCreate) check() error {
	if _, ok := fc.mutation.Label(); !ok {
		return &ValidationError{Name: "label", err: errors.New(`ent: missing required field "Folders.label"`)}
	}
	if v, ok := fc.mutation.Label(); ok {
		if err := folders.LabelValidator(v); err != nil {
			return &ValidationError{Name: "label", err: fmt.Errorf(`ent: validator failed for field "Folders.label": %w`, err)}
		}
	}
	return nil
}

func (fc *FoldersCreate) sqlSave(ctx context.Context) (*Folders, error) {
	if err := fc.check(); err != nil {
		return nil, err
	}
	_node, _spec := fc.createSpec()
	if err := sqlgraph.CreateNode(ctx, fc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	fc.mutation.id = &_node.ID
	fc.mutation.done = true
	return _node, nil
}

func (fc *FoldersCreate) createSpec() (*Folders, *sqlgraph.CreateSpec) {
	var (
		_node = &Folders{config: fc.config}
		_spec = sqlgraph.NewCreateSpec(folders.Table, sqlgraph.NewFieldSpec(folders.FieldID, field.TypeInt))
	)
	if value, ok := fc.mutation.Label(); ok {
		_spec.SetField(folders.FieldLabel, field.TypeString, value)
		_node.Label = value
	}
	if nodes := fc.mutation.ParentIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   folders.ParentTable,
			Columns: []string{folders.ParentColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(folders.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.ParentID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.ChildrenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   folders.ChildrenTable,
			Columns: []string{folders.ChildrenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(folders.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := fc.mutation.HostIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   folders.HostTable,
			Columns: []string{folders.HostColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(hosts.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// FoldersCreateBulk is the builder for creating many Folders entities in bulk.
type FoldersCreateBulk struct {
	config
	err      error
	builders []*FoldersCreate
}

// Save creates the Folders entities in the database.
func (fcb *FoldersCreateBulk) Save(ctx context.Context) ([]*Folders, error) {
	if fcb.err != nil {
		return nil, fcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(fcb.builders))
	nodes := make([]*Folders, len(fcb.builders))
	mutators := make([]Mutator, len(fcb.builders))
	for i := range fcb.builders {
		func(i int, root context.Context) {
			builder := fcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*FoldersMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, fcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, fcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, fcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (fcb *FoldersCreateBulk) SaveX(ctx context.Context) []*Folders {
	v, err := fcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (fcb *FoldersCreateBulk) Exec(ctx context.Context) error {
	_, err := fcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fcb *FoldersCreateBulk) ExecX(ctx context.Context) {
	if err := fcb.Exec(ctx); err != nil {
		panic(err)
	}
}