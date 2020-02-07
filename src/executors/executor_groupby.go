// Copyright 2020 The VectorSQL Authors.
//
// Code is licensed under Apache License, Version 2.0.

package executors

import (
	"planners"
	"processors"
	"transforms"
)

type GroupByExecutor struct {
	ctx  *ExecutorContext
	plan *planners.GroupByPlan
}

func NewGroupByExecutor(ctx *ExecutorContext, plan *planners.GroupByPlan) *GroupByExecutor {
	return &GroupByExecutor{
		ctx:  ctx,
		plan: plan,
	}
}

func (executor *GroupByExecutor) Execute() (processors.IProcessor, error) {
	log := executor.ctx.log
	conf := executor.ctx.conf

	log.Debug("Executor->Enter->LogicalPlan:%s", executor.plan)
	transformCtx := transforms.NewTransformContext(executor.ctx.ctx, log, conf)
	transform := transforms.NewGroupByTransform(transformCtx, executor.plan)
	log.Debug("Executor->Return->Pipeline:%v", transform)
	return transform, nil
}
