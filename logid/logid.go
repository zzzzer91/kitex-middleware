package logid

import (
	"context"
	"strings"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/zzzzer91/zlog"
)

func InjectLogIDMw(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request, response any) error {
		if v, ok := ctx.Value(zlog.EntityFieldNameLogID).(string); ok {
			ctx = metainfo.WithPersistentValue(ctx, zlog.EntityFieldNameLogID.String(), v)
		}
		if v, ok := ctx.Value(zlog.EntityFieldNameTraceID).(string); ok {
			ctx = metainfo.WithPersistentValue(ctx, zlog.EntityFieldNameTraceID.String(), v)
		}
		if v, ok := ctx.Value(zlog.EntityFieldNameRequestID).(string); ok {
			ctx = metainfo.WithPersistentValue(ctx, zlog.EntityFieldNameRequestID.String(), v)
		}
		err := next(ctx, request, response)
		return err
	}
}

func ExtractLogIDMw(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request, response any) error {
		// 传过来的 key 都会变成全大写的
		if v, ok := metainfo.GetPersistentValue(ctx, strings.ToUpper(zlog.EntityFieldNameLogID.String())); ok {
			ctx = context.WithValue(ctx, zlog.EntityFieldNameLogID, v)
		}
		if v, ok := metainfo.GetPersistentValue(ctx, strings.ToUpper(zlog.EntityFieldNameTraceID.String())); ok {
			ctx = context.WithValue(ctx, zlog.EntityFieldNameTraceID, v)
		}
		if v, ok := metainfo.GetPersistentValue(ctx, strings.ToUpper(zlog.EntityFieldNameRequestID.String())); ok {
			ctx = context.WithValue(ctx, zlog.EntityFieldNameRequestID, v)
		}
		err := next(ctx, request, response)
		return err
	}
}
