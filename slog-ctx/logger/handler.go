package logger

import (
	"context"
	"log/slog"
	"slices"
)

type contextField struct {
	label string
	key   any
}

// PrependKey context에는 PrependKey{}:[]slog.Attr 형태로 저장
type PrependKey struct{}

// groupOrAttrs group과 attrs를 저장하는 구조체
type groupOrAttrs struct {
	group string
	attrs []slog.Attr
}

// interface 준수
var _ slog.Handler = (*Handler)(nil)

type Handler struct {
	next                 slog.Handler
	defaultContextFields []contextField
	goas                 []groupOrAttrs
}

// NewHandler custom Handler 생성자
//   - next : 로깅시 사용할 핸들러 json, text 등
//   - HandlerOptionFunc : handler 설정 옵션 함수
func NewHandler(next slog.Handler, ctxFieldOpts ...HandlerOptionFunc) *Handler {
	h := &Handler{next: next}
	for _, fieldOpt := range ctxFieldOpts {
		fieldOpt(h)
	}

	return h
}

// HandlerOptionFunc handler 설정 옵션 함수
type HandlerOptionFunc func(*Handler)

// WithContextField handler에 context key와 label을 추가
//   - key : context에 value를 저장할 때 사용할 key
//   - label : 로깅할 때 사용할 key
func WithContextField(label string, key any) HandlerOptionFunc {
	return func(h *Handler) {
		h.defaultContextFields = append(h.defaultContextFields, contextField{label: label, key: key})
	}
}

// Enabled 주어진 log level을 처리하도록 next handler가 설정되었는지 여부
func (h *Handler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return h.next.Enabled(ctx, lvl)
}

// Handle log record를 처리
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	attrs := make([]slog.Attr, r.NumAttrs())

	// r.attrs를 attrs로 복사
	r.Attrs(func(a slog.Attr) bool {
		attrs = append(attrs, a)
		return true
	})

	// goa 처리: 나중에 추가된 Group부터 안쪽에 추가하기 위해 역순으로 처리
	for i := len(h.goas) - 1; i >= 0; i-- {
		goa := h.goas[i]
		if goa.group != "" {
			attrs = []slog.Attr{{Key: goa.group, Value: slog.GroupValue(attrs...)}}
		} else {
			attrs = slices.Concat(goa.attrs, attrs)
		}
	}

	// prepend 처리
	if v, ok := ctx.Value(PrependKey{}).([]slog.Attr); ok {
		attrs = slices.Concat(v, attrs)
	}

	// defaultContextFields를 attrs 추가
	// 역순으로 추가하여 먼저 추가된 필드가 먼저 로깅되도록 함
	// 다른 attr보다 항상 앞에 로깅되도록 0번째에 추가
	for i := len(h.defaultContextFields) - 1; i >= 0; i-- {
		cf := h.defaultContextFields[i]
		ctxVal := ctx.Value(cf.key)
		if ctxVal == nil { // nil인 경우 로깅하지 않음
			continue
		}
		attr := ConverTContextValueToSlogAttr(cf.label, ctxVal)
		attrs = slices.Insert(attrs, 0, attr)

	}

	// attrs 추가
	rec := r.Clone() // Clone은 완전히 독립적인 복제본 생성
	rec.AddAttrs(attrs...)

	return h.next.Handle(ctx, rec)
}

// WithAttrs implements slog.Handler.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h2 := *h // copy
	h2.goas = append(h2.goas, groupOrAttrs{attrs: attrs})

	return &h2
}

// WithGroup 새 Handler를 반환하며 기존 handler의 ??
func (h *Handler) WithGroup(name string) slog.Handler {
	h2 := *h // copy
	h2.goas = append(h2.goas, groupOrAttrs{group: name})

	return &h2
}

// Prepend context를 이용해 default로 로깅할 필드를 추가해두는 함수로 다른 로깅 필드보다 Prepend로 추가된 필드가 먼저 로깅됨.
func Prepend(parentCtx context.Context, attrs ...slog.Attr) context.Context {
	if parentCtx == nil {
		parentCtx = context.Background()
	}

	if v, ok := parentCtx.Value(PrependKey{}).([]slog.Attr); ok {
		// slices.Concat: 두 슬라이스를 합친 새 슬라이스 반환
		return context.WithValue(parentCtx, PrependKey{}, slices.Concat(v, attrs))
	}

	return context.WithValue(parentCtx, PrependKey{}, attrs)
}
