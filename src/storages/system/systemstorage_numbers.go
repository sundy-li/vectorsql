// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package system

import (
	"base/errors"
	"columns"
	"datablocks"
	"datatypes"
	"datavalues"
	"planners"
	"sessions"
)

type SystemNumbersStorage struct {
	ctx *SystemStorageContext
}

func NewSystemNumbersStorage(ctx *SystemStorageContext) *SystemNumbersStorage {
	return &SystemNumbersStorage{
		ctx: ctx,
	}
}

func (storage *SystemNumbersStorage) Name() string {
	return ""
}

func (storage *SystemNumbersStorage) Columns() []columns.Column {
	return []columns.Column{
		{Name: "number", DataType: datatypes.NewUInt64DataType()},
	}
}

func (storage *SystemNumbersStorage) GetOutputStream(session *sessions.Session, scan *planners.ScanPlan) (datablocks.IDataBlockOutputStream, error) {
	return nil, errors.New("Couldn't find outputstream")
}

func (storage *SystemNumbersStorage) GetInputStream(session *sessions.Session, scan *planners.ScanPlan) (datablocks.IDataBlockInputStream, error) {
	return NewSystemNumbersBlockInputStream(storage), nil
}

type SystemNumbersBlockIntputStream struct {
	storage      *SystemNumbersStorage
	block        *datablocks.DataBlock
	offset       int
	limit        int
	maxBlockSize int
	current      int
}

func NewSystemNumbersBlockInputStream(storage *SystemNumbersStorage) *SystemNumbersBlockIntputStream {
	return &SystemNumbersBlockIntputStream{
		storage:      storage,
		block:        datablocks.NewDataBlock(storage.Columns()),
		maxBlockSize: storage.ctx.conf.Server.DefaultBlockSize,
	}
}

func (stream *SystemNumbersBlockIntputStream) Name() string {
	return "SystemNumbersBlockIntputStream"
}

func (stream *SystemNumbersBlockIntputStream) Read() (block *datablocks.DataBlock, err error) {
	var rows = 0
	block = stream.block.Clone()
	batcher := datablocks.NewBatchWriter(block.Columns())

	for rows < stream.maxBlockSize {
		if err := batcher.WriteRow(datavalues.MakeInt(stream.current)); err != nil {
			return nil, err
		}
		stream.current++
		rows++
	}

	if rows == 0 {
		return nil, nil
	}
	block.WriteBatch(batcher)
	return
}
