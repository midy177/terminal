// Code generated by ent, DO NOT EDIT.

package hosts

import (
	"terminal/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Hosts {
	return predicate.Hosts(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Hosts {
	return predicate.Hosts(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Hosts {
	return predicate.Hosts(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Hosts {
	return predicate.Hosts(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Hosts {
	return predicate.Hosts(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Hosts {
	return predicate.Hosts(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Hosts {
	return predicate.Hosts(sql.FieldLTE(FieldID, id))
}

// Username applies equality check predicate on the "username" field. It's identical to UsernameEQ.
func Username(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldUsername, v))
}

// Address applies equality check predicate on the "address" field. It's identical to AddressEQ.
func Address(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldAddress, v))
}

// Port applies equality check predicate on the "port" field. It's identical to PortEQ.
func Port(v uint) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldPort, v))
}

// Password applies equality check predicate on the "password" field. It's identical to PasswordEQ.
func Password(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldPassword, v))
}

// FolderID applies equality check predicate on the "folder_id" field. It's identical to FolderIDEQ.
func FolderID(v int) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldFolderID, v))
}

// KeyID applies equality check predicate on the "key_id" field. It's identical to KeyIDEQ.
func KeyID(v int) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldKeyID, v))
}

// LabelEQ applies the EQ predicate on the "label" field.
func LabelEQ(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldLabel, v))
}

// LabelNEQ applies the NEQ predicate on the "label" field.
func LabelNEQ(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldNEQ(FieldLabel, v))
}

// LabelIn applies the In predicate on the "label" field.
func LabelIn(vs ...string) predicate.Hosts {
	return predicate.Hosts(sql.FieldIn(FieldLabel, vs...))
}

// LabelNotIn applies the NotIn predicate on the "label" field.
func LabelNotIn(vs ...string) predicate.Hosts {
	return predicate.Hosts(sql.FieldNotIn(FieldLabel, vs...))
}

// LabelGT applies the GT predicate on the "label" field.
func LabelGT(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldGT(FieldLabel, v))
}

// LabelGTE applies the GTE predicate on the "label" field.
func LabelGTE(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldGTE(FieldLabel, v))
}

// LabelLT applies the LT predicate on the "label" field.
func LabelLT(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldLT(FieldLabel, v))
}

// LabelLTE applies the LTE predicate on the "label" field.
func LabelLTE(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldLTE(FieldLabel, v))
}

// LabelContains applies the Contains predicate on the "label" field.
func LabelContains(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldContains(FieldLabel, v))
}

// LabelHasPrefix applies the HasPrefix predicate on the "label" field.
func LabelHasPrefix(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldHasPrefix(FieldLabel, v))
}

// LabelHasSuffix applies the HasSuffix predicate on the "label" field.
func LabelHasSuffix(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldHasSuffix(FieldLabel, v))
}

// LabelEqualFold applies the EqualFold predicate on the "label" field.
func LabelEqualFold(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldEqualFold(FieldLabel, v))
}

// LabelContainsFold applies the ContainsFold predicate on the "label" field.
func LabelContainsFold(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldContainsFold(FieldLabel, v))
}

// UsernameEQ applies the EQ predicate on the "username" field.
func UsernameEQ(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldUsername, v))
}

// UsernameNEQ applies the NEQ predicate on the "username" field.
func UsernameNEQ(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldNEQ(FieldUsername, v))
}

// UsernameIn applies the In predicate on the "username" field.
func UsernameIn(vs ...string) predicate.Hosts {
	return predicate.Hosts(sql.FieldIn(FieldUsername, vs...))
}

// UsernameNotIn applies the NotIn predicate on the "username" field.
func UsernameNotIn(vs ...string) predicate.Hosts {
	return predicate.Hosts(sql.FieldNotIn(FieldUsername, vs...))
}

// UsernameGT applies the GT predicate on the "username" field.
func UsernameGT(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldGT(FieldUsername, v))
}

// UsernameGTE applies the GTE predicate on the "username" field.
func UsernameGTE(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldGTE(FieldUsername, v))
}

// UsernameLT applies the LT predicate on the "username" field.
func UsernameLT(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldLT(FieldUsername, v))
}

// UsernameLTE applies the LTE predicate on the "username" field.
func UsernameLTE(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldLTE(FieldUsername, v))
}

// UsernameContains applies the Contains predicate on the "username" field.
func UsernameContains(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldContains(FieldUsername, v))
}

// UsernameHasPrefix applies the HasPrefix predicate on the "username" field.
func UsernameHasPrefix(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldHasPrefix(FieldUsername, v))
}

// UsernameHasSuffix applies the HasSuffix predicate on the "username" field.
func UsernameHasSuffix(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldHasSuffix(FieldUsername, v))
}

