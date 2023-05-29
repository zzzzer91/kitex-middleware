package logid

import (
	"context"
	"strings"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/zzzzer91/zlog"
)

func InjectLogIdMW(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request, response any) error {
		if v, ok := ctx.Value(zlog.EntityFieldNameLogId).(string); ok {
			ctx = metainfo.WithPersistentValue(ctx, zlog.EntityFieldNameLogId.String(), v)
		}
		if v, ok := ctx.Value(zlog.EntityFieldNameTraceId).(string); ok {
			ctx = metainfo.WithPersistentValue(ctx, zlog.EntityFieldNameTraceId.String(), v)
		}
		if v, ok := ctx.Value(zlog.EntityFieldNameRequestId).(string); ok {
			ctx = metainfo.WithPersistentValue(ctx, zlog.EntityFieldNameRequestId.String(), v)
		}
		err := next(ctx, request, response)
		return err
	}
}

func ExtractLogIdMW(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request, response any) error {
		// 传过来的 key 都会变成全大写的
		if v, ok := metainfo.GetPersistentValue(ctx, strings.ToUpper(zlog.EntityFieldNameLogId.String())); ok {
			ctx = context.WithValue(ctx, zlog.EntityFieldNameLogId, v)
		}
		if v, ok := metainfo.GetPersistentValue(ctx, strings.ToUpper(zlog.EntityFieldNameTraceId.String())); ok {
			ctx = context.WithValue(ctx, zlog.EntityFieldNameTraceId, v)
		}
		if v, ok := metainfo.GetPersistentValue(ctx, strings.ToUpper(zlog.EntityFieldNameRequestId.String())); ok {
			ctx = context.WithValue(ctx, zlog.EntityFieldNameRequestId, v)
		}
		err := next(ctx, request, response)
		return err
	}
}
