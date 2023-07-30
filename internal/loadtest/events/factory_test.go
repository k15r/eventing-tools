package events

import (
	"testing"
	"time"

	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kyma-project/eventing-tools/internal/loadtest/api/subscription/v1alpha2"
)

func TestFactory_reconcile(t *testing.T) {
	type fields struct {
		generators map[NamespaceName]eventGenerator
		senderC    chan Event
	}
	type args struct {
		sub *v1alpha2.Subscription
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "new subscription",
			fields: fields{
				generators: map[NamespaceName]eventGenerator{},
				senderC:    make(chan Event),
			},
			args: args{
				sub: &v1alpha2.Subscription{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "new",
						Namespace: "new",
					},
					Spec: v1alpha2.SubscriptionSpec{
						Sink:         "",
						TypeMatching: v1alpha2.Standard,
						Source:       "source",
						Types:        []string{"foo.bar.v1", "bar.foo.v1"},
						Config:       nil,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			go func(ctx context.Context) {
				for {
					select {
					case <-ctx.Done():
						return
					case e := <-tt.fields.senderC:
						t.Logf("%+v", e)
					}
				}
			}(ctx)
			f := &Factory{
				generators: tt.fields.generators,
				senderC:    tt.fields.senderC,
			}
			if err := f.reconcile(tt.args.sub); (err != nil) != tt.wantErr {
				t.Errorf("reconcile() error = %v, wantErr %v", err, tt.wantErr)
			}
			time.Sleep(1 * time.Second)
			f.Stop()
			cancel()
		})
	}
}