// UsernameEqualFold applies the EqualFold predicate on the "username" field.
func UsernameEqualFold(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldEqualFold(FieldUsername, v))
}

// UsernameContainsFold applies the ContainsFold predicate on the "username" field.
func UsernameContainsFold(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldContainsFold(FieldUsername, v))
}

// AddressEQ applies the EQ predicate on the "address" field.
func AddressEQ(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldAddress, v))
}

// AddressNEQ applies the NEQ predicate on the "address" field.
func AddressNEQ(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldNEQ(FieldAddress, v))
}

// AddressIn applies the In predicate on the "address" field.
func AddressIn(vs ...string) predicate.Hosts {
	return predicate.Hosts(sql.FieldIn(FieldAddress, vs...))
}

// AddressNotIn applies the NotIn predicate on the "address" field.
func AddressNotIn(vs ...string) predicate.Hosts {
	return predicate.Hosts(sql.FieldNotIn(FieldAddress, vs...))
}

// AddressGT applies the GT predicate on the "address" field.
func AddressGT(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldGT(FieldAddress, v))
}

// AddressGTE applies the GTE predicate on the "address" field.
func AddressGTE(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldGTE(FieldAddress, v))
}

// AddressLT applies the LT predicate on the "address" field.
func AddressLT(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldLT(FieldAddress, v))
}

// AddressLTE applies the LTE predicate on the "address" field.
func AddressLTE(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldLTE(FieldAddress, v))
}

// AddressContains applies the Contains predicate on the "address" field.
func AddressContains(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldContains(FieldAddress, v))
}

// AddressHasPrefix applies the HasPrefix predicate on the "address" field.
func AddressHasPrefix(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldHasPrefix(FieldAddress, v))
}

// AddressHasSuffix applies the HasSuffix predicate on the "address" field.
func AddressHasSuffix(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldHasSuffix(FieldAddress, v))
}

// AddressEqualFold applies the EqualFold predicate on the "address" field.
func AddressEqualFold(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldEqualFold(FieldAddress, v))
}

// AddressContainsFold applies the ContainsFold predicate on the "address" field.
func AddressContainsFold(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldContainsFold(FieldAddress, v))
}

// PortEQ applies the EQ predicate on the "port" field.
func PortEQ(v uint) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldPort, v))
}

// PortNEQ applies the NEQ predicate on the "port" field.
func PortNEQ(v uint) predicate.Hosts {
	return predicate.Hosts(sql.FieldNEQ(FieldPort, v))
}

// PortIn applies the In predicate on the "port" field.
func PortIn(vs ...uint) predicate.Hosts {
	return predicate.Hosts(sql.FieldIn(FieldPort, vs...))
}

// PortNotIn applies the NotIn predicate on the "port" field.
func PortNotIn(vs ...uint) predicate.Hosts {
	return predicate.Hosts(sql.FieldNotIn(FieldPort, vs...))
}

// PortGT applies the GT predicate on the "port" field.
func PortGT(v uint) predicate.Hosts {
	return predicate.Hosts(sql.FieldGT(FieldPort, v))
}

// PortGTE applies the GTE predicate on the "port" field.
func PortGTE(v uint) predicate.Hosts {
	return predicate.Hosts(sql.FieldGTE(FieldPort, v))
}

// PortLT applies the LT predicate on the "port" field.
func PortLT(v uint) predicate.Hosts {
	return predicate.Hosts(sql.FieldLT(FieldPort, v))
}

// PortLTE applies the LTE predicate on the "port" field.
func PortLTE(v uint) predicate.Hosts {
	return predicate.Hosts(sql.FieldLTE(FieldPort, v))
}

// PasswordEQ applies the EQ predicate on the "password" field.
func PasswordEQ(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldPassword, v))
}

// PasswordNEQ applies the NEQ predicate on the "password" field.
func PasswordNEQ(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldNEQ(FieldPassword, v))
}

// PasswordIn applies the In predicate on the "password" field.
func PasswordIn(vs ...string) predicate.Hosts {
	return predicate.Hosts(sql.FieldIn(FieldPassword, vs...))
}

// PasswordNotIn applies the NotIn predicate on the "password" field.
func PasswordNotIn(vs ...string) predicate.Hosts {
	return predicate.Hosts(sql.FieldNotIn(FieldPassword, vs...))
}

// PasswordGT applies the GT predicate on the "password" field.
func PasswordGT(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldGT(FieldPassword, v))
}

// PasswordGTE applies the GTE predicate on the "password" field.
func PasswordGTE(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldGTE(FieldPassword, v))
}

// PasswordLT applies the LT predicate on the "password" field.
func PasswordLT(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldLT(FieldPassword, v))
}

// PasswordLTE applies the LTE predicate on the "password" field.
func PasswordLTE(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldLTE(FieldPassword, v))
}

