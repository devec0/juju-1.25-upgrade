// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package state

import (
	"fmt"

	"github.com/juju/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2/txn"

	"github.com/juju/1.25-upgrade/juju2/constraints"
	"github.com/juju/1.25-upgrade/juju2/instance"
)

// constraintsDoc is the mongodb representation of a constraints.Value.
type constraintsDoc struct {
	ModelUUID    string `bson:"model-uuid"`
	Arch         *string
	CpuCores     *uint64
	CpuPower     *uint64
	Mem          *uint64
	RootDisk     *uint64
	InstanceType *string
	Container    *instance.ContainerType
	Tags         *[]string
	Spaces       *[]string
	VirtType     *string
}

func (doc constraintsDoc) value() constraints.Value {
	result := constraints.Value{
		Arch:         doc.Arch,
		CpuCores:     doc.CpuCores,
		CpuPower:     doc.CpuPower,
		Mem:          doc.Mem,
		RootDisk:     doc.RootDisk,
		InstanceType: doc.InstanceType,
		Container:    doc.Container,
		Tags:         doc.Tags,
		Spaces:       doc.Spaces,
		VirtType:     doc.VirtType,
	}
	return result
}

func newConstraintsDoc(st *State, cons constraints.Value) constraintsDoc {
	result := constraintsDoc{
		Arch:         cons.Arch,
		CpuCores:     cons.CpuCores,
		CpuPower:     cons.CpuPower,
		Mem:          cons.Mem,
		RootDisk:     cons.RootDisk,
		InstanceType: cons.InstanceType,
		Container:    cons.Container,
		Tags:         cons.Tags,
		Spaces:       cons.Spaces,
		VirtType:     cons.VirtType,
	}
	return result
}

func createConstraintsOp(st *State, id string, cons constraints.Value) txn.Op {
	return txn.Op{
		C:      constraintsC,
		Id:     id,
		Assert: txn.DocMissing,
		Insert: newConstraintsDoc(st, cons),
	}
}

func setConstraintsOp(st *State, id string, cons constraints.Value) txn.Op {
	return txn.Op{
		C:      constraintsC,
		Id:     id,
		Assert: txn.DocExists,
		Update: bson.D{{"$set", newConstraintsDoc(st, cons)}},
	}
}

func removeConstraintsOp(st *State, id string) txn.Op {
	return txn.Op{
		C:      constraintsC,
		Id:     id,
		Remove: true,
	}
}

func readConstraints(st *State, id string) (constraints.Value, error) {
	constraintsCollection, closer := st.db().GetCollection(constraintsC)
	defer closer()

	doc := constraintsDoc{}
	if err := constraintsCollection.FindId(id).One(&doc); err == mgo.ErrNotFound {
		return constraints.Value{}, errors.NotFoundf("constraints")
	} else if err != nil {
		return constraints.Value{}, err
	}
	return doc.value(), nil
}

func writeConstraints(st *State, id string, cons constraints.Value) error {
	ops := []txn.Op{setConstraintsOp(st, id, cons)}
	if err := st.runTransaction(ops); err != nil {
		return fmt.Errorf("cannot set constraints: %v", err)
	}
	return nil
}