// PasswordContains applies the Contains predicate on the "password" field.
func PasswordContains(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldContains(FieldPassword, v))
}

// PasswordHasPrefix applies the HasPrefix predicate on the "password" field.
func PasswordHasPrefix(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldHasPrefix(FieldPassword, v))
}

// PasswordHasSuffix applies the HasSuffix predicate on the "password" field.
func PasswordHasSuffix(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldHasSuffix(FieldPassword, v))
}

// PasswordIsNil applies the IsNil predicate on the "password" field.
func PasswordIsNil() predicate.Hosts {
	return predicate.Hosts(sql.FieldIsNull(FieldPassword))
}

// PasswordNotNil applies the NotNil predicate on the "password" field.
func PasswordNotNil() predicate.Hosts {
	return predicate.Hosts(sql.FieldNotNull(FieldPassword))
}

// PasswordEqualFold applies the EqualFold predicate on the "password" field.
func PasswordEqualFold(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldEqualFold(FieldPassword, v))
}

// PasswordContainsFold applies the ContainsFold predicate on the "password" field.
func PasswordContainsFold(v string) predicate.Hosts {
	return predicate.Hosts(sql.FieldContainsFold(FieldPassword, v))
}

// FolderIDEQ applies the EQ predicate on the "folder_id" field.
func FolderIDEQ(v int) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldFolderID, v))
}

// FolderIDNEQ applies the NEQ predicate on the "folder_id" field.
func FolderIDNEQ(v int) predicate.Hosts {
	return predicate.Hosts(sql.FieldNEQ(FieldFolderID, v))
}

// FolderIDIn applies the In predicate on the "folder_id" field.
func FolderIDIn(vs ...int) predicate.Hosts {
	return predicate.Hosts(sql.FieldIn(FieldFolderID, vs...))
}

// FolderIDNotIn applies the NotIn predicate on the "folder_id" field.
func FolderIDNotIn(vs ...int) predicate.Hosts {
	return predicate.Hosts(sql.FieldNotIn(FieldFolderID, vs...))
}

// FolderIDIsNil applies the IsNil predicate on the "folder_id" field.
func FolderIDIsNil() predicate.Hosts {
	return predicate.Hosts(sql.FieldIsNull(FieldFolderID))
}

// FolderIDNotNil applies the NotNil predicate on the "folder_id" field.
func FolderIDNotNil() predicate.Hosts {
	return predicate.Hosts(sql.FieldNotNull(FieldFolderID))
}

// KeyIDEQ applies the EQ predicate on the "key_id" field.
func KeyIDEQ(v int) predicate.Hosts {
	return predicate.Hosts(sql.FieldEQ(FieldKeyID, v))
}

// KeyIDNEQ applies the NEQ predicate on the "key_id" field.
func KeyIDNEQ(v int) predicate.Hosts {
	return predicate.Hosts(sql.FieldNEQ(FieldKeyID, v))
}

// KeyIDIn applies the In predicate on the "key_id" field.
func KeyIDIn(vs ...int) predicate.Hosts {
	return predicate.Hosts(sql.FieldIn(FieldKeyID, vs...))
}

// KeyIDNotIn applies the NotIn predicate on the "key_id" field.
func KeyIDNotIn(vs ...int) predicate.Hosts {
	return predicate.Hosts(sql.FieldNotIn(FieldKeyID, vs...))
}

// KeyIDIsNil applies the IsNil predicate on the "key_id" field.
func KeyIDIsNil() predicate.Hosts {
	return predicate.Hosts(sql.FieldIsNull(FieldKeyID))
}

// KeyIDNotNil applies the NotNil predicate on the "key_id" field.
func KeyIDNotNil() predicate.Hosts {
	return predicate.Hosts(sql.FieldNotNull(FieldKeyID))
}

// HasFolder applies the HasEdge predicate on the "folder" edge.
func HasFolder() predicate.Hosts {
	return predicate.Hosts(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, FolderTable, FolderColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasFolderWith applies the HasEdge predicate on the "folder" edge with a given conditions (other predicates).
func HasFolderWith(preds ...predicate.Folders) predicate.Hosts {
	return predicate.Hosts(func(s *sql.Selector) {
		step := newFolderStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasKey applies the HasEdge predicate on the "key" edge.
func HasKey() predicate.Hosts {
	return predicate.Hosts(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, KeyTable, KeyColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasKeyWith applies the HasEdge predicate on the "key" edge with a given conditions (other predicates).
func HasKeyWith(preds ...predicate.Keys) predicate.Hosts {
	return predicate.Hosts(func(s *sql.Selector) {
		step := newKeyStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Hosts) predicate.Hosts {
	return predicate.Hosts(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Hosts) predicate.Hosts {
	return predicate.Hosts(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Hosts) predicate.Hosts {
	return predicate.Hosts(sql.NotPredicates(p))
}
